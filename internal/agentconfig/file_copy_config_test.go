package agentconfig

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFileCopyType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		want    FileCopyType
		wantErr bool
	}{
		{
			name: "inline",
			yaml: `"inline"`,
			want: FileCopyTypeInline,
		},
		{
			name: "path",
			yaml: `"path"`,
			want: FileCopyTypePath,
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
			var got FileCopyType
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

func TestAgentConfig_CopyFilesValidation(t *testing.T) {
	tests := []struct {
		name    string
		yaml    string
		wantErr bool
	}{
		{
			name: "valid inline type",
			yaml: `
copy_files:
  - type: inline
    source: hello
    target: /tmp/hello
`,
		},
		{
			name: "valid path type",
			yaml: `
copy_files:
  - type: path
    source: ./local.txt
    target: /tmp/local.txt
`,
		},
		{
			name: "invalid file copy type",
			yaml: `
copy_files:
  - type: bogus
    source: ./local.txt
    target: /tmp/local.txt
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

func TestFileCopyConfig_Validate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       FileCopyConfig
		wantField string
	}{
		{
			name:      "valid inline type",
			cfg:       FileCopyConfig{Type: FileCopyTypeInline, Source: "hello", Target: "/tmp/hello"},
			wantField: "",
		},
		{
			name:      "valid path type",
			cfg:       FileCopyConfig{Type: FileCopyTypePath, Source: "./local.txt", Target: "/tmp/local.txt"},
			wantField: "",
		},
		{
			name:      "empty source",
			cfg:       FileCopyConfig{Type: FileCopyTypeInline, Source: "", Target: "/tmp/hello"},
			wantField: "source",
		},
		{
			name:      "empty target",
			cfg:       FileCopyConfig{Type: FileCopyTypeInline, Source: "hello", Target: ""},
			wantField: "target",
		},
		{
			name:      "invalid type",
			cfg:       FileCopyConfig{Type: FileCopyType("bogus"), Source: "hello", Target: "/tmp/hello"},
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
