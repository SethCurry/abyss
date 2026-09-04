package runenv

import (
	"context"
	"fmt"

	"github.com/SethCurry/abyss/internal/types"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
)

type ContainerConfig struct {
	Host     *container.HostConfig
	Config   *container.Config
	Endpoint *ContainerEndpoint
}

func (c *ContainerConfig) Validate() *types.ValidationError {
	if c.Host == nil {
		return &types.ValidationError{Type: "missing", Field: "Host", Reason: "HostConfig is required"}
	}
	if c.Config == nil {
		return &types.ValidationError{Type: "missing", Field: "Config", Reason: "ContainerConfig is required"}
	}

	return nil
}

func WithImage(imgName string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		config.Config.Image = imgName
		return nil
	}
}

func WithEnv(envString string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		config.Config.Env = append(config.Config.Env, envString)

		return nil
	}
}

func WithExposePort(port network.Port) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Config.ExposedPorts == nil {
			config.Config.ExposedPorts = make(map[network.Port]struct{})
		}
		config.Config.ExposedPorts[port] = struct{}{}
		return nil
	}
}

func WithLabel(label, value string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Config.Labels == nil {
			config.Config.Labels = make(map[string]string)
		}
		config.Config.Labels[label] = value
		return nil
	}
}

func WithHostPort(port network.Port) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Host.PortBindings == nil {
			config.Host.PortBindings = make(map[network.Port][]network.PortBinding)
		}

		config.Host.PortBindings[port] = make([]network.PortBinding, 0, 1)

		return nil
	}
}

func WithHostBind(from string, to string) ContainerPreBuildStep {
	return func(config *ContainerConfig) error {
		if config.Host.Binds == nil {
			config.Host.Binds = make([]string, 0, 1)
		}
		config.Host.Binds = append(config.Host.Binds, fmt.Sprintf("%s:%s", from, to))
		return nil
	}
}

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

type ContainerBuilder struct {
	config *ContainerConfig
	steps  []ContainerBuildStep
}

func (b *ContainerBuilder) Build(ctx context.Context, client *DockerClient) (*Container, ContainerEndpoint, error) {
	container, endpoint, err := client.StartContainer(ctx, b.config.Config.Image, b.config.Config, b.config.Host, "name", 0, 0)
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
