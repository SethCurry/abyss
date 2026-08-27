package agentconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

/*
 * I want to support these configurations:
 * - apt-get packages
 * - apk packages
 * - Language versions:
 * 	- Python
 *  - Node
 *  - These are just shims to apt-get/apk/etc
 */

type HostMount struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}

type DockerConfig struct {
	Image        string      `yaml:"image"`
	HostMounts   []HostMount `yaml:"host_mounts"`
	AgentCommand []string    `yaml:"agent_command"`
}

type SetupScriptType string

const (
	SetupScriptTypeFile   SetupScriptType = "file"
	SetupScriptTypeInline SetupScriptType = "inline"
)

type SetupScriptsConfig struct {
	// Type is either "inline" or "file"
	Type SetupScriptType `yaml:"type"`
	// Source is the script content or file path
	Source string `yaml:"source"`
}

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

type FileCopyType string

const (
	FileCopyTypeInline FileCopyType = "inline"
	FileCopyTypePath   FileCopyType = "path"
)

type FileCopyConfig struct {
	// Type is currently either "inline" or "path"
	Type FileCopyType `yaml:"type"`
	// Source is the content or file path of the source
	Source string `yaml:"source"`
	// Target is the destination path on the agent
	Target string `yaml:"target"`
}

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

type AgentConfig struct {
	Docker       DockerConfig         `yaml:"docker"`
	SetupScripts []SetupScriptsConfig `yaml:"setup_scripts"`
	CopyFiles    []FileCopyConfig     `yaml:"copy_files"`
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
