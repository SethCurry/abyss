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
