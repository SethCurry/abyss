package agentconfig

import (
	"os"
	"strings"

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

func BuildAgentProxyArgs(cfg *AgentConfig) string {
	agent := cfg.Docker.AgentCommand

	agentArgs := make([]string, len(agent)*2)

	for k, v := range agent {
		startIndex := k * 2
		agentArgs[startIndex] = "--agent"
		agentArgs[startIndex+1] = v
	}

	if cfg.ACP.ToolsOnHost.Files {
		agentArgs = append(agentArgs, "--local-filesystem")
	}

	if cfg.ACP.ToolsOnHost.Terminal {
		agentArgs = append(agentArgs, "--local-terminal")
	}

	if !cfg.Websocket.DisableTLS {
		agentArgs = append(agentArgs,
			"--tls-cert", DefaultTLSServerCertPath,
			"--tls-key", DefaultTLSServerKeyPath,
			"--tls-ca", DefaultTLSCACertPath,
		)
	}

	joinedArgs := strings.Join(agentArgs, " ")

	return "/usr/local/bin/abyss server " + joinedArgs
}
