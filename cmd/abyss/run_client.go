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

	var certs *pacific.Certificates
	var cont *runenv.Container
	var endpoint runenv.ContainerEndpoint

	if cfg.Docker.PersistentName != "" {
		cont, certs, endpoint, err = startPersistentContainer(ctx, docker, configPath, cfg, image, logger)
	} else {
		cont, certs, endpoint, err = createAgentContainer(ctx, docker, configPath, cfg, image, logger)
	}
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

	if err := api.RunClient(ctx, cfg, prompt, wsURL, tlsConfig, logger); err != nil {
		logger.Error().Err(err).Msg("client disconnected with error")
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

// startPersistentContainer re-uses the container named by the config's
// persistent_name, creating it if it doesn't exist yet.
//
// When an existing container is found, the TLS certificates that were
// installed when it was created are pulled back out of the container so this
// session can authenticate against them, and the container is started if it
// isn't already running. Otherwise a new container is created with the
// persistent name.
func startPersistentContainer(
	ctx context.Context,
	docker *runenv.DockerClient,
	configPath string,
	cfg *agentconfig.AgentConfig,
	image string,
	logger zerolog.Logger,
) (*runenv.Container, *pacific.Certificates, runenv.ContainerEndpoint, error) {
	name := cfg.Docker.PersistentName

	inspect, err := docker.FindContainerByName(ctx, name)
	if err != nil {
		logger.Error().Err(err).Str("name", name).Msg("failed to look up persistent container")
		return nil, nil, runenv.ContainerEndpoint{}, types.NewHumanError(err,
			fmt.Sprintf("Failed to look up the persistent container %q.", name),
			"Have you made sure Docker is running and that you have permission to connect?")
	}

	if inspect == nil {
		logger.Info().Str("name", name).Msg("no existing persistent container found, creating one")
		return createAgentContainer(ctx, docker, configPath, cfg, image, logger)
	}

	return reusePersistentContainer(ctx, docker, inspect, cfg, image, logger)
}

// reusePersistentContainer reconnects to an existing persistent container: it
// pulls the TLS certificates from the container's filesystem so the client can
// authenticate against the server certificates it already holds, starts the
// container if it isn't running, and works out how to reach it from the host.
func reusePersistentContainer(
	ctx context.Context,
	docker *runenv.DockerClient,
	inspect *container.InspectResponse,
	cfg *agentconfig.AgentConfig,
	image string,
	logger zerolog.Logger,
) (*runenv.Container, *pacific.Certificates, runenv.ContainerEndpoint, error) {
	name := cfg.Docker.PersistentName

	// Only re-use containers that abyss created; starting a foreign container
	// or pulling files from it would be surprising.
	if inspect.Config == nil || inspect.Config.Labels["abyss_version"] == "" {
		err := fmt.Errorf("container %q exists but is not managed by abyss", name)
		logger.Error().Err(err).Msg("refusing to re-use container")
		return nil, nil, runenv.ContainerEndpoint{}, types.NewHumanError(err,
			fmt.Sprintf("The container %q already exists, but abyss didn't create it.", name),
			"Remove that container or pick a different persistent_name in your agent configuration.")
	}

	if inspect.Config.Image != image {
		logger.Warn().
			Str("container_image", inspect.Config.Image).
			Str("configured_image", image).
			Msg("persistent container was created from a different image; it will keep using the old one")
	}

	cont := docker.GetContainer(inspect.ID)

	// The certificates inside the container were generated when it was
	// created; re-use them so the client can still authenticate.
	var certs *pacific.Certificates
	if !cfg.Websocket.DisableTLS {
		var err error
		certs, err = pullTLSCerts(ctx, cont)
		if err != nil {
			logger.Error().Err(err).
				Str("container_id", inspect.ID).
				Msg("failed to pull TLS certificates from persistent container")
			return nil, nil, runenv.ContainerEndpoint{}, types.NewHumanError(err,
				fmt.Sprintf("Failed to read the TLS certificates from the container %q.", name),
				"It may have been created by an older version of abyss.  Remove the container so a new one can be created.")
		}
	}

	if inspect.State == nil || !inspect.State.Running {
		logger.Info().Str("container_id", inspect.ID).Msg("starting persistent agent container")
		if err := cont.Start(ctx); err != nil {
			return nil, nil, runenv.ContainerEndpoint{}, err
		}
	} else {
		logger.Info().Str("container_id", inspect.ID).Msg("re-using running persistent agent container")
	}

	endpoint, err := docker.EndpointFor(inspect, agentconfig.DefaultServerPort)
	if err != nil {
		return nil, nil, runenv.ContainerEndpoint{}, err
	}

	return cont, certs, endpoint, nil
}

// pullTLSCerts reads the CA certificate and client certificate pair that were
// stored inside a persistent container when it was created.
func pullTLSCerts(ctx context.Context, cont *runenv.Container) (*pacific.Certificates, error) {
	caCert, err := cont.ReadFile(ctx, agentconfig.DefaultTLSCACertPath)
	if err != nil {
		return nil, fmt.Errorf("read CA certificate: %w", err)
	}

	clientCert, err := cont.ReadFile(ctx, agentconfig.DefaultTLSClientCertPath)
	if err != nil {
		return nil, fmt.Errorf("read client certificate: %w", err)
	}

	clientKey, err := cont.ReadFile(ctx, agentconfig.DefaultTLSClientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read client key: %w", err)
	}

	return &pacific.Certificates{
		CACertPEM:     caCert,
		ClientCertPEM: clientCert,
		ClientKeyPEM:  clientKey,
	}, nil
}

// createAgentContainer pulls the image, generates a fresh certificate set,
// and creates a container from the agent config. The container is named after
// the config's persistent_name when one is set.
func createAgentContainer(
	ctx context.Context,
	docker *runenv.DockerClient,
	configPath string,
	cfg *agentconfig.AgentConfig,
	image string,
	logger zerolog.Logger,
) (*runenv.Container, *pacific.Certificates, runenv.ContainerEndpoint, error) {
	err := runenv.PullImage(ctx, docker.Client, image, cfg.Docker.ImagePullPolicy)
	if err != nil {
		logger.Error().Err(err).Str("image", image).Msg("failed to pull Docker image")
		return nil, nil, runenv.ContainerEndpoint{}, types.NewHumanError(
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
			return nil, nil, runenv.ContainerEndpoint{}, err
		}
	}

	cont, endpoint, err := startAgentContainer(ctx, docker, configPath, cfg, image, certs, logger)
	if err != nil {
		return nil, nil, runenv.ContainerEndpoint{}, err
	}

	return cont, certs, endpoint, nil
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

	// A persistent name keeps the same container across sessions; without
	// one, Docker generates a throwaway name.
	config.Name = cfg.Docker.PersistentName
	config.Persistent = cfg.Docker.PersistentName != ""

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
