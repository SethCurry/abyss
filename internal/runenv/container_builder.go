package runenv

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/SethCurry/abyss/internal/agentconfig"
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
		if err := step(ctx, container); err != nil {
			return nil, endpoint, err
		}
	}

	return container, endpoint, nil
}

type ContainerPreBuildStep func(*ContainerConfig) error
type ContainerBuildStep func(context.Context, *Container) error

// WithSetupScripts returns a build step that copies each setup script into the
// container and executes it with bash. Scripts with Type "inline" use Source as
// the script contents; scripts with Type "file" read the contents from the host
// path in Source. Each script is written to a temporary path inside the
// container, made executable, and run before the agent starts.
func WithSetupScripts(scripts []agentconfig.SetupScriptsConfig) ContainerBuildStep {
	return func(ctx context.Context, container *Container) error {
		for i, s := range scripts {
			content, err := setupScriptContent(s)
			if err != nil {
				return err
			}

			target := fmt.Sprintf("/tmp/abyss-setup-%d.sh", i)
			if err := container.CopyFileFromHost(ctx, content, target, 0o755); err != nil {
				return fmt.Errorf("copy setup script to %q: %w", target, err)
			}

			stdout, stderr, err := container.ExecBash(ctx, fmt.Sprintf("chmod +x %s && bash %s", target, target))
			if err != nil {
				container.logger.Error().
					Err(err).
					Str("target", target).
					Str("stdout", stdout).
					Str("stderr", stderr).
					Msg("setup script failed")
				return fmt.Errorf("setup script %q failed: %w", target, err)
			}

			container.logger.Debug().Str("target", target).Str("stdout", stdout).Str("stderr", stderr).Msg("setup script completed")
		}

		return nil
	}
}

// setupScriptContent resolves the contents of a setup script based on its type.
func setupScriptContent(s agentconfig.SetupScriptsConfig) ([]byte, error) {
	switch s.Type {
	case agentconfig.SetupScriptTypeInline, "":
		return []byte(s.Source), nil
	case agentconfig.SetupScriptTypeFile:
		source := s.Source
		if strings.HasPrefix(source, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("get user home directory: %w", err)
			}
			source = strings.Replace(source, "~", home, 1)
		}

		info, err := os.Stat(source)
		if err != nil {
			return nil, fmt.Errorf("stat setup script source %q: %w", source, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("setup script source %q is a directory, only files are supported", source)
		}

		content, err := os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("read setup script source %q: %w", source, err)
		}
		return content, nil
	default:
		return nil, fmt.Errorf("unsupported setup script type %q", s.Type)
	}
}

// WithCopyFiles returns a build step that copies each file declared in the
// agent config into the container. Files with Type "inline" use Source as the
// file contents; files with Type "path" read the contents from the host path in
// Source. Each file is written to its Target path inside the container, creating
// parent directories as needed. When a "path" entry points to a directory, the
// directory is copied recursively into the container at Target.
func WithCopyFiles(files []agentconfig.FileCopyConfig) ContainerBuildStep {
	return func(ctx context.Context, container *Container) error {
		for _, f := range files {
			var content []byte
			var mode os.FileMode = 0o644

			switch f.Type {
			case agentconfig.FileCopyTypeInline:
				content = []byte(f.Source)
			case agentconfig.FileCopyTypePath:
				source := f.Source
				if strings.HasPrefix(source, "~") {
					home, err := os.UserHomeDir()
					if err != nil {
						return fmt.Errorf("get user home directory: %w", err)
					}
					source = strings.Replace(source, "~", home, 1)
				}

				info, err := os.Stat(source)
				if err != nil {
					return fmt.Errorf("stat copy file source %q: %w", source, err)
				}
				if info.IsDir() {
					if err := container.CopyFromHost(ctx, container.ID(), source, f.Target); err != nil {
						return fmt.Errorf("copy directory %q to %q: %w", source, f.Target, err)
					}
					container.logger.Debug().Str("target", f.Target).Str("type", string(f.Type)).Msg("copied directory into container")
					continue
				}
				mode = info.Mode()
				content, err = os.ReadFile(source)
				if err != nil {
					return fmt.Errorf("read copy file source %q: %w", source, err)
				}
			default:
				return fmt.Errorf("unsupported copy file type %q for target %q", f.Type, f.Target)
			}

			if err := container.CopyFileFromHost(ctx, content, f.Target, mode); err != nil {
				return fmt.Errorf("copy file to %q: %w", f.Target, err)
			}

			container.logger.Debug().Str("target", f.Target).Str("type", string(f.Type)).Msg("copied file into container")
		}

		return nil
	}
}
