package agentconfig

import "testing"

func TestToolsOnHostConfig_Validate(t *testing.T) {
	assertNoError(t, ToolsOnHostConfig{}.Validate())
	assertNoError(t, ToolsOnHostConfig{Files: true, Terminal: true}.Validate())
}

func TestACPConfig_Validate(t *testing.T) {
	assertNoError(t, ACPConfig{}.Validate())
	assertNoError(t, ACPConfig{ToolsOnHost: ToolsOnHostConfig{Files: true}}.Validate())
}
