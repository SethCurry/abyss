package acptools

import (
	"context"
	"testing"
	"time"

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
		Command:        "sh",
		Args:           []string{"-c", "printf '0123456789ABCDEF'"},
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