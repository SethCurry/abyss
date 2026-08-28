package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SethCurry/abyss/internal/acptools"
	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/SethCurry/abyss/internal/api/agentapi"
	"github.com/SethCurry/abyss/internal/erres"
	"github.com/SethCurry/abyss/internal/runenv"
	api "github.com/SethCurry/abyss/internal/websockets/wsacp"
	"github.com/moby/moby/api/types/container"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

const (
	defaultImage = "ghcr.io/sethcurry/abyss-pi:latest"
	serverPort   = 8080
)

func main() {
	logFile, err := openLogFile()
	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	logOut := io.MultiWriter(logFile, os.Stderr)
	logger := zerolog.New(logOut).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	log.Logger = logger

	cmd := &cli.Command{
		Name:        "abyss",
		Usage:       "Agent Runtime Environment(s)",
		ArgsUsage:   "",
		Description: "Manage agents just like containers.",
		Version:     time.Now().Format(time.RFC3339),
		Authors:     []any{"Seth Curry"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"l"},
				Usage:   "The level to log at.  One of trace, debug, info, warn, error, fatal, disabled",
				Value:   "debug",
				Sources: cli.EnvVars("ABYSS_LOG_LEVEL"),
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			level, err := zerolog.ParseLevel(cmd.String("log-level"))
			if err != nil {
				return nil, fmt.Errorf("invalid log level %q: %w", cmd.String("log-level"), err)
			}
			logger = logger.Level(level)
			log.Logger = logger
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Starts the client-side proxy, which will create its own server-side proxy container running your agent.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"f"},
						Usage:    "The path to the agent configuration YAML file.",
						Required: true,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					configPath := cmd.String("config")
					agentCfg, err := agentconfig.FromYAMLFile(configPath)
					if err != nil {
						logger.Error().Err(err).Str("config_path", configPath).Msg("failed to load agent config")
						return err
					}
					logger.Debug().
						Str("config_path", configPath).
						Str("image", agentCfg.Docker.Image).
						Strs("agent_command", agentCfg.Docker.AgentCommand).
						Msg("starting Docker agent")
					return runClient(ctx, agentCfg, logger)
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Starts the agent-side proxy.  You should never need to manually invoke this.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "addr",
						Aliases: []string{"a"},
						Usage:   "The address to run the HTTP server on.",
						Value:   ":8080",
					},
					&cli.StringSliceFlag{
						Name:     "agent",
						Aliases:  []string{"g"},
						Usage:    "The agent command to run.  Specify this flag multiple times if there is more than one part to the command (e.g. \"npx my-package\" would be \"-g npx -g my-package\"",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "local-terminal",
						Aliases: []string{""},
						Usage:   "Run terminal ACP commands on the client rather than this server",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "local-filesystem",
						Aliases: []string{""},
						Usage:   "Run filesystem ACP commands on the client rather than this server",
						Value:   false,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					agentCmd := cmd.StringSlice("agent")
					logger.Debug().Strs("agent_command", agentCmd).Msg("starting agent and websocket server")

					var localTerminal *acptools.TerminalTools
					var localFilesystem *acptools.FilesystemTools

					if cmd.Bool("local-terminal") {
						localTerminal = acptools.NewTerminalTools(logger)
					}
					if cmd.Bool("local-filesystem") {
						localFilesystem = acptools.NewFilesystemTools(logger)
					}

					httpSrv := agentapi.NewServer(agentCmd, localTerminal, localFilesystem)
					return httpSrv.Serve(cmd.String("addr"))
				},
			},
			{
				Name:    "docker",
				Aliases: []string{"d"},
				Usage:   "Docker-related utilities.",
				Commands: []*cli.Command{
					{
						Name:    "ps",
						Aliases: []string{"p"},
						Usage:   "List Abyss containers that are currently running.",
						Action: func(ctx context.Context, cmd *cli.Command) error {
							docker, err := runenv.NewDockerClient()
							if err != nil {
								return fmt.Errorf("failed to connect to docker: %w", err)
							}

							containers, err := docker.AbyssContainers(ctx)
							if err != nil {
								return fmt.Errorf("failed to list containers: %w", err)
							}

							for _, v := range containers {
								joinedNames := strings.Join(v.Names, ", ")
								fmt.Println(v.ID + ": " + joinedNames)
							}
							return nil
						},
					},
				},
			},
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		if humanErr, ok := errors.AsType[erres.HumanError](err); ok {
			logger.Error().Str("error", humanErr.HumanError()).Msg("command failed")
		} else {
			logger.Error().Err(err).Msg("command failed")
		}
	}
}

// newFileLogger returns a zerolog.Logger configured to write debug-level logs
// to ~/.local/share/abyss/abyss.log, creating the directory if necessary.
func openLogFile() (io.WriteCloser, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	logDir := filepath.Join(home, ".local", "var", "abyss")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to make parent directories of %q: %w", logDir, err)
	}

	logFilePath := filepath.Join(logDir, "abyss.log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed open log file %q: %w", logFilePath, err)
	}

	return logFile, nil
}

// runClient starts a Docker container running the abyss server command and
// bridges it to a client over stdio.
func runClient(ctx context.Context, cfg *agentconfig.AgentConfig, logger zerolog.Logger) error {
	docker, err := runenv.NewDockerClient()
	if err != nil {
		return err
	}
	defer docker.Close()

	image := cfg.Docker.Image
	if image == "" {
		image = defaultImage
	}
	agent := cfg.Docker.AgentCommand

	agentArgs := make([]string, len(agent)*2)

	for k, v := range agent {
		startIndex := k * 2
		agentArgs[startIndex] = "--agent"
		agentArgs[startIndex+1] = v
	}

	if cfg.ACP.ToolsOnHost.Files {
		agentArgs = append(agentArgs, "--local-filesystem")
	}

	if cfg.ACP.ToolsOnHost.Terminal {
		agentArgs = append(agentArgs, "--local-terminal")
	}

	joinedArgs := strings.Join(agentArgs, " ")

	config := &container.Config{
		Entrypoint: []string{"/bin/bash", "-c"},
		Cmd:        []string{"/usr/local/bin/abyss server " + joinedArgs},
	}

	hostConfig := docker.ApplyHostMounts(&cfg.Docker, nil)

	container, endpoint, err := docker.StartContainer(ctx, image, config, hostConfig, "", serverPort, 0)
	if err != nil {
		logger.Error().Err(err).Msg("failed to start container")
		return err
	}

	if err := copyFiles(ctx, docker, container, cfg.CopyFiles, logger); err != nil {
		logger.Error().Err(err).Msg("failed to copy files into container")
		return err
	}

	if err := runSetupScripts(ctx, docker, container, cfg.SetupScripts, logger); err != nil {
		logger.Error().Err(err).Msg("failed to run setup scripts in container")
		return err
	}

	wsURL := "ws://" + endpoint.String() + "/ws"
	logger.Info().
		Str("url", wsURL).
		Str("container_id", endpoint.ContainerID).
		Msg("connecting to agent container")

	if err := api.RunClient(ctx, wsURL, logger); err != nil {
		logger.Error().Err(err).Msg("client disconnected with error")
	}

	/*
		logger.Info().Str("container_id", endpoint.ContainerID).Msg("stopping agent container")
		if stopErr := container.Stop(ctx, 10*time.Second); stopErr != nil {
			logger.Error().Err(stopErr).Str("container_id", endpoint.ContainerID).Msg("failed to stop container")
			return stopErr
		}
	*/

	return nil
}

// runSetupScripts copies each setup script declared in the agent config into
// the running container and executes it with bash. Scripts with Type "inline"
// use Source as the script contents; scripts with Type "file" read the contents
// from the host path in Source. Each script is written to a temporary path
// inside the container, made executable, and run before the agent starts.
func runSetupScripts(ctx context.Context, docker *runenv.DockerClient, container *runenv.Container, scripts []agentconfig.SetupScriptsConfig, logger zerolog.Logger) error {
	for i, s := range scripts {
		var content []byte

		switch s.Type {
		case agentconfig.SetupScriptTypeInline, "":
			content = []byte(s.Source)
		case agentconfig.SetupScriptTypeFile:
			if strings.HasPrefix(s.Source, "~") {
				home, err := os.UserHomeDir()
				if err != nil {
					logger.Error().Err(err).Str("source", s.Source).Msg("failed to get user home directory")
					return fmt.Errorf("get user home directory: %w", err)
				}
				s.Source = strings.Replace(s.Source, "~", home, 1)
			}
			info, err := os.Stat(s.Source)
			if err != nil {
				logger.Error().Err(err).Str("source", s.Source).Msg("failed to stat setup script source")
				return fmt.Errorf("stat setup script source %q: %w", s.Source, err)
			}
			if info.IsDir() {
				logger.Error().Str("source", s.Source).Msg("setup script source is a directory")
				return fmt.Errorf("setup script source %q is a directory, only files are supported", s.Source)
			}
			content, err = os.ReadFile(s.Source)
			if err != nil {
				logger.Error().Err(err).Str("source", s.Source).Msg("failed to read setup script source")
				return fmt.Errorf("read setup script source %q: %w", s.Source, err)
			}
		default:
			logger.Error().Str("type", string(s.Type)).Msg("unsupported setup script type")
			return fmt.Errorf("unsupported setup script type %q", s.Type)
		}

		target := fmt.Sprintf("/tmp/abyss-setup-%d.sh", i)
		if err := container.CopyFileFromHost(ctx, content, target, 0o755); err != nil {
			logger.Error().Err(err).Str("target", target).Msg("failed to copy setup script into container")
			return fmt.Errorf("copy setup script to %q: %w", target, err)
		}

		logger.Debug().Str("target", target).Str("type", string(s.Type)).Msg("executing setup script in container")
		stdout, stderr, err := container.ExecBash(ctx, fmt.Sprintf("chmod +x %s && bash %s", target, target))
		if err != nil {
			logger.Error().
				Err(err).
				Str("target", target).
				Str("stdout", stdout).
				Str("stderr", stderr).
				Msg("setup script failed")
			return fmt.Errorf("setup script %q failed: %w", target, err)
		}

		logger.Debug().Str("target", target).Str("stdout", stdout).Str("stderr", stderr).Msg("setup script completed")
	}

	return nil
}

// copyFiles copies each file declared in the agent config into the running
// container. Files with Type "inline" use Source as the file contents; files
// with Type "path" read the contents from the host path in Source. Each file
// is written to its Target path inside the container, creating parent
// directories as needed. When a "path" entry points to a directory, the
// directory is copied recursively into the container at Target, which is
// treated as the destination directory (its parent directories are created as
// needed).
func copyFiles(ctx context.Context, docker *runenv.DockerClient, container *runenv.Container, files []agentconfig.FileCopyConfig, logger zerolog.Logger) error {
	for _, f := range files {
		var content []byte
		var mode os.FileMode = 0o644

		switch f.Type {
		case agentconfig.FileCopyTypeInline:
			content = []byte(f.Source)
		case agentconfig.FileCopyTypePath:
			if strings.HasPrefix(f.Source, "~") {
				home, err := os.UserHomeDir()
				if err != nil {
					logger.Error().Err(err).Str("source", f.Source).Msg("failed to get user home directory")
					return fmt.Errorf("get user home directory: %w", err)
				}
				f.Source = strings.Replace(f.Source, "~", home, 1)
			}

			info, err := os.Stat(f.Source)
			if err != nil {
				logger.Error().Err(err).Str("source", f.Source).Msg("failed to stat copy file source")
				return fmt.Errorf("stat copy file source %q: %w", f.Source, err)
			}
			if info.IsDir() {
				// Recursively copy the directory into the container at Target.
				// CopyFromHost creates Target (and any parents) inside the
				// container before extracting, and preserves the directory
				// structure relative to Source.
				if err := container.CopyFromHost(ctx, container.ID(), f.Source, f.Target); err != nil {
					logger.Error().Err(err).Str("source", f.Source).Str("target", f.Target).Msg("failed to copy directory into container")
					return fmt.Errorf("copy directory %q to %q: %w", f.Source, f.Target, err)
				}
				logger.Debug().Str("target", f.Target).Str("type", string(f.Type)).Msg("copied directory into container")
				continue
			}
			mode = info.Mode()
			content, err = os.ReadFile(f.Source)
			if err != nil {
				logger.Error().Err(err).Str("source", f.Source).Msg("failed to read copy file source")
				return fmt.Errorf("read copy file source %q: %w", f.Source, err)
			}
		default:
			logger.Error().Str("type", string(f.Type)).Msg("unsupported copy file type")
			return fmt.Errorf("unsupported copy file type %q for target %q", f.Type, f.Target)
		}

		if err := container.CopyFileFromHost(ctx, content, f.Target, mode); err != nil {
			logger.Error().Err(err).Str("target", f.Target).Msg("failed to copy file into container")
			return fmt.Errorf("copy file to %q: %w", f.Target, err)
		}

		logger.Debug().Str("target", f.Target).Str("type", string(f.Type)).Msg("copied file into container")
	}

	return nil
}
