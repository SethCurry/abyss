package agentconfig

import (
	"fmt"

	"github.com/SethCurry/abyss/internal/types"
	"gopkg.in/yaml.v3"
)

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
