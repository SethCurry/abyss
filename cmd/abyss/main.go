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

	logOut := zerolog.ConsoleWriter{Out: io.MultiWriter(logFile, os.Stderr)}
	globalLogger := zerolog.New(logOut).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	log.Logger = globalLogger

	globalLogger, closeLogger := timber.CreateLogger(zerolog.DebugLevel)
	defer closeLogger()

	cmd := &cli.Command{
		Name:        "abyss",
		Usage:       "A tool for managing and connecting to agents running in containers.",
		Description: "abyss creates Docker containers for you, copies files, creates bind mounts, executes setup scripts, and proxies your ACP connection with mutual TLS authentication.",
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
			log.Logger = globalLogger.Level(level)

			_ = timber.CleanLogDir(10)
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:        "client",
				Aliases:     []string{"c"},
				Usage:       "Starts the host-side proxy that your editor connects to.",
				Description: "Creates a Docker container, starts the container-side proxy inside of it, and proxies your ACP connection into the container.",
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
						log.Logger.Error().Err(err).Str("config_path", configPath).Msg("failed to load agent config")
						return err
					} else {
						log.Logger.Debug().
							Str("config_path", configPath).
							Msg("loaded config")
					}
					return runClient(ctx, "", agentCfg, log.Logger)
				},
			},
			{
				Name:    "oneshot",
				Aliases: []string{"p"},
				Usage:   "Executes a single agent turn, batch-style.",
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
						log.Logger.Error().
							Err(err).
							Str("config_path", configPath).
							Msg("failed to load agent config")
						return err
					}
					log.Logger.Debug().
						Str("config_path", configPath).
						Msg("loaded config")
					prompt := strings.Join(cmd.Args().Slice(), " ")
					return runClient(ctx, prompt, agentCfg, log.Logger)
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
						Name:  "local-terminal",
						Usage: "Run terminal ACP commands on the client rather than this server",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "local-filesystem",
						Usage: "Run filesystem ACP commands on the client rather than this server",
						Value: false,
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
					log.Logger.Info().Str("start_file_path", agentconfig.DefaultStartFilePath).Msg("waiting for start file to appear to start")
					for {
						if _, err := os.Stat(agentconfig.DefaultStartFilePath); err == nil {
							log.Logger.Info().
								Str("path", agentconfig.DefaultStartFilePath).
								Msg("found start file, starting")
							break
						}
						log.Logger.Debug().
							Str("path", agentconfig.DefaultStartFilePath).
							Msg("start file not found, waiting")
						time.Sleep(agentconfig.WaitForStartFileSleepDuration)
					}

					agentCmd := cmd.StringSlice("agent")

					var localTerminal *acptools.TerminalTools
					var localFilesystem *acptools.FilesystemTools

					if cmd.Bool("local-terminal") {
						log.Logger.Debug().Msg("local-terminal enabled")
						localTerminal = acptools.NewTerminalTools(log.Logger)
					}
					if cmd.Bool("local-filesystem") {
						log.Logger.Debug().Msg("local-filesystem enabled")
						localFilesystem = acptools.NewFilesystemTools(log.Logger)
					}

					httpSrv := agentapi.NewServer(agentCmd, localTerminal, localFilesystem)

					tlsCert := cmd.String("tls-cert")
					tlsKey := cmd.String("tls-key")
					tlsCA := cmd.String("tls-ca")

					if tlsCert != "" && tlsKey != "" && tlsCA != "" {
						log.Logger.Info().
							Str("tls-cert", tlsCert).
							Str("tls-key", tlsKey).
							Str("tls-ca", tlsCA).
							Msg("TLS enabled")
						tlsConfig, err := pacific.LoadServerTLSConfig(
							tlsCert,
							tlsKey,
							tlsCA)
						if err != nil {
							return erres.NewHumanError(fmt.Errorf("failed to load TLS config: %w", err), "Failed to load TLS config.  Ensure that the files exist and that you have permissions to read them.")
						}
						return httpSrv.ServeTLS(cmd.String("addr"), tlsConfig)
					}

					return httpSrv.Serve(cmd.String("addr"))
				},
			},
			{
				Name:    "docker",
				Aliases: []string{"d"},
				Usage:   "Clean up old Docker containers, see running abyss containers, etc.",
				Commands: []*cli.Command{
					{
						Name:        "ps",
						Aliases:     []string{"p"},
						Usage:       "List Abyss containers that are currently running.",
						Description: "Finds all running containers with the `abyss` label.",
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
						Name:        "gc",
						Usage:       "Stop all abyss containers.",
						Description: "Stops all containers with the `abyss` label.",
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
								log.Logger.Info().
									Str("container_id", v.ID).
									Strs("container_names", v.Names).
									Msg("stopping container")
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
			log.Logger.Error().Err(err).Str("human_error", humanErr.HumanError()).Msg("command failed")
			fmt.Fprintf(os.Stderr, "%s\n", humanErr.HumanError())
		} else {
			log.Logger.Error().Err(err).Msg("command failed")
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
