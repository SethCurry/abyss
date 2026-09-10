package wsacp

import (
	"context"
	"os"
	"time"

	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
	"github.com/rs/zerolog"
)

// HostProxy is the entrypoint for the ACP client like Zed.
// It stores the Websocket-proxied connection to the real agent as
// well as the real connection to the ACP client.
type HostProxy struct {
	logger zerolog.Logger

	// conn is the actual client connected over stdio
	conn *acp.AgentSideConnection

	// underlying is the websocket proxy to the container side of abyss
	underlying *ProxiedACPAgent

	// router is used for inbound messages
	router *wsrouter.ACPRouter

	// TODO we shouldn't have both underlying and router here
	// The abstraction is leaking
}

var (
	_ acp.Agent             = (*HostProxy)(nil)
	_ acp.AgentLoader       = (*HostProxy)(nil)
	_ acp.AgentExperimental = (*HostProxy)(nil)
)

// NewHostProxy creates a new host-side proxy that accepts ACP input
// over stdio, and communicates to the container-side proxy via websocket.
func NewHostProxy(underlying *ProxiedACPAgent, router *wsrouter.ACPRouter, logger zerolog.Logger) *HostProxy {
	return &HostProxy{
		logger:     logger,
		underlying: underlying,
		router:     router,
	}
}

func (w *HostProxy) UserMessage(ctx context.Context, sessionID string, message string) error {
	return w.conn.SessionUpdate(ctx, acp.SessionNotification{
		SessionId: acp.SessionId(sessionID),
		Update: acp.SessionUpdate{
			AgentMessageChunk: &acp.SessionUpdateAgentMessageChunk{
				Content: acp.TextBlock(message),
			},
		},
	})
}

// SetSessionMode implements acp.Agent.
func (w *HostProxy) SetSessionMode(ctx context.Context, params acp.SetSessionModeRequest) (acp.SetSessionModeResponse, error) {
	w.logger.Debug().
		Str("method", "SetSessionMode").
		Str("session_id", string(params.SessionId)).
		Str("mode_id", string(params.ModeId)).
		Msg("handling request")
	return w.underlying.SetSessionMode(ctx, params)
}

// UnstableForkSession implements acp.AgentExperimental.
func (w *HostProxy) UnstableForkSession(ctx context.Context, params acp.UnstableForkSessionRequest) (acp.UnstableForkSessionResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableForkSession").
		Str("session_id", string(params.SessionId)).
		Str("cwd", params.Cwd).
		Msg("handling request")
	return w.underlying.UnstableForkSession(ctx, params)
}

// ListSessions implements acp.Agent.
func (w *HostProxy) ListSessions(ctx context.Context, params acp.ListSessionsRequest) (acp.ListSessionsResponse, error) {
	w.logger.Debug().
		Str("method", "ListSessions").
		Msg("handling request")
	return w.underlying.ListSessions(ctx, params)
}

// ResumeSession implements acp.Agent.
func (w *HostProxy) ResumeSession(ctx context.Context, params acp.ResumeSessionRequest) (acp.ResumeSessionResponse, error) {
	w.logger.Debug().
		Str("method", "ResumeSession").
		Str("session_id", string(params.SessionId)).
		Str("cwd", params.Cwd).
		Msg("handling request")
	return w.underlying.ResumeSession(ctx, params)
}

// SetSessionConfigOption implements acp.Agent.
func (w *HostProxy) SetSessionConfigOption(ctx context.Context, params acp.SetSessionConfigOptionRequest) (acp.SetSessionConfigOptionResponse, error) {
	w.logger.Debug().
		Str("method", "SetSessionConfigOption").
		Msg("handling request")
	return w.underlying.SetSessionConfigOption(ctx, params)
}

// UnstableDidChangeDocument implements acp.AgentExperimental.
func (w *HostProxy) UnstableDidChangeDocument(ctx context.Context, params acp.UnstableDidChangeDocumentNotification) error {
	w.logger.Debug().
		Str("method", "UnstableDidChangeDocument").
		Str("session_id", string(params.SessionId)).
		Str("uri", string(params.Uri)).
		Msg("handling notification")
	return w.underlying.UnstableDidChangeDocument(ctx, params)
}

// UnstableDidCloseDocument implements acp.AgentExperimental.
func (w *HostProxy) UnstableDidCloseDocument(ctx context.Context, params acp.UnstableDidCloseDocumentNotification) error {
	w.logger.Debug().
		Str("method", "UnstableDidCloseDocument").
		Str("session_id", string(params.SessionId)).
		Str("uri", string(params.Uri)).
		Msg("handling notification")
	return w.underlying.UnstableDidCloseDocument(ctx, params)
}

// UnstableDidFocusDocument implements acp.AgentExperimental.
func (w *HostProxy) UnstableDidFocusDocument(ctx context.Context, params acp.UnstableDidFocusDocumentNotification) error {
	w.logger.Debug().
		Str("method", "UnstableDidFocusDocument").
		Str("session_id", string(params.SessionId)).
		Str("uri", string(params.Uri)).
		Msg("handling notification")
	return w.underlying.UnstableDidFocusDocument(ctx, params)
}

// UnstableDidOpenDocument implements acp.AgentExperimental.
func (w *HostProxy) UnstableDidOpenDocument(ctx context.Context, params acp.UnstableDidOpenDocumentNotification) error {
	w.logger.Debug().
		Str("method", "UnstableDidOpenDocument").
		Str("session_id", string(params.SessionId)).
		Str("uri", string(params.Uri)).
		Msg("handling notification")
	return w.underlying.UnstableDidOpenDocument(ctx, params)
}

// UnstableDidSaveDocument implements acp.AgentExperimental.
func (w *HostProxy) UnstableDidSaveDocument(ctx context.Context, params acp.UnstableDidSaveDocumentNotification) error {
	w.logger.Debug().
		Str("method", "UnstableDidSaveDocument").
		Str("session_id", string(params.SessionId)).
		Str("uri", string(params.Uri)).
		Msg("handling notification")
	return w.underlying.UnstableDidSaveDocument(ctx, params)
}

// Logout implements acp.Agent.
func (w *HostProxy) Logout(ctx context.Context, params acp.LogoutRequest) (acp.LogoutResponse, error) {
	w.logger.Debug().
		Str("method", "Logout").
		Msg("handling request")
	return w.underlying.Logout(ctx, params)
}

// UnstableAcceptNes implements acp.AgentExperimental.
func (w *HostProxy) UnstableAcceptNes(ctx context.Context, params acp.UnstableAcceptNesNotification) error {
	w.logger.Debug().
		Str("method", "UnstableAcceptNes").
		Msg("handling notification")
	return w.underlying.UnstableAcceptNes(ctx, params)
}

// UnstableCloseNes implements acp.AgentExperimental.
func (w *HostProxy) UnstableCloseNes(ctx context.Context, params acp.UnstableCloseNesRequest) (acp.UnstableCloseNesResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableCloseNes").
		Msg("handling request")
	return w.underlying.UnstableCloseNes(ctx, params)
}

// UnstableRejectNes implements acp.AgentExperimental.
func (w *HostProxy) UnstableRejectNes(ctx context.Context, params acp.UnstableRejectNesNotification) error {
	w.logger.Debug().
		Str("method", "UnstableRejectNes").
		Msg("handling notification")
	return w.underlying.UnstableRejectNes(ctx, params)
}

// UnstableStartNes implements acp.AgentExperimental.
func (w *HostProxy) UnstableStartNes(ctx context.Context, params acp.UnstableStartNesRequest) (acp.UnstableStartNesResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableStartNes").
		Msg("handling request")
	return w.underlying.UnstableStartNes(ctx, params)
}

// UnstableSuggestNes implements acp.AgentExperimental.
func (w *HostProxy) UnstableSuggestNes(ctx context.Context, params acp.UnstableSuggestNesRequest) (acp.UnstableSuggestNesResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableSuggestNes").
		Msg("handling request")
	return w.underlying.UnstableSuggestNes(ctx, params)
}

// UnstableDisableProvider implements acp.AgentExperimental.
func (w *HostProxy) UnstableDisableProvider(ctx context.Context, params acp.UnstableDisableProviderRequest) (acp.UnstableDisableProviderResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableDisableProvider").
		Msg("handling request")
	return w.underlying.UnstableDisableProvider(ctx, params)
}

// UnstableListProviders implements acp.AgentExperimental.
func (w *HostProxy) UnstableListProviders(ctx context.Context, params acp.UnstableListProvidersRequest) (acp.UnstableListProvidersResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableListProviders").
		Msg("handling request")
	return w.underlying.UnstableListProviders(ctx, params)
}

// UnstableSetProvider implements acp.AgentExperimental.
func (w *HostProxy) UnstableSetProvider(ctx context.Context, params acp.UnstableSetProviderRequest) (acp.UnstableSetProviderResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableSetProvider").
		Msg("handling request")
	return w.underlying.UnstableSetProvider(ctx, params)
}

// UnstableDeleteSession implements acp.AgentExperimental.
func (w *HostProxy) UnstableDeleteSession(ctx context.Context, params acp.UnstableDeleteSessionRequest) (acp.UnstableDeleteSessionResponse, error) {
	w.logger.Debug().
		Str("method", "UnstableDeleteSession").
		Str("session_id", string(params.SessionId)).
		Msg("handling request")
	return w.underlying.UnstableDeleteSession(ctx, params)
}

// CloseSession implements acp.Agent.
func (w *HostProxy) CloseSession(ctx context.Context, params acp.CloseSessionRequest) (acp.CloseSessionResponse, error) {
	w.logger.Debug().
		Str("method", "CloseSession").
		Str("session_id", string(params.SessionId)).
		Msg("handling request")
	return w.underlying.CloseSession(ctx, params)
}

// SetAgentConnection stores the ACP agent-side connection used to forward
// client capability requests received over the websocket to the client over
// stdio.
func (w *HostProxy) SetAgentConnection(conn *acp.AgentSideConnection) {
	w.logger.Debug().Msg("agent connection set")
	w.conn = conn
	w.router.SetClient(conn)
}

func (w *HostProxy) Initialize(ctx context.Context, params acp.InitializeRequest) (acp.InitializeResponse, error) {
	w.logger.Debug().
		Str("method", "Initialize").
		Msg("handling request")
	return w.underlying.Initialize(ctx, params)
}

func (w *HostProxy) NewSession(ctx context.Context, params acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	if _, err := os.Stat(params.Cwd); err != nil {
		err = os.MkdirAll(params.Cwd, 0755)
		if err != nil {
			w.logger.Error().
				Err(err).
				Str("method", "NewSession").
				Msg("failed to create cwd")
			return acp.NewSessionResponse{}, err
		}
	}
	w.logger.Debug().
		Str("method", "NewSession").
		Msg("handling request")
	newSession, err := w.underlying.NewSession(ctx, params)
	if err != nil {
		return acp.NewSessionResponse{}, err
	}

	go func() {
		time.Sleep(time.Second * 1)

		_ = w.UserMessage(context.Background(), string(newSession.SessionId), "Welcome to abyss!")
	}()
	return newSession, nil
}

func (w *HostProxy) Authenticate(ctx context.Context, params acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	w.logger.Debug().
		Str("method", "Authenticate").
		Msg("handling request")
	return w.underlying.Authenticate(ctx, params)
}

func (w *HostProxy) LoadSession(ctx context.Context, params acp.LoadSessionRequest) (acp.LoadSessionResponse, error) {
	w.logger.Debug().
		Str("method", "LoadSession").
		Str("session_id", string(params.SessionId)).
		Msg("handling request")
	return w.underlying.LoadSession(ctx, params)
}

func (w *HostProxy) Cancel(ctx context.Context, params acp.CancelNotification) error {
	w.logger.Debug().
		Str("method", "Cancel").
		Str("session_id", string(params.SessionId)).
		Msg("handling notification")
	return w.underlying.Cancel(ctx, params)
}

func (w *HostProxy) Prompt(ctx context.Context, params acp.PromptRequest) (acp.PromptResponse, error) {
	w.logger.Debug().
		Str("method", "Prompt").
		Msg("handling request")
	return w.underlying.Prompt(ctx, params)
}
