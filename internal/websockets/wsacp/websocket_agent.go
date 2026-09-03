package wsacp

import (
	"context"
	"os"

	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

type WebsocketAgent struct {
	logger     zerolog.Logger
	conn       *acp.AgentSideConnection
	underlying *ProxiedACPAgent
	router     *wsrouter.ACPRouter
}

var (
	_ acp.Agent             = (*WebsocketAgent)(nil)
	_ acp.AgentLoader       = (*WebsocketAgent)(nil)
	_ acp.AgentExperimental = (*WebsocketAgent)(nil)
)

// NewWebsocketAgent creates an agent-side ACP proxy that bridges a websocket
// connection to a client over stdio.
func NewWebsocketAgent(underlying *ProxiedACPAgent, logger zerolog.Logger) *WebsocketAgent {
	return &WebsocketAgent{
		logger:     logger,
		underlying: underlying,
	}
}

// SetSessionMode implements acp.Agent.
func (w *WebsocketAgent) SetSessionMode(ctx context.Context, params acp.SetSessionModeRequest) (acp.SetSessionModeResponse, error) {
	w.logger.Debug().Str("method", "SetSessionMode").Msg("handling request")
	return w.underlying.SetSessionMode(ctx, params)
}

// UnstableForkSession implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableForkSession(ctx context.Context, params acp.UnstableForkSessionRequest) (acp.UnstableForkSessionResponse, error) {
	w.logger.Debug().Str("method", "UnstableForkSession").Msg("handling request")
	return w.underlying.UnstableForkSession(ctx, params)
}

// ListSessions implements acp.Agent.
func (w *WebsocketAgent) ListSessions(ctx context.Context, params acp.ListSessionsRequest) (acp.ListSessionsResponse, error) {
	w.logger.Debug().Str("method", "ListSessions").Msg("handling request")
	return w.underlying.ListSessions(ctx, params)
}

// ResumeSession implements acp.Agent.
func (w *WebsocketAgent) ResumeSession(ctx context.Context, params acp.ResumeSessionRequest) (acp.ResumeSessionResponse, error) {
	w.logger.Debug().Str("method", "ResumeSession").Msg("handling request")
	return w.underlying.ResumeSession(ctx, params)
}

// SetSessionConfigOption implements acp.Agent.
func (w *WebsocketAgent) SetSessionConfigOption(ctx context.Context, params acp.SetSessionConfigOptionRequest) (acp.SetSessionConfigOptionResponse, error) {
	w.logger.Debug().Str("method", "SetSessionConfigOption").Msg("handling request")
	return w.underlying.SetSessionConfigOption(ctx, params)
}

// UnstableDidChangeDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidChangeDocument(ctx context.Context, params acp.UnstableDidChangeDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidChangeDocument").Msg("handling notification")
	return w.underlying.UnstableDidChangeDocument(ctx, params)
}

// UnstableDidCloseDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidCloseDocument(ctx context.Context, params acp.UnstableDidCloseDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidCloseDocument").Msg("handling notification")
	return w.underlying.UnstableDidCloseDocument(ctx, params)
}

// UnstableDidFocusDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidFocusDocument(ctx context.Context, params acp.UnstableDidFocusDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidFocusDocument").Msg("handling notification")
	return w.underlying.UnstableDidFocusDocument(ctx, params)
}

// UnstableDidOpenDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidOpenDocument(ctx context.Context, params acp.UnstableDidOpenDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidOpenDocument").Msg("handling notification")
	return w.underlying.UnstableDidOpenDocument(ctx, params)
}

// UnstableDidSaveDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidSaveDocument(ctx context.Context, params acp.UnstableDidSaveDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidSaveDocument").Msg("handling notification")
	return w.underlying.UnstableDidSaveDocument(ctx, params)
}

// Logout implements acp.Agent.
func (w *WebsocketAgent) Logout(ctx context.Context, params acp.LogoutRequest) (acp.LogoutResponse, error) {
	w.logger.Debug().Str("method", "Logout").Msg("handling request")
	return w.underlying.Logout(ctx, params)
}

// UnstableAcceptNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableAcceptNes(ctx context.Context, params acp.UnstableAcceptNesNotification) error {
	w.logger.Debug().Str("method", "UnstableAcceptNes").Msg("handling notification")
	return w.underlying.UnstableAcceptNes(ctx, params)
}

// UnstableCloseNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableCloseNes(ctx context.Context, params acp.UnstableCloseNesRequest) (acp.UnstableCloseNesResponse, error) {
	w.logger.Debug().Str("method", "UnstableCloseNes").Msg("handling request")
	return w.underlying.UnstableCloseNes(ctx, params)
}

// UnstableRejectNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableRejectNes(ctx context.Context, params acp.UnstableRejectNesNotification) error {
	w.logger.Debug().Str("method", "UnstableRejectNes").Msg("handling notification")
	return w.underlying.UnstableRejectNes(ctx, params)
}

// UnstableStartNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableStartNes(ctx context.Context, params acp.UnstableStartNesRequest) (acp.UnstableStartNesResponse, error) {
	w.logger.Debug().Str("method", "UnstableStartNes").Msg("handling request")
	return w.underlying.UnstableStartNes(ctx, params)
}

// UnstableSuggestNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableSuggestNes(ctx context.Context, params acp.UnstableSuggestNesRequest) (acp.UnstableSuggestNesResponse, error) {
	w.logger.Debug().Str("method", "UnstableSuggestNes").Msg("handling request")
	return w.underlying.UnstableSuggestNes(ctx, params)
}

// UnstableDisableProvider implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDisableProvider(ctx context.Context, params acp.UnstableDisableProviderRequest) (acp.UnstableDisableProviderResponse, error) {
	w.logger.Debug().Str("method", "UnstableDisableProvider").Msg("handling request")
	return w.underlying.UnstableDisableProvider(ctx, params)
}

// UnstableListProviders implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableListProviders(ctx context.Context, params acp.UnstableListProvidersRequest) (acp.UnstableListProvidersResponse, error) {
	w.logger.Debug().Str("method", "UnstableListProviders").Msg("handling request")
	return w.underlying.UnstableListProviders(ctx, params)
}

// UnstableSetProvider implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableSetProvider(ctx context.Context, params acp.UnstableSetProviderRequest) (acp.UnstableSetProviderResponse, error) {
	w.logger.Debug().Str("method", "UnstableSetProvider").Msg("handling request")
	return w.underlying.UnstableSetProvider(ctx, params)
}

// UnstableDeleteSession implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDeleteSession(ctx context.Context, params acp.UnstableDeleteSessionRequest) (acp.UnstableDeleteSessionResponse, error) {
	w.logger.Debug().Str("method", "UnstableDeleteSession").Msg("handling request")
	return w.underlying.UnstableDeleteSession(ctx, params)
}

// CloseSession implements acp.Agent.
func (w *WebsocketAgent) CloseSession(ctx context.Context, params acp.CloseSessionRequest) (acp.CloseSessionResponse, error) {
	w.logger.Debug().Str("method", "CloseSession").Msg("handling request")
	return w.underlying.CloseSession(ctx, params)
}

// SetAgentConnection stores the ACP agent-side connection used to forward
// client capability requests received over the websocket to the client over
// stdio.
func (w *WebsocketAgent) SetAgentConnection(conn *acp.AgentSideConnection) {
	w.logger.Debug().Msg("agent connection set")
	w.conn = conn
	w.router.SetClient(conn)
}

func (w *WebsocketAgent) Initialize(ctx context.Context, params acp.InitializeRequest) (acp.InitializeResponse, error) {
	w.logger.Debug().Str("method", "Initialize").Msg("handling request")
	return w.underlying.Initialize(ctx, params)
}

func (w *WebsocketAgent) NewSession(ctx context.Context, params acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	if _, err := os.Stat(params.Cwd); err != nil {
		err = os.MkdirAll(params.Cwd, 0755)
		if err != nil {
			w.logger.Error().Err(err).Str("method", "NewSession").Msg("failed to create cwd")
			return acp.NewSessionResponse{}, err
		}
	}
	w.logger.Debug().Str("method", "NewSession").Msg("handling request")
	return w.underlying.NewSession(ctx, params)
}

func (w *WebsocketAgent) Authenticate(ctx context.Context, params acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	w.logger.Debug().Str("method", "Authenticate").Msg("handling request")
	return w.underlying.Authenticate(ctx, params)
}

func (w *WebsocketAgent) LoadSession(ctx context.Context, params acp.LoadSessionRequest) (acp.LoadSessionResponse, error) {
	w.logger.Debug().Str("method", "LoadSession").Msg("handling request")
	return w.underlying.LoadSession(ctx, params)
}

func (w *WebsocketAgent) Cancel(ctx context.Context, params acp.CancelNotification) error {
	w.logger.Debug().Str("method", "Cancel").Msg("handling notification")
	return w.underlying.Cancel(ctx, params)
}

func (w *WebsocketAgent) Prompt(ctx context.Context, params acp.PromptRequest) (acp.PromptResponse, error) {
	w.logger.Debug().Str("method", "Prompt").Msg("handling request")
	return w.underlying.Prompt(ctx, params)
}
