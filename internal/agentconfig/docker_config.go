package agentconfig

import (
	"os"

	"github.com/SethCurry/abyss/internal/types"
)

// HostMount maps a directory on the host into the Docker container.
type HostMount struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

// Validate implements types.Validator by checking that the mount source
// exists and that both source and destination are non-empty.
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

// ImagePullPolicy controls when Docker pulls the configured image.
type ImagePullPolicy string

const (
	// ImagePullPolicyAlways always pulls the configured image.
	ImagePullPolicyAlways ImagePullPolicy = "Always"
	// ImagePullPolicyIfNotPresent pulls the image only when it is not
	// present locally.
	ImagePullPolicyIfNotPresent ImagePullPolicy = "IfNotPresent"
	// ImagePullPolicyNever never pulls the image.
	ImagePullPolicyNever ImagePullPolicy = "Never"
)

// Validate implements types.Validator by rejecting values outside the
// allowed set. An empty value is treated as unset and is valid.
func (p ImagePullPolicy) Validate() error {
	switch p {
	case "", ImagePullPolicyAlways, ImagePullPolicyIfNotPresent, ImagePullPolicyNever:
		return nil
	default:
		return types.NewValidationError(p, "image_pull_policy", "must be one of Always, IfNotPresent, or Never")
	}
}

// DockerConfig holds the Docker-specific configuration for an agent.
type DockerConfig struct {
	// The Docker image to use.  Can be short or long, Docker will
	// resolve it for short names.
	Image           string          `yaml:"image"`
	ImagePullPolicy ImagePullPolicy `yaml:"image_pull_policy"`
	HostMounts      []HostMount     `yaml:"host_mounts"`
	AgentCommand    []string        `yaml:"agent_command"`
}

// Validate implements types.Validator by checking the image and each host
// mount.
func (d DockerConfig) Validate() error {
	if d.Image == "" {
		return types.NewValidationError(d, "image", "Cannot be an empty string")
	}

	if err := d.ImagePullPolicy.Validate(); err != nil {
		return err
	}

	for _, mount := range d.HostMounts {
		if err := mount.Validate(); err != nil {
			return err
		}
	}

	return nil
}
