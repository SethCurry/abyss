package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/SethCurry/abyss/internal/acp/acptools"
	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/SethCurry/abyss/internal/api/agentapi"
	"github.com/SethCurry/abyss/internal/api/pacific"
	"github.com/SethCurry/abyss/internal/erres"
	"github.com/SethCurry/abyss/internal/runenv"
	"github.com/SethCurry/abyss/internal/timber"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	logFile, err := timber.OpenLogFile()
	if err != nil {
		panic(err)
	}

	defer func() {
		defErr := logFile.Close()
		if defErr != nil {
			log.Error().Err(err).Msg("failed to close log file")
		}
	}()

	logOut := io.MultiWriter(logFile, os.Stderr)
	logger := zerolog.New(logOut).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	log.Logger = logger

	cmd := &cli.Command{
		Name:        "abyss",
		Usage:       "Agent Runtime Environment(s)",
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

			_ = timber.CleanLogDir(10)
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:        "client",
				Aliases:     []string{"c"},
				Usage:       "Starts the host-side proxy that your ACP client will connect directly to.",
				Description: "Starts the host-side proxy that your ACP client will connect to via stdio.  It will create a Docker container, start the container-side proxy inside of it, and proxy your ACP connection into the container.",
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
					return runClient(ctx, "", agentCfg, logger)
				},
			},
			{
				Name:    "oneshot",
				Aliases: []string{"p"},
				Usage:   "Executes a single agent turn based on the provided prompt.",
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
						logger.Error().
							Err(err).
							Str("config_path", configPath).
							Msg("failed to load agent config")
						return err
					}
					logger.Debug().
						Str("config_path", configPath).
						Str("image", agentCfg.Docker.Image).
						Strs("agent_command", agentCfg.Docker.AgentCommand).
						Msg("starting Docker agent")
					prompt := strings.Join(cmd.Args().Slice(), " ")
					return runClient(ctx, prompt, agentCfg, logger)
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
						Usage:   "The address to run the HTTP server on, formatted like 127.0.0.1:8080",
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
					&cli.StringFlag{
						Name:  "tls-cert",
						Usage: "Path to the TLS server certificate.",
					},
					&cli.StringFlag{
						Name:  "tls-key",
						Usage: "Path to the TLS server key.",
					},
					&cli.StringFlag{
						Name:  "tls-ca",
						Usage: "Path to the CA certificate used to verify client certificates.",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					logger.Info().Str("path", agentconfig.DefaultStartFilePath).Msg("waiting for file to appear to start")
					for {
						if _, err := os.Stat(agentconfig.DefaultStartFilePath); err == nil {
							logger.Info().Str("path", agentconfig.DefaultStartFilePath).Msg("found start file, starting")
							break
						}
						logger.Debug().Str("path", agentconfig.DefaultStartFilePath).Msg("start file not found, waiting")
						time.Sleep(agentconfig.WaitForStartFileSleepDuration)
					}
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

					tlsCert := cmd.String("tls-cert")
					tlsKey := cmd.String("tls-key")
					tlsCA := cmd.String("tls-ca")

					if tlsCert != "" && tlsKey != "" && tlsCA != "" {
						tlsConfig, err := pacific.LoadServerTLSConfig(tlsCert, tlsKey, tlsCA)
						if err != nil {
							return fmt.Errorf("load TLS config: %w", err)
						}
						return httpSrv.ServeTLS(cmd.String("addr"), tlsConfig)
					}

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
					{
						Name:  "gc",
						Usage: "Garbage collect all abyss containers.",
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
								logger.Info().Str("container_id", v.ID).Strs("container_names", v.Names).Msg("stopping container")
								cont := docker.GetContainer(v.ID)
								err = cont.Stop(ctx, time.Second*5)
								if err != nil {
									return fmt.Errorf("failed to stop container %q: %w", v.ID, err)
								}
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
			logger.Error().Err(err).Str("human_error", humanErr.HumanError()).Msg("command failed")
		} else {
			logger.Error().Err(err).Msg("command failed")
		}
	}
}

// installTLSCerts returns a build step that copies the server certificate,
// server key, and CA certificate into the container so the server can serve
// mutual TLS.
func installTLSCerts(certs *pacific.Certificates) runenv.ContainerBuildStep {
	return func(ctx context.Context, container *runenv.Container) error {
		files := []struct {
			path    string
			content []byte
		}{
			{path: agentconfig.DefaultTLSServerCertPath, content: certs.ServerCertPEM},
			{path: agentconfig.DefaultTLSServerKeyPath, content: certs.ServerKeyPEM},
			{path: agentconfig.DefaultTLSCACertPath, content: certs.CACertPEM},
		}

		for _, f := range files {
			if err := container.CopyFileFromHost(ctx, f.content, f.path, 0o600); err != nil {
				return fmt.Errorf("copy %q into container: %w", f.path, err)
			}
		}

		return nil
	}
}
