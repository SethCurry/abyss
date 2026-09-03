package wsacp

import (
	"context"

	"github.com/SethCurry/abyss/internal/acptools"
	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

type WebsocketAgentClient struct {
	logger        zerolog.Logger
	underlying    *ProxiedACPClient
	acpConn       *acp.ClientSideConnection
	terminalTools *acptools.TerminalTools
	fileTools     *acptools.FilesystemTools
	router        *wsrouter.ACPRouter
}

var _ acp.Client = (*WebsocketAgentClient)(nil)

// NewWebsocketAgentClient creates a client-side ACP proxy that bridges a
// websocket connection to an agent over stdio.
func NewWebsocketAgentClient(underlying *ProxiedACPClient, terminalTools *acptools.TerminalTools, fileTools *acptools.FilesystemTools, logger zerolog.Logger) *WebsocketAgentClient {
	return &WebsocketAgentClient{
		logger:        logger,
		underlying:    underlying,
		terminalTools: terminalTools,
		fileTools:     fileTools,
	}
}

// SetClientConnection stores the ACP client-side connection used to forward
// agent requests received over the websocket to the agent over stdio.
func (e *WebsocketAgentClient) SetClientConnection(conn *acp.ClientSideConnection) {
	e.acpConn = conn
	e.router.SetAgent(conn)
}

func (e *WebsocketAgentClient) RequestPermission(ctx context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	e.logger.Debug().Str("method", "RequestPermission").Msg("handling request")
	return e.underlying.RequestPermission(ctx, params)
}

func (e *WebsocketAgentClient) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	e.logger.Debug().Str("method", "SessionUpdate").Msg("handling notification")
	return e.underlying.SessionUpdate(ctx, params)
}

func (e *WebsocketAgentClient) WriteTextFile(ctx context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	if e.fileTools != nil {
		e.logger.Debug().Str("method", "WriteTextFile").Str("handler", "fileTools").Msg("handling request")
		return e.fileTools.WriteTextFile(ctx, params)
	}

	e.logger.Debug().Str("method", "WriteTextFile").Str("handler", "client").Msg("handling request")
	return e.underlying.WriteTextFile(ctx, params)
}

func (e *WebsocketAgentClient) ReadTextFile(ctx context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	if e.fileTools != nil {
		e.logger.Debug().Str("method", "ReadTextFile").Str("handler", "fileTools").Msg("handling request")
		return e.fileTools.ReadTextFile(ctx, params)
	}

	e.logger.Debug().Str("method", "ReadTextFile").Str("handler", "client").Msg("handling request")
	return e.underlying.ReadTextFile(ctx, params)
}

func (e *WebsocketAgentClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	e.logger.Debug().Str("method", "CreateTerminal").Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().Str("method", "CreateTerminal").Str("handler", "terminalTools").Msg("handling request")
		return e.terminalTools.CreateTerminal(ctx, params)
	}

	return e.underlying.CreateTerminal(ctx, params)
}

func (e *WebsocketAgentClient) TerminalOutput(ctx context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	e.logger.Debug().Str("method", "TerminalOutput").Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().Str("method", "TerminalOutput").Str("handler", "terminalTools").Msg("handling request")
		return e.terminalTools.TerminalOutput(ctx, params)
	}

	return e.underlying.TerminalOutput(ctx, params)
}

func (e *WebsocketAgentClient) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	e.logger.Debug().Str("method", "ReleaseTerminal").Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().Str("method", "ReleaseTerminal").Str("handler", "terminalTools").Msg("handling request")
		return e.terminalTools.ReleaseTerminal(ctx, params)
	}

	return e.underlying.ReleaseTerminal(ctx, params)
}

func (e *WebsocketAgentClient) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	e.logger.Debug().Str("method", "WaitForTerminalExit").Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().Str("method", "WaitForTerminalExit").Str("handler", "terminalTools").Msg("handling request")
		return e.terminalTools.WaitForTerminalExit(ctx, params)
	}

	return e.underlying.WaitForTerminalExit(ctx, params)
}

// KillTerminal implements acp.Client.
func (e *WebsocketAgentClient) KillTerminal(ctx context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	e.logger.Debug().Str("method", "KillTerminal").Msg("handling request")

	if e.terminalTools != nil {
		e.logger.Debug().Str("method", "KillTerminal").Str("handler", "terminalTools").Msg("handling request")
		return e.terminalTools.KillTerminal(ctx, params)
	}

	return e.underlying.KillTerminal(ctx, params)
}
