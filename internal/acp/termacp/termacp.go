package termacp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/SethCurry/abyss/internal/acptools"
	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

// NewTermACPClient creates a terminal-backed ACP client. Session updates and
// permission prompts are rendered to os.Stdout, and permission selections are
// read from os.Stdin. File and terminal operations are backed by the built-in
// acptools implementations.
// Compile-time check that TermACPClient satisfies the acp.Client interface.
var _ acp.Client = (*TermACPClient)(nil)

func NewTermACPClient() *TermACPClient {
	logger := zerolog.New(io.Discard)
	return &TermACPClient{
		logger: logger,
		out:    os.Stdout,
		in:     bufio.NewReader(os.Stdin),
		fs:     acptools.NewFilesystemTools(logger),
		term:   acptools.NewTerminalTools(logger),
	}
}

// TermACPClient implements acp.Client against a local terminal: it prints
// session updates, prompts the user for permission decisions, and delegates
// file and terminal operations to acptools.
type TermACPClient struct {
	logger zerolog.Logger
	out    io.Writer
	in     *bufio.Reader
	fs     *acptools.FilesystemTools
	term   *acptools.TerminalTools
}

// RequestPermission prompts the user on the terminal to choose one of the
// provided permission options for a tool call. An invalid or empty selection
// (or a cancelled context) is reported as cancelled rather than erroring, so
// the agent can wind down the turn cleanly.
func (e *TermACPClient) RequestPermission(ctx context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	title := toolCallTitle(params.ToolCall)

	fmt.Fprintf(e.out, "\n\033[1mPermission required\033[0m")
	if title != "" {
		fmt.Fprintf(e.out, ": %s", title)
	}
	fmt.Fprintln(e.out)

	for i, opt := range params.Options {
		fmt.Fprintf(e.out, "  [%d] %s\n", i+1, opt.Name)
	}
	fmt.Fprintf(e.out, "Select an option (1-%d) [cancel]: ", len(params.Options))

	line, ok := e.readLine(ctx)
	if !ok || line == "" {
		return acp.RequestPermissionResponse{
			Outcome: acp.RequestPermissionOutcome{
				Cancelled: &acp.RequestPermissionOutcomeCancelled{},
			},
		}, nil
	}

	idx, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || idx < 1 || idx > len(params.Options) {
		return acp.RequestPermissionResponse{
			Outcome: acp.RequestPermissionOutcome{
				Cancelled: &acp.RequestPermissionOutcomeCancelled{},
			},
		}, nil
	}

	return acp.RequestPermissionResponse{
		Outcome: acp.RequestPermissionOutcome{
			Selected: &acp.RequestPermissionOutcomeSelected{
				OptionId: params.Options[idx-1].OptionId,
			},
		},
	}, nil
}

// SessionUpdate renders a session update to the terminal. Only the update
// kinds relevant to a terminal client are rendered; the rest are ignored.
func (e *TermACPClient) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	switch {
	case params.Update.AgentMessageChunk != nil:
		fmt.Fprint(e.out, contentBlockText(params.Update.AgentMessageChunk.Content))
	case params.Update.AgentThoughtChunk != nil:
		text := contentBlockText(params.Update.AgentThoughtChunk.Content)
		if text != "" {
			fmt.Fprintf(e.out, "\033[2m%s\033[0m", text)
		}
	case params.Update.UserMessageChunk != nil:
		fmt.Fprint(e.out, contentBlockText(params.Update.UserMessageChunk.Content))
	case params.Update.ToolCall != nil:
		tc := params.Update.ToolCall
		fmt.Fprintf(e.out, "\n\033[1mtool:\033[0m %s\n", tc.Title)
		for _, c := range tc.Content {
			fmt.Fprint(e.out, renderToolCallContent(c))
		}
	case params.Update.ToolCallUpdate != nil:
		tc := params.Update.ToolCallUpdate
		if tc.Title != nil {
			fmt.Fprintf(e.out, "\n\033[1mtool:\033[0m %s\n", *tc.Title)
		}
		if tc.Status != nil {
			fmt.Fprintf(e.out, "\033[2m[%s]\033[0m\n", *tc.Status)
		}
		for _, c := range tc.Content {
			fmt.Fprint(e.out, renderToolCallContent(c))
		}
	case params.Update.Plan != nil:
		fmt.Fprintln(e.out, "\n\033[1mplan\033[0m")
		for _, entry := range params.Update.Plan.Entries {
			fmt.Fprintf(e.out, "  [%s] %s\n", entry.Status, entry.Content)
		}
	}
	return nil
}

// WriteTextFile delegates to the local filesystem implementation.
func (e *TermACPClient) WriteTextFile(ctx context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	return e.fs.WriteTextFile(ctx, params)
}

// ReadTextFile delegates to the local filesystem implementation.
func (e *TermACPClient) ReadTextFile(ctx context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	return e.fs.ReadTextFile(ctx, params)
}

// CreateTerminal delegates to the local terminal implementation.
func (e *TermACPClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	return e.term.CreateTerminal(ctx, params)
}

// TerminalOutput delegates to the local terminal implementation.
func (e *TermACPClient) TerminalOutput(ctx context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	return e.term.TerminalOutput(ctx, params)
}

// ReleaseTerminal delegates to the local terminal implementation.
func (e *TermACPClient) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	return e.term.ReleaseTerminal(ctx, params)
}

// WaitForTerminalExit delegates to the local terminal implementation.
func (e *TermACPClient) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	return e.term.WaitForTerminalExit(ctx, params)
}

// KillTerminal implements acp.Client.
func (e *TermACPClient) KillTerminal(ctx context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	return e.term.KillTerminal(ctx, params)
}

// readLine reads a single line of user input, returning ("", false) if the
// context is cancelled before input arrives.
func (e *TermACPClient) readLine(ctx context.Context) (string, bool) {
	type result struct {
		line string
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		line, err := e.in.ReadString('\n')
		ch <- result{line, err}
	}()

	select {
	case <-ctx.Done():
		return "", false
	case r := <-ch:
		if r.err != nil && r.err != io.EOF {
			return "", false
		}
		return strings.TrimRight(r.line, "\r\n"), true
	}
}

// toolCallTitle extracts a human-readable title from a tool call update.
func toolCallTitle(tc acp.ToolCallUpdate) string {
	if tc.Title != nil {
		return *tc.Title
	}
	return string(tc.ToolCallId)
}

// contentBlockText returns the text of a text content block, or the empty
// string for non-text blocks.
func contentBlockText(block acp.ContentBlock) string {
	if block.Text != nil {
		return block.Text.Text
	}
	return ""
}

// renderToolCallContent formats a tool call content entry for the terminal.
func renderToolCallContent(c acp.ToolCallContent) string {
	switch {
	case c.Content != nil:
		return contentBlockText(c.Content.Content) + "\n"
	case c.Diff != nil:
		var b strings.Builder
		fmt.Fprintf(&b, "diff: %s\n", c.Diff.Path)
		if c.Diff.OldText != nil {
			fmt.Fprintf(&b, "\033[31m-%s\033[0m\n", *c.Diff.OldText)
		}
		fmt.Fprintf(&b, "\033[32m+%s\033[0m\n", c.Diff.NewText)
		return b.String()
	case c.Terminal != nil:
		return fmt.Sprintf("terminal: %s\n", c.Terminal.TerminalId)
	}
	return ""
}