package main

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/SethCurry/abyss/internal/api/pacific"
	"github.com/SethCurry/abyss/internal/erres"
	"github.com/SethCurry/abyss/internal/runenv"
	api "github.com/SethCurry/abyss/internal/websockets/wsacp"
	"github.com/moby/moby/api/types/container"
	"github.com/rs/zerolog"
)

// runClient starts a Docker container running the abyss server command and
// bridges it to a client over stdio.
//
// If prompt is not empty, it will be used as a one-shot prompt to the server.
func runClient(ctx context.Context, prompt string, cfg *agentconfig.AgentConfig, logger zerolog.Logger) error {
	docker, err := runenv.NewDockerClient()
	if err != nil {
		return erres.NewHumanError("Failed to connect to Docker.\nHave you made sure Docker is running and that you have permission to connect?", err)
	}
	defer func() {
		closeErr := docker.Close()
		if closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed to close Docker client")
		}
	}()

	image := cfg.Docker.Image
	if image == "" {
		image = defaultImage
	}
	agent := cfg.Docker.AgentCommand

	// Generate a certificate set for mutual TLS unless the user disabled it.
	var certs *pacific.Certificates
	if !cfg.Websocket.DisableTLS {
		certs, err = pacific.GenerateCertificates()
		if err != nil {
			logger.Error().Err(err).Msg("failed to generate TLS certificates")
			return err
		}
	}

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

	if certs != nil {
		agentArgs = append(agentArgs,
			"--tls-cert", tlsServerCertPath,
			"--tls-key", tlsServerKeyPath,
			"--tls-ca", tlsCACertPath,
		)
	}

	joinedArgs := strings.Join(agentArgs, " ")

	config := &runenv.ContainerConfig{
		Config: &container.Config{
			Entrypoint: []string{"/bin/bash", "-c"},
			Cmd:        []string{"/usr/local/bin/abyss server " + joinedArgs},
		},
		Host:          docker.ApplyHostMounts(&cfg.Docker, nil),
		ContainerPort: serverPort,
	}

	builder, err := runenv.NewContainerBuilder(config, runenv.WithImage(image))
	if err != nil {
		logger.Error().Err(err).Msg("failed to build container config")
		return err
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
		return err
	}

	time.Sleep(time.Second * 1)

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
	if stopErr := cont.Stop(ctx, 10*time.Second); stopErr != nil {
		logger.Error().Err(stopErr).Str("container_id", endpoint.ContainerID).Msg("failed to stop container")
		return stopErr
	}

	return nil
}
