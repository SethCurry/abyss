package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/SethCurry/abyss/internal/api/pacific"
	"github.com/SethCurry/abyss/internal/runenv"
	"github.com/SethCurry/abyss/internal/types"
	api "github.com/SethCurry/abyss/internal/websockets/wsacp"
	"github.com/moby/moby/api/types/container"
	"github.com/rs/zerolog"
)

// runClient starts a Docker container running the abyss server command and
// bridges it to a client over stdio.
//
// If prompt is not empty, it will be used as a one-shot prompt to the server.
func runClient(
	ctx context.Context, prompt string, configPath string, cfg *agentconfig.AgentConfig, logger zerolog.Logger,
) error {
	docker, err := runenv.NewDockerClient()
	if err != nil {
		return types.NewHumanError(err,
			"Failed to connect to Docker.",
			"Have you made sure Docker is running and that you have permission to connect?")
	}
	defer func() {
		closeErr := docker.Close()
		if closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed to close Docker client")
		}
	}()

	image := cfg.Docker.Image
	if image == "" {
		image = agentconfig.DefaultImage
	}

	err = runenv.PullImage(ctx, docker.Client, image, cfg.Docker.ImagePullPolicy)
	if err != nil {
		logger.Error().Err(err).Str("image", image).Msg("failed to pull Docker image")
		return types.NewHumanError(
			fmt.Errorf("failed to pull Docker image: %w", err),
			fmt.Sprintf("Failed to pull Docker image %q.  Ensure that the image exists and that you have permission to pull it.",
				image))
	}

	// Generate a certificate set for mutual TLS unless the user disabled it.
	var certs *pacific.Certificates
	if !cfg.Websocket.DisableTLS {
		certs, err = pacific.GenerateCertificates()
		if err != nil {
			logger.Error().Err(err).Msg("failed to generate TLS certificates")
			return err
		}
	}

	cont, endpoint, err := startAgentContainer(ctx, docker, configPath, cfg, image, certs, logger)
	if err != nil {
		return err
	}

	err = cont.CreateAgentStartFile(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create agent start file")
		return err
	}

	// Sleep for as long as the agent does between checking for the file
	// Prevents a race-condition where the start file was created but the
	// container-side proxy hasn't seen it yet.
	// TODO replace this with a dial retry loop
	time.Sleep(agentconfig.WaitForStartFileSleepDuration)

	scheme := "ws"
	var tlsConfig *tls.Config
	if certs != nil {
		scheme = "wss"
		tlsConfig, err = certs.ClientTLSConfig()
		if err != nil {
			logger.Error().Err(err).Msg("failed to build client TLS config")
			return err
		}
	}

	wsURL := scheme + "://" + endpoint.String() + "/ws"
	logger.Info().
		Str("url", wsURL).
		Str("container_id", endpoint.ContainerID).
		Msg("connecting to agent container")

	// TODO clean this up, there's no need to have an if here
	if prompt == "" {
		if err := api.RunClient(ctx, wsURL, tlsConfig, logger); err != nil {
			logger.Error().Err(err).Msg("client disconnected with error")
		}
	} else {
		if err := api.Oneshot(ctx, prompt, wsURL, tlsConfig, logger); err != nil {
			logger.Error().Err(err).Msg("oneshot failed")
		}
	}

	logger.Info().Str("container_id", endpoint.ContainerID).Msg("stopping agent container")
	if stopErr := cont.Stop(ctx, 10*time.Second, false); stopErr != nil {
		logger.Error().
			Err(stopErr).
			Str("container_id", endpoint.ContainerID).
			Msg("failed to stop container")
		return stopErr
	}

	return nil
}

// startAgentContainer builds and starts the agent container, so runClient can
// focus on wiring up the client connection rather than container assembly.
func startAgentContainer(
	ctx context.Context,
	docker *runenv.DockerClient,
	configPath string,
	cfg *agentconfig.AgentConfig,
	image string,
	certs *pacific.Certificates,
	logger zerolog.Logger,
) (*runenv.Container, runenv.ContainerEndpoint, error) {
	config := &runenv.ContainerConfig{
		Config: &container.Config{
			Entrypoint: []string{"/bin/bash"},
			Cmd:        []string{"-c", agentconfig.BuildAgentProxyArgs(cfg)},
		},
		ContainerPort: agentconfig.DefaultServerPort,
	}

	builder, err := runenv.NewContainerBuilder(
		configPath,
		config,
		runenv.WithImage(image),
		runenv.WithHostMounts(&cfg.Docker),
		runenv.WithExposeContainerPort(8080),
	)
	if err != nil {
		logger.Error().Err(err).Msg("failed to build container config")
		return nil, runenv.ContainerEndpoint{}, types.NewHumanError(err,
			"Failed to build container config.",
			"This is likely an issue with abyss itself or the Docker image you are using.",
			"Please report a bug if you have time.")
	}

	builder.AddSteps(
		runenv.WithCopyFiles(cfg.CopyFiles),
		runenv.WithSetupScripts(cfg.SetupScripts),
	)

	if certs != nil {
		builder.AddStep(installTLSCerts(certs))
	}

	cont, endpoint, err := builder.Build(ctx, docker)
	if err != nil {
		logger.Error().Err(err).Msg("failed to start container")
		return nil, runenv.ContainerEndpoint{}, err
	}

	return cont, endpoint, nil
}
