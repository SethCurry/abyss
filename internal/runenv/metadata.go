package runenv

import (
	"github.com/moby/moby/api/types/container"
)

// NewErrInvalidContainerMetadata builds an error for invalid or missing
// container labels.
func NewErrInvalidContainerMetadata(msg string) *ErrInvalidContainerMetadata {
	return &ErrInvalidContainerMetadata{
		message: msg,
	}
}

// ErrInvalidContainerMetadata is returned when a container's abyss labels are
// malformed or absent.
type ErrInvalidContainerMetadata struct {
	message string
}

// Error formats the invalid metadata error.
func (e *ErrInvalidContainerMetadata) Error() string {
	return "cannot get abyss container labels: " + e.message
}

// GetLabelMetadata extracts abyss metadata from a container's labels.
// Returns an error for missing labels
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

// LabelMetadata holds the abyss metadata stored in container labels.
type LabelMetadata struct {
	AbyssVersion    string
	AgentConfigPath string
	AgentConfigHash string
}

// ToMap returns the metadata as a map of container labels.
func (l *LabelMetadata) ToMap() map[string]string {
	return map[string]string{
		"abyss_version":           l.AbyssVersion,
		"abyss_agent_config_path": l.AgentConfigPath,
		"abyss_agent_config_hash": l.AgentConfigHash,
	}
}

// AddToMap writes the metadata labels into the given map.
func (l *LabelMetadata) AddToMap(underlying map[string]string) {
	underlying["abyss_version"] = l.AbyssVersion
	underlying["abyss_agent_config_path"] = l.AgentConfigPath
	underlying["abyss_agent_config_hash"] = l.AgentConfigHash
}
