package agentconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SethCurry/abyss/internal/types"
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

func assertValidationError(t *testing.T, err error, field string) {
	t.Helper()
	var ve *types.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected *types.ValidationError, got %T: %v", err, err)
	}
	if ve.Field != field {
		t.Fatalf("expected field %q, got %q", field, ve.Field)
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

func TestToolsOnHostConfig_Validate(t *testing.T) {
	assertNoError(t, ToolsOnHostConfig{}.Validate())
	assertNoError(t, ToolsOnHostConfig{Files: true, Terminal: true}.Validate())
}

func TestACPConfig_Validate(t *testing.T) {
	assertNoError(t, ACPConfig{}.Validate())
	assertNoError(t, ACPConfig{ToolsOnHost: ToolsOnHostConfig{Files: true}}.Validate())
}

func TestHostMount_Validate(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tests := []struct {
		name      string
		mount     HostMount
		wantField string
	}{
		{
			name:      "valid",
			mount:     HostMount{Source: existingFile, Destination: "/tmp/dest"},
			wantField: "",
		},
		{
			name:      "empty source",
			mount:     HostMount{Source: "", Destination: "/tmp/dest"},
			wantField: "source",
		},
		{
			name:      "source does not exist",
			mount:     HostMount{Source: filepath.Join(tempDir, "missing.txt"), Destination: "/tmp/dest"},
			wantField: "source",
		},
		{
			name:      "empty destination",
			mount:     HostMount{Source: existingFile, Destination: ""},
			wantField: "destination",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mount.Validate()
			if tt.wantField == "" {
				assertNoError(t, err)
				return
			}
			assertValidationError(t, err, tt.wantField)
		})
	}
}

func TestDockerConfig_Validate(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tests := []struct {
		name      string
		cfg       DockerConfig
		wantField string
	}{
		{
			name:      "valid",
			cfg:       DockerConfig{Image: "ubuntu:latest"},
			wantField: "",
		},
		{
			name:      "valid with mount",
			cfg:       DockerConfig{Image: "ubuntu:latest", HostMounts: []HostMount{{Source: existingFile, Destination: "/tmp/dest"}}},
			wantField: "",
		},
		{
			name:      "empty image",
			cfg:       DockerConfig{Image: ""},
			wantField: "image",
		},
		{
			name:      "invalid mount",
			cfg:       DockerConfig{Image: "ubuntu:latest", HostMounts: []HostMount{{Source: "", Destination: "/tmp/dest"}}},
			wantField: "source",
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

func TestAgentConfig_Validate(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tests := []struct {
		name      string
		cfg       AgentConfig
		wantField string
	}{
		{
			name: "valid",
			cfg: AgentConfig{
				Docker: DockerConfig{Image: "ubuntu:latest"},
			},
			wantField: "",
		},
		{
			name: "valid full",
			cfg: AgentConfig{
				Docker: DockerConfig{
					Image:      "ubuntu:latest",
					HostMounts: []HostMount{{Source: existingFile, Destination: "/tmp/dest"}},
				},
				SetupScripts: []SetupScriptsConfig{{Type: SetupScriptTypeInline, Source: "echo hello"}},
				CopyFiles:    []FileCopyConfig{{Type: FileCopyTypeInline, Source: "hello", Target: "/tmp/hello"}},
			},
			wantField: "",
		},
		{
			name:      "invalid docker image",
			cfg:       AgentConfig{Docker: DockerConfig{Image: ""}},
			wantField: "image",
		},
		{
			name: "invalid setup script",
			cfg: AgentConfig{
				Docker:       DockerConfig{Image: "ubuntu:latest"},
				SetupScripts: []SetupScriptsConfig{{Type: SetupScriptTypeInline, Source: ""}},
			},
			wantField: "source",
		},
		{
			name: "invalid copy file",
			cfg: AgentConfig{
				Docker:    DockerConfig{Image: "ubuntu:latest"},
				CopyFiles: []FileCopyConfig{{Type: FileCopyTypeInline, Source: "hello", Target: ""}},
			},
			wantField: "target",
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