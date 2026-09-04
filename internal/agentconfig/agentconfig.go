package agentconfig

import (
	"os"

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

type WebsocketConfig struct {
	DisableTLS bool `yaml:"disable_tls"`
}

type AgentConfig struct {
	Docker       DockerConfig         `yaml:"docker"`
	SetupScripts []SetupScriptsConfig `yaml:"setup_scripts"`
	CopyFiles    []FileCopyConfig     `yaml:"copy_files"`
	ACP          ACPConfig            `yaml:"acp"`
	Websocket    WebsocketConfig      `yaml:"websocket"`
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
