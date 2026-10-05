package runenv

import (
	"strings"
	"testing"

	"github.com/moby/moby/api/types/container"
)

func TestNewErrInvalidContainerMetadata(t *testing.T) {
	err := NewErrInvalidContainerMetadata("some reason")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.message != "some reason" {
		t.Fatalf("message = %q, want %q", err.message, "some reason")
	}
}

func TestErrInvalidContainerMetadata_Error(t *testing.T) {
	err := NewErrInvalidContainerMetadata("some reason")
	want := "cannot get abyss container labels: some reason"
	if got := err.Error(); got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func TestGetLabelMetadata(t *testing.T) {
	tests := []struct {
		name    string
		summary *container.Summary
		want    *LabelMetadata
		wantErr string
	}{
		{
			name:    "nil labels",
			summary: &container.Summary{Labels: nil},
			wantErr: "container has no labels",
		},
		{
			name:    "missing abyss_version",
			summary: &container.Summary{Labels: map[string]string{}},
			wantErr: "container has no abyss_version label",
		},
		{
			name: "missing abyss_agent_config_path",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version": "v1",
			}},
			wantErr: "container has no abyss_agent_config_path label",
		},
		{
			name: "missing abyss_agent_config_hash",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
			}},
			wantErr: "container has no abyss_agent_config_hash label",
		},
		{
			name: "missing persistent label",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
			}},
			want: &LabelMetadata{
				AbyssVersion:    "v1",
				AgentConfigPath: "/path",
				AgentConfigHash: "hash",
				Persistent:      false,
			},
		},
		{
			name: "persistent label true",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
				"abyss_persistent":        "true",
			}},
			want: &LabelMetadata{
				AbyssVersion:    "v1",
				AgentConfigPath: "/path",
				AgentConfigHash: "hash",
				Persistent:      true,
			},
		},
		{
			name: "persistent label false",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
				"abyss_persistent":        "false",
			}},
			want: &LabelMetadata{
				AbyssVersion:    "v1",
				AgentConfigPath: "/path",
				AgentConfigHash: "hash",
				Persistent:      false,
			},
		},
		{
			name: "invalid persistent label",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
				"abyss_persistent":        "not-a-bool",
			}},
			wantErr: "abyss_persistent label is not a valid boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLabelMetadata(tt.summary)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %q, want to contain %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if *got != *tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestIsPersistent(t *testing.T) {
	tests := []struct {
		name    string
		summary *container.Summary
		want    bool
	}{
		{
			name:    "nil labels",
			summary: &container.Summary{Labels: nil},
			want:    false,
		},
		{
			name: "persistent",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
				"abyss_persistent":        "true",
			}},
			want: true,
		},
		{
			name: "not persistent",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
				"abyss_persistent":        "false",
			}},
			want: false,
		},
		{
			name: "missing persistent label",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version":           "v1",
				"abyss_agent_config_path": "/path",
				"abyss_agent_config_hash": "hash",
			}},
			want: false,
		},
		{
			name: "malformed metadata",
			summary: &container.Summary{Labels: map[string]string{
				"abyss_version": "v1",
			}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPersistent(tt.summary); got != tt.want {
				t.Errorf("IsPersistent() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestLabelMetadata_ToMap(t *testing.T) {
	meta := &LabelMetadata{
		AbyssVersion:    "v1",
		AgentConfigPath: "/path",
		AgentConfigHash: "hash",
		Persistent:      true,
	}
	want := map[string]string{
		"abyss_version":           "v1",
		"abyss_agent_config_path": "/path",
		"abyss_agent_config_hash": "hash",
		"abyss_persistent":        "true",
	}

	got := meta.ToMap()
	if len(got) != len(want) {
		t.Fatalf("ToMap() returned %d entries, want %d", len(got), len(want))
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("ToMap()[%q] = %q, want %q", k, got[k], v)
		}
	}
}

func TestLabelMetadata_AddToMap(t *testing.T) {
	meta := &LabelMetadata{
		AbyssVersion:    "v1",
		AgentConfigPath: "/path",
		AgentConfigHash: "hash",
		Persistent:      true,
	}

	t.Run("into empty map", func(t *testing.T) {
		m := map[string]string{}
		meta.AddToMap(m)
		want := map[string]string{
			"abyss_version":           "v1",
			"abyss_agent_config_path": "/path",
			"abyss_agent_config_hash": "hash",
			"abyss_persistent":        "true",
		}
		for k, v := range want {
			if m[k] != v {
				t.Errorf("m[%q] = %q, want %q", k, m[k], v)
			}
		}
		if len(m) != len(want) {
			t.Errorf("len(m) = %d, want %d", len(m), len(want))
		}
	})

	t.Run("overwrites existing keys and keeps others", func(t *testing.T) {
		m := map[string]string{
			"abyss_version":    "old",
			"abyss_persistent": "false",
			"other":            "keep",
		}
		meta.AddToMap(m)
		if m["abyss_version"] != "v1" {
			t.Errorf("abyss_version = %q, want %q", m["abyss_version"], "v1")
		}
		if m["abyss_persistent"] != "true" {
			t.Errorf("abyss_persistent = %q, want %q", m["abyss_persistent"], "true")
		}
		if m["other"] != "keep" {
			t.Errorf("other = %q, want %q", m["other"], "keep")
		}
	})
}

func TestLabelMetadata_ToMapAddToMapConsistent(t *testing.T) {
	meta := &LabelMetadata{
		AbyssVersion:    "v1",
		AgentConfigPath: "/path",
		AgentConfigHash: "hash",
		Persistent:      true,
	}
	m := map[string]string{"existing": "value"}
	meta.AddToMap(m)

	toMap := meta.ToMap()
	for k, v := range toMap {
		if m[k] != v {
			t.Errorf("AddToMap[%q] = %q, ToMap()[%q] = %q", k, m[k], k, v)
		}
	}
}
