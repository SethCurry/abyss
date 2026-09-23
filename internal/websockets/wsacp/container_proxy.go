// Package wsacp implements WebSocket proxies for the Agent Client Protocol.
package wsacp

import (
	"context"

	"github.com/SethCurry/abyss/internal/acp/acptools"
	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

// ContainerProxy forwards ACP requests received over a websocket to the
// agent, dispatching terminal and filesystem operations to in-process tools
// when available.
type ContainerProxy struct {
	logger        zerolog.Logger
	underlying    *ProxiedACPClient
	acpConn       *acp.ClientSideConnection
	terminalTools *acptools.TerminalTools
	fileTools     *acptools.FilesystemTools
	router        *wsrouter.ACPRouter
}

var _ acp.Client = (*ContainerProxy)(nil)

// NewContainerProxy constructs a ContainerProxy from its dependencies.
func NewContainerProxy(
	underlying *ProxiedACPClient,
	router *wsrouter.ACPRouter,
	terminalTools *acptools.TerminalTools,
	fileTools *acptools.FilesystemTools,
	logger zerolog.Logger) *ContainerProxy {
	return &ContainerProxy{
		logger:        logger,
		underlying:    underlying,
		terminalTools: terminalTools,
		fileTools:     fileTools,
		router:        router,
	}
}

// SetClientConnection stores the ACP client-side connection used to forward
// agent requests received over the websocket to the agent over stdio.
func (e *ContainerProxy) SetClientConnection(conn *acp.ClientSideConnection) {
	e.acpConn = conn
	e.router.SetAgent(conn)
}

// RequestPermission forwards permission requests to the underlying client.
func (e *ContainerProxy) RequestPermission(
	ctx context.Context, params acp.RequestPermissionRequest,
) (acp.RequestPermissionResponse, error) {
	e.logger.Debug().
		Str("method", "RequestPermission").
		Msg("handling request")
	return e.underlying.RequestPermission(ctx, params)
}

// SessionUpdate forwards session update notifications to the underlying client.
func (e *ContainerProxy) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	e.logger.Debug().
		Str("method", "SessionUpdate").
		Msg("handling notification")
	return e.underlying.SessionUpdate(ctx, params)
}

// WriteTextFile writes a text file via the in-process file tools when
// available, falling back to the underlying client otherwise.
func (e *ContainerProxy) WriteTextFile(
	ctx context.Context,
	params acp.WriteTextFileRequest,
) (acp.WriteTextFileResponse, error) {
	if e.fileTools != nil {
		e.logger.Debug().
			Str("method", "WriteTextFile").
			Str("handler", "fileTools").
			Msg("handling request")
		return e.fileTools.WriteTextFile(ctx, params)
	}

	e.logger.Debug().
		Str("method", "WriteTextFile").
		Str("handler", "client").
		Msg("handling request")
	return e.underlying.WriteTextFile(ctx, params)
}

// ReadTextFile reads a text file via the in-process file tools when
// available, falling back to the underlying client otherwise.
func (e *ContainerProxy) ReadTextFile(
	ctx context.Context, params acp.ReadTextFileRequest,
) (acp.ReadTextFileResponse, error) {
	if e.fileTools != nil {
		e.logger.Debug().
			Str("method", "ReadTextFile").
			Str("handler", "fileTools").
			Msg("handling request")
		return e.fileTools.ReadTextFile(ctx, params)
	}

	e.logger.Debug().
		Str("method", "ReadTextFile").
		Str("handler", "client").
		Msg("handling request")
	return e.underlying.ReadTextFile(ctx, params)
}

// CreateTerminal creates a terminal via the in-process terminal tools when
// available, falling back to the underlying client otherwise.
func (e *ContainerProxy) CreateTerminal(
	ctx context.Context, params acp.CreateTerminalRequest,
) (acp.CreateTerminalResponse, error) {
	e.logger.Debug().
		Str("method", "CreateTerminal").
		Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().
			Str("method", "CreateTerminal").
			Str("handler", "terminalTools").
			Msg("handling request")
		return e.terminalTools.CreateTerminal(ctx, params)
	}

	return e.underlying.CreateTerminal(ctx, params)
}

// TerminalOutput fetches terminal output via the in-process terminal tools
// when available, falling back to the underlying client otherwise.
func (e *ContainerProxy) TerminalOutput(
	ctx context.Context, params acp.TerminalOutputRequest,
) (acp.TerminalOutputResponse, error) {
	e.logger.Debug().
		Str("method", "TerminalOutput").
		Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().
			Str("method", "TerminalOutput").
			Str("handler", "terminalTools").
			Msg("handling request")
		return e.terminalTools.TerminalOutput(ctx, params)
	}

	return e.underlying.TerminalOutput(ctx, params)
}

// ReleaseTerminal releases a terminal via the in-process terminal tools when
// available, falling back to the underlying client otherwise.
func (e *ContainerProxy) ReleaseTerminal(
	ctx context.Context, params acp.ReleaseTerminalRequest,
) (acp.ReleaseTerminalResponse, error) {
	e.logger.Debug().
		Str("method", "ReleaseTerminal").
		Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().
			Str("method", "ReleaseTerminal").
			Str("handler", "terminalTools").
			Msg("handling request")
		return e.terminalTools.ReleaseTerminal(ctx, params)
	}

	return e.underlying.ReleaseTerminal(ctx, params)
}

// WaitForTerminalExit waits for terminal exit via the in-process terminal
// tools when available, falling back to the underlying client otherwise.
func (e *ContainerProxy) WaitForTerminalExit(
	ctx context.Context, params acp.WaitForTerminalExitRequest,
) (acp.WaitForTerminalExitResponse, error) {
	e.logger.Debug().
		Str("method", "WaitForTerminalExit").
		Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().
			Str("method", "WaitForTerminalExit").
			Str("handler", "terminalTools").
			Msg("handling request")
		return e.terminalTools.WaitForTerminalExit(ctx, params)
	}

	return e.underlying.WaitForTerminalExit(ctx, params)
}

// KillTerminal implements acp.Client.
func (e *ContainerProxy) KillTerminal(
	ctx context.Context, params acp.KillTerminalRequest,
) (acp.KillTerminalResponse, error) {
	e.logger.Debug().
		Str("method", "KillTerminal").
		Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().
			Str("method", "KillTerminal").
			Str("handler", "terminalTools").
			Msg("handling request")
		return e.terminalTools.KillTerminal(ctx, params)
	}

	return e.underlying.KillTerminal(ctx, params)
}
