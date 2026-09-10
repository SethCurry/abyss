package agentconfig

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAgentConfig_SetupScriptsValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid file type",
			yaml: `
setup_scripts:
  - type: file
    source: /tmp/script.sh
`,
		},
		{
			name: "valid inline type",
			yaml: `
setup_scripts:
  - type: inline
    source: echo hello
`,
		},
		{
			name: "invalid setup script type",
			yaml: `
setup_scripts:
  - type: bogus
    source: echo hello
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg AgentConfig
			err := yaml.Unmarshal([]byte(tt.yaml), &cfg)
			if tt.wantErr {
				assertError(t, err)
				return
			}
			assertNoError(t, err)
		})
	}
}

func TestSetupScriptType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    SetupScriptType
		wantErr bool
	}{
		{
			name: "file",
			yaml: `"file"`,
			want: SetupScriptTypeFile,
		},
		{
			name: "inline",
			yaml: `"inline"`,
			want: SetupScriptTypeInline,
		},
		{
			name:    "invalid value",
			yaml:    `"bogus"`,
			wantErr: true,
		},
		{
			name:    "empty value",
			yaml:    `""`,
			wantErr: true,
		},
		{
			name:    "non-string node",
			yaml:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got SetupScriptType
			err := yaml.Unmarshal([]byte(tt.yaml), &got)
			if tt.wantErr {
				assertError(t, err)
				return
			}
			assertNoError(t, err)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetupScriptsConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       SetupScriptsConfig
		wantField string
	}{
		{
			name:      "valid file type",
			cfg:       SetupScriptsConfig{Type: SetupScriptTypeFile, Source: "/tmp/script.sh"},
			wantField: "",
		},
		{
			name:      "valid inline type",
			cfg:       SetupScriptsConfig{Type: SetupScriptTypeInline, Source: "echo hello"},
			wantField: "",
		},
		{
			name:      "valid empty type",
			cfg:       SetupScriptsConfig{Type: "", Source: "echo hello"},
			wantField: "",
		},
		{
			name:      "empty source",
			cfg:       SetupScriptsConfig{Type: SetupScriptTypeInline, Source: ""},
			wantField: "source",
		},
		{
			name:      "invalid type",
			cfg:       SetupScriptsConfig{Type: SetupScriptType("bogus"), Source: "echo hello"},
			wantField: "type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantField == "" {
				assertNoError(t, err)
				return
			}
			assertValidationError(t, err, tt.wantField)
		})
	}
}
