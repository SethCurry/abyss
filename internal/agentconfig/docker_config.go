package agentconfig

import (
	"os"

	"github.com/SethCurry/abyss/internal/types"
)

type HostMount struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

func (h HostMount) Validate() error {
	// TODO throw an error if source doesn't exist
	if h.Source == "" {
		return types.NewValidationError(h, "source", "Cannot be an empty string")
	}

	if _, err := os.Stat(h.Source); err != nil && os.IsNotExist(err) {
		return types.NewValidationError(h, "source", "Path does not exist")
	}

	if h.Destination == "" {
		return types.NewValidationError(h, "destination", "Cannot be an empty string")
	}

	return nil
}

type DockerConfig struct {
	// The Docker image to use.  Can be short or long, Docker will
	// resolve it for short names.
	Image        string      `yaml:"image"`
	HostMounts   []HostMount `yaml:"host_mounts"`
	AgentCommand []string    `yaml:"agent_command"`
}

// Validate implements types.Validator by checking the image and each host
// mount.
func (d DockerConfig) Validate() error {
	if d.Image == "" {
		return types.NewValidationError(d, "image", "Cannot be an empty string")
	}

	for _, mount := range d.HostMounts {
		if err := mount.Validate(); err != nil {
			return err
		}
	}

	return nil
}
