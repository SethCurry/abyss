package agentconfig

// ToolsOnHostConfig holds the boolean flags that toggle the ACP tools that
// run on the host, such as file and terminal access.
type ToolsOnHostConfig struct {
	Files    bool `yaml:"files"`
	Terminal bool `yaml:"terminal"`
}

// Validate implements types.Validator. There is nothing to validate for
// boolean flags.
func (t ToolsOnHostConfig) Validate() error {
	return nil
}

// ACPConfig is the root of the ACP-specific agent configuration.
type ACPConfig struct {
	ToolsOnHost ToolsOnHostConfig `yaml:"tools_on_host"`
}

// Validate implements types.Validator by validating the nested tools-on-host
// config.
func (a ACPConfig) Validate() error {
	return a.ToolsOnHost.Validate()
}
