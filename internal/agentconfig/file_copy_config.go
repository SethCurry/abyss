package agentconfig

import (
	"fmt"

	"github.com/SethCurry/abyss/internal/types"
	"gopkg.in/yaml.v3"
)

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
