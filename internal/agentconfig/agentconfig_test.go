package agentconfig

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
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

func TestUnmarshalYAML_ErrorMessages(t *testing.T) {
	t.Run("setup script error message", func(t *testing.T) {
		var st SetupScriptType
		err := yaml.Unmarshal([]byte(`"bogus"`), &st)
		assertError(t, err)
		if !strings.Contains(err.Error(), "invalid setup script type") {
			t.Fatalf("error %q does not mention setup script type", err)
		}
	})

	t.Run("file copy error message", func(t *testing.T) {
		var ft FileCopyType
		err := yaml.Unmarshal([]byte(`"bogus"`), &ft)
		assertError(t, err)
		if !strings.Contains(err.Error(), "invalid file copy type") {
			t.Fatalf("error %q does not mention file copy type", err)
		}
	})
}