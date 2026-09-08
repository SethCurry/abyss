package acptools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

func newTestTerminalTools() *TerminalTools {
	return NewTerminalTools(zerolog.Nop())
}

func TestCreateTerminalAndWaitExit(t *testing.T) {
	tt := newTestTerminalTools()
	ctx := context.Background()

	resp, err := tt.CreateTerminal(ctx, acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "echo hello; exit 7"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	id := resp.TerminalId
	if id == "" {
		t.Fatal("expected terminal id")
	}

	wait, err := tt.WaitForTerminalExit(ctx, acp.WaitForTerminalExitRequest{TerminalId: id})
	if err != nil {
		t.Fatalf("WaitForTerminalExit: %v", err)
	}
	if wait.ExitCode == nil || *wait.ExitCode != 7 {
		t.Fatalf("expected exit code 7, got %+v", wait.ExitCode)
	}
	if wait.Signal != nil {
		t.Fatalf("expected nil signal, got %q", *wait.Signal)
	}

	out, err := tt.TerminalOutput(ctx, acp.TerminalOutputRequest{TerminalId: id})
	if err != nil {
		t.Fatalf("TerminalOutput: %v", err)
	}
	if out.Output != "hello\n" {
		t.Fatalf("expected output %q, got %q", "hello\n", out.Output)
	}
	if out.Truncated {
		t.Fatal("expected not truncated")
	}
	if out.ExitStatus == nil || out.ExitStatus.ExitCode == nil || *out.ExitStatus.ExitCode != 7 {
		t.Fatalf("expected exit status code 7, got %+v", out.ExitStatus)
	}

	// After release, the terminal should no longer be findable.
	if _, err := tt.ReleaseTerminal(ctx, acp.ReleaseTerminalRequest{TerminalId: id}); err != nil {
		t.Fatalf("ReleaseTerminal: %v", err)
	}
	if _, err := tt.TerminalOutput(ctx, acp.TerminalOutputRequest{TerminalId: id}); err == nil {
		t.Fatal("expected error after release")
	}
}

func TestKillTerminal(t *testing.T) {
	tt := newTestTerminalTools()
	ctx := context.Background()

	resp, err := tt.CreateTerminal(ctx, acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "sleep 30"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	id := resp.TerminalId

	if _, err := tt.KillTerminal(ctx, acp.KillTerminalRequest{TerminalId: id}); err != nil {
		t.Fatalf("KillTerminal: %v", err)
	}

	wait, err := tt.WaitForTerminalExit(ctx, acp.WaitForTerminalExitRequest{TerminalId: id})
	if err != nil {
		t.Fatalf("WaitForTerminalExit: %v", err)
	}
	if wait.Signal == nil {
		t.Fatal("expected a signal after kill")
	}
	if wait.ExitCode != nil {
		t.Fatalf("expected nil exit code when signaled, got %d", *wait.ExitCode)
	}
}

func TestWaitForTerminalExitCancelledContext(t *testing.T) {
	tt := newTestTerminalTools()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "sleep 30"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}

	if _, err := tt.WaitForTerminalExit(ctx, acp.WaitForTerminalExitRequest{TerminalId: resp.TerminalId}); err == nil {
		t.Fatal("expected context cancellation error")
	}

	// Clean up the lingering process.
	_, _ = tt.KillTerminal(context.Background(), acp.KillTerminalRequest{TerminalId: resp.TerminalId})
}

func TestTerminalOutputByteLimit(t *testing.T) {
	tt := newTestTerminalTools()
	ctx := context.Background()

	limit := 8
	resp, err := tt.CreateTerminal(ctx, acp.CreateTerminalRequest{
		Command:         "sh",
		Args:            []string{"-c", "printf '0123456789ABCDEF'"},
		OutputByteLimit: &limit,
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}

	if _, err := tt.WaitForTerminalExit(ctx, acp.WaitForTerminalExitRequest{TerminalId: resp.TerminalId}); err != nil {
		t.Fatalf("WaitForTerminalExit: %v", err)
	}

	out, err := tt.TerminalOutput(ctx, acp.TerminalOutputRequest{TerminalId: resp.TerminalId})
	if err != nil {
		t.Fatalf("TerminalOutput: %v", err)
	}
	if len(out.Output) > limit {
		t.Fatalf("output exceeded limit: %q (len %d)", out.Output, len(out.Output))
	}
	if !out.Truncated {
		t.Fatal("expected output to be truncated")
	}
	if out.Output != "89ABCDEF" {
		t.Fatalf("expected last 8 bytes, got %q", out.Output)
	}
}

// --- unit tests for the helper functions/types ---

func TestNewTerminalTools(t *testing.T) {
	tt := NewTerminalTools(zerolog.Nop())
	if tt == nil {
		t.Fatal("expected non-nil TerminalTools")
	}
	if tt.terminals == nil {
		t.Fatal("expected initialized terminals map")
	}
	if len(tt.terminals) != 0 {
		t.Fatalf("expected empty terminals map, got %d entries", len(tt.terminals))
	}
}

func TestNewTerminalId(t *testing.T) {
	id, err := newTerminalId()
	if err != nil {
		t.Fatalf("newTerminalId: %v", err)
	}
	// Should be 32 hex characters (16 bytes).
	if len(id) != 32 {
		t.Fatalf("expected 32 hex chars, got %d (%q)", len(id), id)
	}
	for _, r := range id {
		isHex := (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')
		if !isHex {
			t.Fatalf("expected lowercase hex, got non-hex char %q in %q", r, id)
		}
	}

	// Uniqueness across many calls.
	seen := make(map[string]struct{}, 1000)
	for i := 0; i < 1000; i++ {
		id, err := newTerminalId()
		if err != nil {
			t.Fatalf("newTerminalId[%d]: %v", i, err)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id %q at iteration %d", id, i)
		}
		seen[id] = struct{}{}
	}
}

func TestBuildEnvOverridesAndAdds(t *testing.T) {
	// Pick a name extremely unlikely to already be in the environment.
	name := "ABYSS_TEST_BUILDENV_UNIQUE_VAR_42"
	if _, ok := os.LookupEnv(name); ok {
		t.Fatalf("test precondition failed: %s already set", name)
	}

	env := buildEnv([]acp.EnvVariable{
		{Name: name, Value: "newvalue"},
	})

	found := false
	for _, kv := range env {
		if kv == name+"=newvalue" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected %q in env, got %v", name+"=newvalue", env)
	}
}

func TestBuildEnvOverridesExisting(t *testing.T) {
	// Use PATH, which should always be present in the test environment.
	env := buildEnv([]acp.EnvVariable{
		{Name: "PATH", Value: "/overridden"},
	})
	pathSeen := false
	dupPath := 0
	for _, kv := range env {
		if k, _, _ := strings.Cut(kv, "="); k == "PATH" {
			dupPath++
			if kv == "PATH=/overridden" {
				pathSeen = true
			}
		}
	}
	if !pathSeen {
		t.Fatalf("expected PATH to be overridden, got %v", env)
	}
	if dupPath != 1 {
		t.Fatalf("expected exactly one PATH entry, got %d", dupPath)
	}
}

func TestBuildEnvPreservesExisting(t *testing.T) {
	name := "ABYSS_TEST_BUILDENV_PRESERVE_99"
	if _, ok := os.LookupEnv(name); ok {
		t.Fatalf("test precondition failed: %s already set", name)
	}
	// Ensure there is at least one existing variable to preserve by checking
	// that the current environment's entries (other than our overrides) are
	// still represented.
	want := "ABYSS_TEST_BUILDENV_PRESERVE_99"
	_ = os.Setenv(want, "baseval")
	defer func() {
		_ = os.Unsetenv(want)
	}()

	env := buildEnv(nil)
	foundBase := false
	for _, kv := range env {
		if kv == want+"=baseval" {
			foundBase = true
		}
	}
	if !foundBase {
		t.Fatalf("expected existing env var %q preserved", want+"=baseval")
	}
}

func TestCappedBufferUnlimited(t *testing.T) {
	b := &cappedBuffer{limit: 0}
	data := make([]byte, 4096)
	for i := range data {
		data[i] = 'x'
	}
	if n, err := b.Write(data); n != len(data) || err != nil {
		t.Fatalf("Write returned (%d, %v), want (%d, nil)", n, err, len(data))
	}
	snap, truncated := b.snapshot()
	if truncated {
		t.Fatal("expected not truncated with unlimited buffer")
	}
	if len(snap) != len(data) {
		t.Fatalf("expected %d bytes, got %d", len(data), len(snap))
	}
}

func TestCappedBufferTruncatesToLastBytes(t *testing.T) {
	b := &cappedBuffer{limit: 4}
	if _, err := b.Write([]byte("0123456789")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	snap, truncated := b.snapshot()
	if !truncated {
		t.Fatal("expected truncated")
	}
	if snap != "6789" {
		t.Fatalf("expected %q, got %q", "6789", snap)
	}
}

func TestCappedBufferNoTruncateWhenUnderLimit(t *testing.T) {
	b := &cappedBuffer{limit: 100}
	if _, err := b.Write([]byte("short")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	snap, truncated := b.snapshot()
	if truncated {
		t.Fatal("expected not truncated")
	}
	if snap != "short" {
		t.Fatalf("expected %q, got %q", "short", snap)
	}
}

func TestCappedBufferTruncatesOnUTF8Boundary(t *testing.T) {
	// "é" is 2 bytes (0xC3 0xA9). With a limit of 3 bytes, after writing two
	// "é"s (4 bytes) we must keep the last rune boundary, so the retained data
	// is exactly one "é" (2 bytes), not a sliced continuation byte.
	b := &cappedBuffer{limit: 3}
	if _, err := b.Write([]byte("éé")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	snap, truncated := b.snapshot()
	if !truncated {
		t.Fatal("expected truncated")
	}
	if !utf8.ValidString(snap) {
		t.Fatalf("retained output is not valid UTF-8: %q", snap)
	}
	if snap != "é" {
		t.Fatalf("expected single é, got %q", snap)
	}
}

func TestCappedBufferSnapshotIsCopy(t *testing.T) {
	b := &cappedBuffer{limit: 0}
	if _, err := b.Write([]byte("abc")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	snap, _ := b.snapshot()
	// Mutating the returned snapshot must not affect subsequent snapshots.
	//nolint:ineffassign,staticcheck
	snap = snap + "xyz"
	snap2, _ := b.snapshot()
	if snap2 != "abc" {
		t.Fatalf("snapshot changed after mutating earlier result: got %q", snap2)
	}
}

func TestCappedBufferConcurrentWrites(t *testing.T) {
	b := &cappedBuffer{limit: 0}
	const goroutines = 16
	const writesPerG = 100
	const payload = "hello"

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < writesPerG; j++ {
				if _, err := b.Write([]byte(payload)); err != nil {
					t.Errorf("Write: %v", err)
				}
			}
		}()
	}
	wg.Wait()

	snap, _ := b.snapshot()
	want := goroutines * writesPerG * len(payload)
	if len(snap) != want {
		t.Fatalf("expected %d bytes, got %d", want, len(snap))
	}
}

func TestCreateTerminalEmptyCommand(t *testing.T) {
	tt := newTestTerminalTools()
	_, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "",
	})
	if err == nil {
		t.Fatal("expected error for empty command")
	}
	if !strings.Contains(err.Error(), "command is required") {
		t.Fatalf("expected 'command is required' error, got %v", err)
	}
}

func TestCreateTerminalBadCommand(t *testing.T) {
	tt := newTestTerminalTools()
	_, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "/nonexistent/program/that/does/not/exist",
	})
	if err == nil {
		t.Fatal("expected error for non-existent command")
	}
	if !strings.Contains(err.Error(), "failed to start command") {
		t.Fatalf("expected 'failed to start command' error, got %v", err)
	}
}

func TestCreateTerminalUsesCwd(t *testing.T) {
	tt := newTestTerminalTools()
	dir := t.TempDir()
	cwd := dir
	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "pwd"},
		Cwd:     &cwd,
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	if _, err := tt.WaitForTerminalExit(context.Background(), acp.WaitForTerminalExitRequest{TerminalId: resp.TerminalId}); err != nil {
		t.Fatalf("WaitForTerminalExit: %v", err)
	}
	out, err := tt.TerminalOutput(context.Background(), acp.TerminalOutputRequest{TerminalId: resp.TerminalId})
	if err != nil {
		t.Fatalf("TerminalOutput: %v", err)
	}
	// Resolve dir the same way the shell would; on some systems TempDir is a
	// symlink, so compare via filepath.EvalSymlinks.
	expected, _ := filepath.EvalSymlinks(dir)
	got, _ := filepath.EvalSymlinks(strings.TrimSpace(out.Output))
	if expected != got {
		t.Fatalf("expected cwd %q, got %q", expected, got)
	}
}

func TestCreateTerminalUsesEnv(t *testing.T) {
	tt := newTestTerminalTools()
	val := "envvalue123"
	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "printf $ABYSS_TEST_TERM_ENV"},
		Env: []acp.EnvVariable{
			{Name: "ABYSS_TEST_TERM_ENV", Value: val},
		},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	if _, err := tt.WaitForTerminalExit(context.Background(), acp.WaitForTerminalExitRequest{TerminalId: resp.TerminalId}); err != nil {
		t.Fatalf("WaitForTerminalExit: %v", err)
	}
	out, err := tt.TerminalOutput(context.Background(), acp.TerminalOutputRequest{TerminalId: resp.TerminalId})
	if err != nil {
		t.Fatalf("TerminalOutput: %v", err)
	}
	if out.Output != val {
		t.Fatalf("expected %q, got %q", val, out.Output)
	}
}

func TestTerminalOutputNotFound(t *testing.T) {
	tt := newTestTerminalTools()
	_, err := tt.TerminalOutput(context.Background(), acp.TerminalOutputRequest{TerminalId: "does-not-exist"})
	if err == nil {
		t.Fatal("expected error for missing terminal")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' error, got %v", err)
	}
}

func TestTerminalOutputBeforeExitHasNoExitStatus(t *testing.T) {
	tt := newTestTerminalTools()
	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "sleep 0.2; echo done"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	defer func() {
		_, _ = tt.KillTerminal(context.Background(), acp.KillTerminalRequest{TerminalId: resp.TerminalId})
	}()

	out, err := tt.TerminalOutput(context.Background(), acp.TerminalOutputRequest{TerminalId: resp.TerminalId})
	if err != nil {
		t.Fatalf("TerminalOutput: %v", err)
	}
	if out.ExitStatus != nil {
		t.Fatalf("expected nil ExitStatus before exit, got %+v", out.ExitStatus)
	}
}

func TestReleaseTerminalNotFound(t *testing.T) {
	tt := newTestTerminalTools()
	_, err := tt.ReleaseTerminal(context.Background(), acp.ReleaseTerminalRequest{TerminalId: "nope"})
	if err == nil {
		t.Fatal("expected error for missing terminal")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' error, got %v", err)
	}
}

func TestWaitForTerminalExitNotFound(t *testing.T) {
	tt := newTestTerminalTools()
	_, err := tt.WaitForTerminalExit(context.Background(), acp.WaitForTerminalExitRequest{TerminalId: "nope"})
	if err == nil {
		t.Fatal("expected error for missing terminal")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' error, got %v", err)
	}
}

func TestKillTerminalNotFound(t *testing.T) {
	tt := newTestTerminalTools()
	_, err := tt.KillTerminal(context.Background(), acp.KillTerminalRequest{TerminalId: "nope"})
	if err == nil {
		t.Fatal("expected error for missing terminal")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected 'not found' error, got %v", err)
	}
}

func TestReleaseTerminalKillsRunningProcess(t *testing.T) {
	tt := newTestTerminalTools()
	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "sleep 30"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	id := resp.TerminalId

	if _, err := tt.ReleaseTerminal(context.Background(), acp.ReleaseTerminalRequest{TerminalId: id}); err != nil {
		t.Fatalf("ReleaseTerminal: %v", err)
	}

	// The released terminal is removed and no longer findable.
	if _, ok := tt.getTerminal(id); ok {
		t.Fatal("expected terminal to be removed after release")
	}
}

func TestKillTerminalPreservesTerminal(t *testing.T) {
	tt := newTestTerminalTools()
	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "sleep 30"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	id := resp.TerminalId

	if _, err := tt.KillTerminal(context.Background(), acp.KillTerminalRequest{TerminalId: id}); err != nil {
		t.Fatalf("KillTerminal: %v", err)
	}

	// Unlike ReleaseTerminal, KillTerminal keeps the terminal registered.
	if _, ok := tt.getTerminal(id); !ok {
		t.Fatal("expected terminal to remain after kill")
	}
	if _, err := tt.KillTerminal(context.Background(), acp.KillTerminalRequest{TerminalId: id}); err != nil {
		t.Fatalf("KillTerminal second call: %v", err)
	}
}

func TestTerminalKillIsNoOpAfterExit(t *testing.T) {
	tt := newTestTerminalTools()
	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "exit 0"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	id := resp.TerminalId

	if _, err := tt.WaitForTerminalExit(context.Background(), acp.WaitForTerminalExitRequest{TerminalId: id}); err != nil {
		t.Fatalf("WaitForTerminalExit: %v", err)
	}

	// Killing an already-exited terminal must not panic or error.
	if _, err := tt.KillTerminal(context.Background(), acp.KillTerminalRequest{TerminalId: id}); err != nil {
		t.Fatalf("KillTerminal after exit: %v", err)
	}

	// Output is still available.
	out, err := tt.TerminalOutput(context.Background(), acp.TerminalOutputRequest{TerminalId: id})
	if err != nil {
		t.Fatalf("TerminalOutput: %v", err)
	}
	if out.ExitStatus == nil || out.ExitStatus.ExitCode == nil || *out.ExitStatus.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %+v", out.ExitStatus)
	}
}

func TestGetTerminal(t *testing.T) {
	tt := newTestTerminalTools()
	if _, ok := tt.getTerminal("missing"); ok {
		t.Fatal("expected not found for missing id")
	}

	resp, err := tt.CreateTerminal(context.Background(), acp.CreateTerminalRequest{
		Command: "sh",
		Args:    []string{"-c", "true"},
	})
	if err != nil {
		t.Fatalf("CreateTerminal: %v", err)
	}
	if _, ok := tt.getTerminal(resp.TerminalId); !ok {
		t.Fatal("expected to find created terminal")
	}
	_, _ = tt.ReleaseTerminal(context.Background(), acp.ReleaseTerminalRequest{TerminalId: resp.TerminalId})
	if _, ok := tt.getTerminal(resp.TerminalId); ok {
		t.Fatal("expected not found after release")
	}
}

func TestExitStatusNilProcessState(t *testing.T) {
	exitCode, signal := exitStatus(nil, nil)
	if exitCode != nil {
		t.Fatalf("expected nil exitCode, got %v", *exitCode)
	}
	if signal != nil {
		t.Fatalf("expected nil signal, got %q", *signal)
	}
}
