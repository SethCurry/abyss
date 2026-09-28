// Package main is the entry point for the abyss binary.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/SethCurry/abyss/internal/acptools"
	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/SethCurry/abyss/internal/api/agentapi"
	"github.com/SethCurry/abyss/internal/api/pacific"
	"github.com/SethCurry/abyss/internal/constants"
	"github.com/SethCurry/abyss/internal/plugin"
	"github.com/SethCurry/abyss/internal/runenv"
	"github.com/SethCurry/abyss/internal/timber"
	"github.com/SethCurry/abyss/internal/types"
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

	//nolint:lll
	cmd := &cli.Command{
		Name:  "abyss",
		Usage: "Run your AI agent in its own private sandbox.",
		Description: `Abyss gives your AI agent a tidy home of its own: a Docker container, which
is a sealed-off workspace where the agent can think, build, and experiment
without touching anything else on your computer.

A small YAML file drives everything: it says which agent to run and which
folders on your computer the agent is allowed to see. Point your editor at
abyss, and it quietly does the behind-the-scenes work for you: creating the
container, sharing in those folders, running any setup steps, and carrying
the conversation between your editor and the agent over an encrypted
connection.

The commands you'll use most:

  client    connect your editor to an agent
  oneshot   ask an agent one question and print the answer
  docker    see and tidy up the containers abyss has created

New here? The guides at https://abyss.scurry.io/ walk you through everything.`,
		Version: constants.Version,
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

			if err := timber.CleanLogDir(10); err != nil {
				log.Logger.Warn().Err(err).Msg("failed to clean log directory")
			}

			log.Logger.Info().Str("version", constants.Version).Msg("starting abyss")
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Connect your editor to an agent running inside a container.",
				Description: `This is the command your editor runs for you, so once things are set up you
won't need to type it yourself — your editor starts it whenever you chat with
your agent.

When it starts, abyss creates a brand-new container for the agent, sets it
up exactly the way your configuration file describes, and then bridges your
editor's connection to the agent living inside. When the chat ends, the
container is shut down again.

For example, to start the agent described by my-agent.yaml:

  abyss client -f my-agent.yaml`,
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
					return runClient(ctx, "", configPath, agentCfg, log.Logger)
				},
			},
			{
				Name:    "oneshot",
				Aliases: []string{"p"},
				Usage:   "Ask an agent one question and print the answer.",
				Description: `A quick way to try out a configuration: oneshot spins up an agent's container
just like client does, sends it a single prompt, prints the answer (and the
logs) right in your terminal, and then packs the container away again.

For example:

  abyss oneshot -f my-agent.yaml "What is the capital of France?"

Because it shows you everything that's going on, oneshot is also the first
tool to reach for when a configuration isn't behaving.`,
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
					return runClient(ctx, prompt, configPath, agentCfg, log.Logger)
				},
			},
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Run abyss's agent-side half inside the container (abyss starts this for you).",
				Description: `This is the piece of abyss that lives inside the agent's container. It waits
for abyss to give it the go-ahead, then brings your agent to life and
carries the conversation between the agent and your editor.

You should never need to run this command yourself: whenever abyss creates a
container, it starts the server inside automatically. If you've spotted it,
it's probably because you were peeking at the programs running inside the
container, and everything is working exactly as it should.`,
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

					plugins, err := plugin.NewACPManager(ctx)
					if err != nil {
						return err
					}

					httpSrv := agentapi.NewServer(agentCmd, localTerminal, localFilesystem, plugins)

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
							return types.NewHumanError(fmt.Errorf("failed to load TLS config: %w", err), "Failed to load TLS config.  Ensure that the files exist and that you have permissions to read them.")
						}
						return httpSrv.ServeTLS(cmd.String("addr"), tlsConfig)
					}

					return httpSrv.Serve(cmd.String("addr"))
				},
			},
			{
				Name:    "docker",
				Aliases: []string{"d"},
				Usage:   "See and tidy up the containers abyss has created.",
				Description: `Every time an agent session starts, abyss creates a Docker container for
the agent to live in. These commands help you peek at the ones that are
running and send any stragglers home when you're done.`,
				Commands: []*cli.Command{
					{
						Name:    "ps",
						Aliases: []string{"p"},
						Usage:   "List the containers abyss is currently running.",
						Description: `Shows every abyss container that is still up and running, along with each
one's ID and name. It's handy for a peek behind the scenes, or for finding a
container's ID — Docker's own commands (like docker logs) will ask for it.`,
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
						Usage: "Stop every container that abyss is running.",
						Description: `Tells Docker to stop and remove all containers
Abyss has created (including those currently running).

Only containers that abyss itself created are stopped; everything else on
your computer is left completely alone.`,
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
								err = cont.Stop(ctx, time.Second*5, true)
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
		if humanErr, ok := errors.AsType[types.HumanError](err); ok {
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
	files := []struct {
		path    string
		content []byte
	}{
		{path: agentconfig.DefaultTLSServerCertPath, content: certs.ServerCertPEM},
		{path: agentconfig.DefaultTLSServerKeyPath, content: certs.ServerKeyPEM},
		{path: agentconfig.DefaultTLSCACertPath, content: certs.CACertPEM},
	}

	steps := make([]runenv.ContainerBuildStep, 0, len(files))
	for _, f := range files {
		steps = append(steps, func(ctx context.Context, container *runenv.Container) error {
			if err := container.CopyFileFromHost(ctx, f.content, f.path, 0o600); err != nil {
				return fmt.Errorf("copy %q into container: %w", f.path, err)
			}
			return nil
		})
	}

	return runenv.NewParallelStep(steps...)
}
