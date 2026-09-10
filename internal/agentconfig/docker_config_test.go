package agentconfig

import (
	"os"
	"path/filepath"
	"testing"
)

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
