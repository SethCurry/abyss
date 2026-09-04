package runenv

import (
	"context"

	"github.com/SethCurry/abyss/internal/types"
	"github.com/moby/moby/api/types/container"
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
