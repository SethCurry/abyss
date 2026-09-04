package agentconfig

import (
	"fmt"
	"os"

	"github.com/SethCurry/abyss/internal/types"
	"gopkg.in/yaml.v3"
)

/*
 * TODO features to add
 * I want to support these configurations:
 * - apt-get packages
 * - apk packages
 * - Language versions:
 * 	- Python
 *  - Node
 *  - These are just shims to apt-get/apk/etc
 */

type ToolsOnHostConfig struct {
	Files    bool `yaml:"files"`
	Terminal bool `yaml:"terminal"`
}

// Validate implements types.Validator. There is nothing to validate for
// boolean flags.
func (t ToolsOnHostConfig) Validate() error {
	return nil
}

type ACPConfig struct {
	ToolsOnHost ToolsOnHostConfig `yaml:"tools_on_host"`
}

// Validate implements types.Validator by validating the nested tools-on-host
// config.
func (a ACPConfig) Validate() error {
	return a.ToolsOnHost.Validate()
}

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

// SetupScriptType is an enum type.  It exists
// solely so UnmarshalYAML can validate the type field.
type SetupScriptType string

const (
	SetupScriptTypeFile   SetupScriptType = "file"
	SetupScriptTypeInline SetupScriptType = "inline"
)

// UnmarshalYAML validates that the Type field is either "file" or "inline".
func (t *SetupScriptType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	switch SetupScriptType(s) {
	case SetupScriptTypeFile, SetupScriptTypeInline:
		*t = SetupScriptType(s)
		return nil
	default:
		return fmt.Errorf("invalid setup script type %q: must be %q or %q", s, SetupScriptTypeFile, SetupScriptTypeInline)
	}
}

type SetupScriptsConfig struct {
	// Type is either "inline" or "file"
	Type SetupScriptType `yaml:"type"`
	// Source is the script content or file path
	Source string `yaml:"source"`
}

// Validate implements types.Validator by checking the source and type. An
// empty type is allowed and treated as inline.
func (s SetupScriptsConfig) Validate() error {
	if s.Source == "" {
		return types.NewValidationError(s, "source", "Cannot be an empty string")
	}

	switch s.Type {
	case "", SetupScriptTypeFile, SetupScriptTypeInline:
		return nil
	default:
		return types.NewValidationError(s, "type", fmt.Sprintf("must be %q or %q", SetupScriptTypeFile, SetupScriptTypeInline))
	}
}

// FileCopyType is an enum type.  It exists
// solely so UnmarshalYAML can validate the type field.
type FileCopyType string

const (
	FileCopyTypeInline FileCopyType = "inline"
	FileCopyTypePath   FileCopyType = "path"
)

// UnmarshalYAML validates that the Type field is either "inline" or "path".
func (t *FileCopyType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	switch FileCopyType(s) {
	case FileCopyTypeInline, FileCopyTypePath:
		*t = FileCopyType(s)
		return nil
	default:
		return fmt.Errorf("invalid file copy type %q: must be %q or %q", s, FileCopyTypeInline, FileCopyTypePath)
	}
}

type FileCopyConfig struct {
	// Type is currently either "inline" or "path"
	Type FileCopyType `yaml:"type"`
	// Source is the content or file path of the source
	Source string `yaml:"source"`
	// Target is the destination path on the agent
	Target string `yaml:"target"`
}

// Validate implements types.Validator by checking the source, target, and
// type.
func (f FileCopyConfig) Validate() error {
	if f.Source == "" {
		return types.NewValidationError(f, "source", "Cannot be an empty string")
	}

	if f.Target == "" {
		return types.NewValidationError(f, "target", "Cannot be an empty string")
	}

	switch f.Type {
	case FileCopyTypeInline, FileCopyTypePath:
		return nil
	default:
		return types.NewValidationError(f, "type", fmt.Sprintf("must be %q or %q", FileCopyTypeInline, FileCopyTypePath))
	}
}

type AgentConfig struct {
	Docker       DockerConfig         `yaml:"docker"`
	SetupScripts []SetupScriptsConfig `yaml:"setup_scripts"`
	CopyFiles    []FileCopyConfig     `yaml:"copy_files"`
	ACP          ACPConfig            `yaml:"acp"`
}

// Validate implements types.Validator by validating each nested config.
func (a AgentConfig) Validate() error {
	if err := a.Docker.Validate(); err != nil {
		return err
	}

	for _, script := range a.SetupScripts {
		if err := script.Validate(); err != nil {
			return err
		}
	}

	for _, file := range a.CopyFiles {
		if err := file.Validate(); err != nil {
			return err
		}
	}

	return a.ACP.Validate()
}

// FromYAMLFile reads the YAML file at the given path and unmarshals it
// into an *AgentConfig.
func FromYAMLFile(path string) (*AgentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AgentConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
