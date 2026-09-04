package runenv

import (
	"context"
	"fmt"

	"github.com/SethCurry/abyss/internal/types"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// ContainerConfig holds the Docker host and container configuration used to
// build a container.
type ContainerConfig struct {
	Host     *container.HostConfig
	Config   *container.Config
	Endpoint *ContainerEndpoint
}

// Validate ensures the required host and container configs are present.
func (c *ContainerConfig) Validate() *types.ValidationError {
	if c.Host == nil {
		return &types.ValidationError{Type: "missing", Field: "Host", Reason: "HostConfig is required"}
	}
	if c.Config == nil {
		return &types.ValidationError{Type: "missing", Field: "Config", Reason: "ContainerConfig is required"}
	}

	return nil
}

// WithImage returns a pre-build step that sets the container image.
func WithImage(imgName string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		config.Config.Image = imgName
		return nil
	}
}

// WithEnv returns a pre-build step that appends an environment variable.
func WithEnv(envString string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		config.Config.Env = append(config.Config.Env, envString)

		return nil
	}
}

// WithExposePort returns a pre-build step that marks a port as exposed.
func WithExposePort(port network.Port) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Config.ExposedPorts == nil {
			config.Config.ExposedPorts = make(map[network.Port]struct{})
		}
		config.Config.ExposedPorts[port] = struct{}{}
		return nil
	}
}

// WithLabel returns a pre-build step that sets a container label.
func WithLabel(label, value string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Config.Labels == nil {
			config.Config.Labels = make(map[string]string)
		}
		config.Config.Labels[label] = value
		return nil
	}
}

// WithHostPort returns a pre-build step that reserves a host port binding.
func WithHostPort(port network.Port) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Host.PortBindings == nil {
			config.Host.PortBindings = make(map[network.Port][]network.PortBinding)
		}

		config.Host.PortBindings[port] = make([]network.PortBinding, 0, 1)

		return nil
	}
}

// WithHostBind returns a pre-build step that binds a host path into the container.
func WithHostBind(from string, to string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		cleanedFrom, err := cleanPath(from)
		if err != nil {
			return fmt.Errorf("failed to clean from path %q: %w", from, err)
		}

		cleanedTo, err := cleanPath(to)
		if err != nil {
			return fmt.Errorf("failed to clean to path %q: %w", to, err)
		}
		if config.Host.Binds == nil {
			config.Host.Binds = make([]string, 0, 1)
		}
		config.Host.Binds = append(config.Host.Binds, fmt.Sprintf("%s:%s", cleanedFrom, cleanedTo))
		return nil
	}
}

// NewContainerBuilder creates a builder, applying any pre-build steps to the config.
func NewContainerBuilder(config *ContainerConfig, steps ...ContainerPreBuildStep) (*ContainerBuilder, error) {
	if config == nil {
		config = &ContainerConfig{
			Host:   &container.HostConfig{},
			Config: &container.Config{},
		}
	}

	if config.Host == nil {
		config.Host = &container.HostConfig{}
	}
	if config.Config == nil {
		config.Config = &container.Config{}
	}

	for _, step := range steps {
		if err := step(config); err != nil {
			return nil, err
		}
	}

	return &ContainerBuilder{config: config}, nil
}

// ContainerBuilder accumulates build steps and starts a container from its config.
type ContainerBuilder struct {
	config *ContainerConfig
	steps  []ContainerBuildStep
}

// AddStep appends a single build step to the builder.
func (b *ContainerBuilder) AddStep(step ContainerBuildStep) {
	b.steps = append(b.steps, step)
}

// AddSteps appends multiple build steps to the builder.
func (b *ContainerBuilder) AddSteps(newSteps ...ContainerBuildStep) {
	b.steps = append(b.steps, newSteps...)
}

// Build starts the container and runs each build step against it.
func (b *ContainerBuilder) Build(ctx context.Context, cli *DockerClient) (*Container, ContainerEndpoint, error) {
	container, endpoint, err := cli.StartContainer(ctx, b.config.Config.Image, b.config.Config, b.config.Host, "name", 0, 0)
	if err != nil {
		return nil, endpoint, err
	}

	_, err = cli.client.ContainerStart(ctx, container.ID(), client.ContainerStartOptions{})
	if err != nil {
		return nil, endpoint, err
	}

	for _, step := range b.steps {
		if err := step(container); err != nil {
			return nil, endpoint, err
		}
	}

	return container, endpoint, nil
}

type ContainerPreBuildStep func(*ContainerConfig) error
type ContainerBuildStep func(*Container) error
