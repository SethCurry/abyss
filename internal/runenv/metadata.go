package runenv

import (
	"github.com/moby/moby/api/types/container"
)

func NewErrInvalidContainerMetadata(msg string) *ErrInvalidContainerMetadata {
	return &ErrInvalidContainerMetadata{
		message: msg,
	}
}

type ErrInvalidContainerMetadata struct {
	message string
}

func (e *ErrInvalidContainerMetadata) Error() string {
	return "cannot get abyss container labels: " + e.message
}

func GetLabelMetadata(summary *container.Summary) (*LabelMetadata, error) {
	labels := summary.Labels
	if labels == nil {
		return nil, NewErrInvalidContainerMetadata("container has no labels")
	}

	abyssVersion, ok := labels["abyss_version"]
	if !ok {
		return nil, NewErrInvalidContainerMetadata("container has no abyss_version label")
	}
	agentConfigPath, ok := labels["abyss_agent_config_path"]
	if !ok {
		return nil, NewErrInvalidContainerMetadata("container has no abyss_agent_config_path label")
	}
	agentConfigHash, ok := labels["abyss_agent_config_hash"]
	if !ok {
		return nil, NewErrInvalidContainerMetadata("container has no abyss_agent_config_hash label")
	}

	return &LabelMetadata{
		AbyssVersion:    abyssVersion,
		AgentConfigPath: agentConfigPath,
		AgentConfigHash: agentConfigHash,
	}, nil
}

type LabelMetadata struct {
	AbyssVersion    string
	AgentConfigPath string
	AgentConfigHash string
}

func (l *LabelMetadata) ToMap() map[string]string {
	return map[string]string{
		"abyss_version":           l.AbyssVersion,
		"abyss_agent_config_path": l.AgentConfigPath,
		"abyss_agent_config_hash": l.AgentConfigHash,
	}
}

func (l *LabelMetadata) AddToMap(underlying map[string]string) {
	underlying["abyss_version"] = l.AbyssVersion
	underlying["abyss_agent_config_path"] = l.AgentConfigPath
	underlying["abyss_agent_config_hash"] = l.AgentConfigHash
}
