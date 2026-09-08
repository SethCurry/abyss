package acptools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

// TerminalTools implements the terminal-related methods of acp.Client. It runs
// commands as real OS processes and tracks their output and exit status,
// keyed by a generated terminal ID.
type TerminalTools struct {
	logger zerolog.Logger

	mu        sync.Mutex
	terminals map[string]*terminal
}

// NewTerminalTools constructs a TerminalTools with the given logger.
func NewTerminalTools(logger zerolog.Logger) *TerminalTools {
	return &TerminalTools{
		logger:    logger,
		terminals: make(map[string]*terminal),
	}
}

// terminal represents a single running (or completed) command process.
type terminal struct {
	cmd *exec.Cmd

	// output is the combined stdout/stderr capture buffer. It is safe for
	// concurrent use because cappedBuffer guards its own state.
	output *cappedBuffer

	// done is closed once the process has exited and its exit status recorded.
	done chan struct{}

	mu       sync.Mutex // guards exitCode and signal
	exitCode *int
	signal   *string
}

// cappedBuffer is an io.Writer that retains at most limit bytes of the most
// recent output, dropping older data from the beginning. A limit of 0 means
// unlimited. Truncation is performed on UTF-8 character boundaries so the
// retained bytes always decode to a valid string.
type cappedBuffer struct {
	mu        sync.Mutex
	limit     int
	buf       []byte
	truncated bool
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf = append(b.buf, p...)
	if b.limit > 0 && len(b.buf) > b.limit {
		// Keep only the last `limit` bytes, snapping to a rune boundary so we
		// don't slice a multi-byte UTF-8 sequence in half.
		start := len(b.buf) - b.limit
		for start < len(b.buf) && !utf8.RuneStart(b.buf[start]) {
			start++
		}
		b.buf = append([]byte(nil), b.buf[start:]...)
		b.truncated = true
	}
	return len(p), nil
}

// snapshot returns a copy of the currently captured output and whether any
// data has been dropped due to the byte limit.
func (b *cappedBuffer) snapshot() (string, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return string(b.buf), b.truncated
}

// CreateTerminal starts a new command process and returns a terminal ID that
// can be used to inspect its output and status.
func (e *TerminalTools) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	e.logger.Debug().Str("method", "CreateTerminal").Str("command", params.Command).Msg("handling request")

	if params.Command == "" {
		return acp.CreateTerminalResponse{}, fmt.Errorf("command is required")
	}

	id, err := newTerminalId()
	if err != nil {
		e.logger.Error().Err(err).Str("method", "CreateTerminal").Msg("failed to generate terminal id")
		return acp.CreateTerminalResponse{}, fmt.Errorf("failed to generate terminal id: %w", err)
	}

	// Use exec.Command (not CommandContext) so the process outlives this
	// request; the terminal's lifecycle is managed via the other methods.
	cmd := exec.Command(params.Command, params.Args...)
	if params.Cwd != nil && *params.Cwd != "" {
		cmd.Dir = *params.Cwd
	}
	if len(params.Env) > 0 {
		cmd.Env = buildEnv(params.Env)
	}

	limit := 0
	if params.OutputByteLimit != nil {
		limit = *params.OutputByteLimit
	}
	out := &cappedBuffer{limit: limit}
	cmd.Stdout = out
	cmd.Stderr = out

	t := &terminal{
		cmd:    cmd,
		output: out,
		done:   make(chan struct{}),
	}

	if err := cmd.Start(); err != nil {
		e.logger.Error().Err(err).Str("method", "CreateTerminal").Msg("failed to start command")
		return acp.CreateTerminalResponse{}, fmt.Errorf("failed to start command %q: %w", params.Command, err)
	}

	go func() {
		waitErr := cmd.Wait()
		t.mu.Lock()
		t.exitCode, t.signal = exitStatus(cmd.ProcessState, waitErr)
		t.mu.Unlock()
		close(t.done)
	}()

	e.mu.Lock()
	e.terminals[id] = t
	e.mu.Unlock()

	e.logger.Debug().Str("method", "CreateTerminal").Str("terminalId", id).Msg("terminal created")
	return acp.CreateTerminalResponse{
		TerminalId: id,
	}, nil
}

// TerminalOutput returns the output captured so far for a terminal along with
// its exit status, if the command has completed.
func (e *TerminalTools) TerminalOutput(ctx context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	e.logger.Debug().Str("method", "TerminalOutput").Str("terminalId", params.TerminalId).Msg("handling request")

	t, ok := e.getTerminal(params.TerminalId)
	if !ok {
		return acp.TerminalOutputResponse{}, fmt.Errorf("terminal %q not found", params.TerminalId)
	}

	output, truncated := t.output.snapshot()

	var exited bool
	select {
	case <-t.done:
		exited = true
	default:
	}

	resp := acp.TerminalOutputResponse{
		Output:    output,
		Truncated: truncated,
	}

	if exited {
		t.mu.Lock()
		exitCode := t.exitCode
		signal := t.signal
		t.mu.Unlock()
		resp.ExitStatus = &acp.TerminalExitStatus{
			ExitCode: exitCode,
			Signal:   signal,
		}
	}

	return resp, nil
}

// ReleaseTerminal frees the resources associated with a terminal, killing the
// underlying process if it is still running.
func (e *TerminalTools) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	e.logger.Debug().Str("method", "ReleaseTerminal").Str("terminalId", params.TerminalId).Msg("handling request")

	e.mu.Lock()
	t, ok := e.terminals[params.TerminalId]
	if ok {
		delete(e.terminals, params.TerminalId)
	}
	e.mu.Unlock()

	if !ok {
		return acp.ReleaseTerminalResponse{}, fmt.Errorf("terminal %q not found", params.TerminalId)
	}

	t.kill(e.logger)
	return acp.ReleaseTerminalResponse{}, nil
}

// WaitForTerminalExit blocks until the terminal's command exits (or the
// context is cancelled) and returns its exit status.
func (e *TerminalTools) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	e.logger.Debug().Str("method", "WaitForTerminalExit").Str("terminalId", params.TerminalId).Msg("handling request")

	t, ok := e.getTerminal(params.TerminalId)
	if !ok {
		return acp.WaitForTerminalExitResponse{}, fmt.Errorf("terminal %q not found", params.TerminalId)
	}

	select {
	case <-t.done:
	case <-ctx.Done():
		return acp.WaitForTerminalExitResponse{}, ctx.Err()
	}

	t.mu.Lock()
	exitCode := t.exitCode
	signal := t.signal
	t.mu.Unlock()

	return acp.WaitForTerminalExitResponse{
		ExitCode: exitCode,
		Signal:   signal,
	}, nil
}

// KillTerminal terminates the terminal's process without releasing the
// terminal; its output and exit status remain available afterwards.
func (e *TerminalTools) KillTerminal(ctx context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	e.logger.Debug().Str("method", "KillTerminal").Str("terminalId", params.TerminalId).Msg("handling request")

	t, ok := e.getTerminal(params.TerminalId)
	if !ok {
		return acp.KillTerminalResponse{}, fmt.Errorf("terminal %q not found", params.TerminalId)
	}

	t.kill(e.logger)
	return acp.KillTerminalResponse{}, nil
}

// getTerminal returns the terminal registered under id.
func (e *TerminalTools) getTerminal(id string) (*terminal, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	t, ok := e.terminals[id]
	return t, ok
}

// kill terminates the process if it is still running. It is a no-op if the
// process has already exited.
func (t *terminal) kill(logger zerolog.Logger) {
	// If the process has already exited there is nothing to do.
	select {
	case <-t.done:
		return
	default:
	}

	t.mu.Lock()
	proc := t.cmd.Process
	t.mu.Unlock()

	if proc == nil {
		return
	}

	if err := proc.Kill(); err != nil {
		logger.Error().Err(err).Str("method", "kill").Msg("failed to kill process")
	}
}

// buildEnv returns an environment slice for a child process: the current
// process environment with the provided variables overriding any existing
// entries.
func buildEnv(vars []acp.EnvVariable) []string {
	envMap := make(map[string]string)
	for _, kv := range os.Environ() {
		k, v, _ := strings.Cut(kv, "=")
		envMap[k] = v
	}
	for _, v := range vars {
		envMap[v.Name] = v.Value
	}
	env := make([]string, 0, len(envMap))
	for k, v := range envMap {
		env = append(env, k+"="+v)
	}
	return env
}

// newTerminalId generates a random hex identifier for a new terminal.
func newTerminalId() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
