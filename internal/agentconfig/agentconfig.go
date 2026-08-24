package agentconfig

import (
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

type DockerConfig struct {
	Image        string   `yaml:"image"`
	HostMounts   []string `yaml:"host_mounts"`
	AgentCommand []string `yaml:"agent_command"`
}

type SetupScriptsConfig struct {
	// Type is either "inline" or "file"
	Type string `yaml:"type"`
	// Source is the script content or file path
	Source string `yaml:"source"`
}

type FileCopyConfig struct {
	// Type is currently either "inline" or "path"
	Type string `yaml:"type"`
	// Source is the content or file path of the source
	Source string `yaml:"source"`
	// Target is the destination path on the agent
	Target string `yaml:"target"`
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
