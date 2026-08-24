package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type WebsocketAgent struct {
	logger zerolog.Logger
	conn   *acp.AgentSideConnection
	ws     *wsConn
}

var (
	_ acp.Agent             = (*WebsocketAgent)(nil)
	_ acp.AgentLoader       = (*WebsocketAgent)(nil)
	_ acp.AgentExperimental = (*WebsocketAgent)(nil)
)

// NewWebsocketAgent creates an agent-side ACP proxy that bridges a websocket
// connection to a client over stdio.
func NewWebsocketAgent(conn *websocket.Conn, logger zerolog.Logger) *WebsocketAgent {
	return &WebsocketAgent{
		logger: logger,
		ws:     newWSConn(conn, logger),
	}
}

// roundTrip sends a request over the websocket and waits for the matching
// response, returning the unmarshaled response value.
func roundTrip[T wsmessage.ACPMessageType, R wsmessage.ACPMessageType](w *WebsocketAgent, ctx context.Context, req T, respType wsmessage.MessageType) (R, error) {
	var zero R
	resp, err := w.ws.request(ctx, messageTypeFor(req), req, respType)
	if err != nil {
		return zero, err
	}
	return *resp.(*R), nil
}

// notify sends a notification over the websocket without waiting for a
// response.
func notify[T wsmessage.ACPMessageType](w *WebsocketAgent, ctx context.Context, req T) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	return w.ws.notify(messageTypeFor(req), data)
}

// SetSessionMode implements acp.Agent.
func (w *WebsocketAgent) SetSessionMode(ctx context.Context, params acp.SetSessionModeRequest) (acp.SetSessionModeResponse, error) {
	w.logger.Debug().Str("method", "SetSessionMode").Msg("handling request")
	resp, err := roundTrip[acp.SetSessionModeRequest, acp.SetSessionModeResponse](w, ctx, params, wsmessage.SetSessionModeResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "SetSessionMode").Msg("request failed")
	}
	return resp, err
}

// UnstableForkSession implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableForkSession(ctx context.Context, params acp.UnstableForkSessionRequest) (acp.UnstableForkSessionResponse, error) {
	w.logger.Debug().Str("method", "UnstableForkSession").Msg("handling request")
	resp, err := roundTrip[acp.UnstableForkSessionRequest, acp.UnstableForkSessionResponse](w, ctx, params, wsmessage.UnstableForkSessionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableForkSession").Msg("request failed")
	}
	return resp, err
}

// ListSessions implements acp.Agent.
func (w *WebsocketAgent) ListSessions(ctx context.Context, params acp.ListSessionsRequest) (acp.ListSessionsResponse, error) {
	w.logger.Debug().Str("method", "ListSessions").Msg("handling request")
	resp, err := roundTrip[acp.ListSessionsRequest, acp.ListSessionsResponse](w, ctx, params, wsmessage.ListSessionsResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "ListSessions").Msg("request failed")
	}
	return resp, err
}

// ResumeSession implements acp.Agent.
func (w *WebsocketAgent) ResumeSession(ctx context.Context, params acp.ResumeSessionRequest) (acp.ResumeSessionResponse, error) {
	w.logger.Debug().Str("method", "ResumeSession").Msg("handling request")
	resp, err := roundTrip[acp.ResumeSessionRequest, acp.ResumeSessionResponse](w, ctx, params, wsmessage.ResumeSessionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "ResumeSession").Msg("request failed")
	}
	return resp, err
}

// SetSessionConfigOption implements acp.Agent.
func (w *WebsocketAgent) SetSessionConfigOption(ctx context.Context, params acp.SetSessionConfigOptionRequest) (acp.SetSessionConfigOptionResponse, error) {
	w.logger.Debug().Str("method", "SetSessionConfigOption").Msg("handling request")
	resp, err := roundTrip[acp.SetSessionConfigOptionRequest, acp.SetSessionConfigOptionResponse](w, ctx, params, wsmessage.SetSessionConfigOptionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "SetSessionConfigOption").Msg("request failed")
	}
	return resp, err
}

// UnstableDidChangeDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidChangeDocument(ctx context.Context, params acp.UnstableDidChangeDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidChangeDocument").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDidChangeDocument").Msg("notification failed")
	}
	return err
}

// UnstableDidCloseDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidCloseDocument(ctx context.Context, params acp.UnstableDidCloseDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidCloseDocument").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDidCloseDocument").Msg("notification failed")
	}
	return err
}

// UnstableDidFocusDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidFocusDocument(ctx context.Context, params acp.UnstableDidFocusDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidFocusDocument").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDidFocusDocument").Msg("notification failed")
	}
	return err
}

// UnstableDidOpenDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidOpenDocument(ctx context.Context, params acp.UnstableDidOpenDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidOpenDocument").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDidOpenDocument").Msg("notification failed")
	}
	return err
}

// UnstableDidSaveDocument implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDidSaveDocument(ctx context.Context, params acp.UnstableDidSaveDocumentNotification) error {
	w.logger.Debug().Str("method", "UnstableDidSaveDocument").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDidSaveDocument").Msg("notification failed")
	}
	return err
}

// Logout implements acp.Agent.
func (w *WebsocketAgent) Logout(ctx context.Context, params acp.LogoutRequest) (acp.LogoutResponse, error) {
	w.logger.Debug().Str("method", "Logout").Msg("handling request")
	resp, err := roundTrip[acp.LogoutRequest, acp.LogoutResponse](w, ctx, params, wsmessage.LogoutResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "Logout").Msg("request failed")
	}
	return resp, err
}

// UnstableAcceptNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableAcceptNes(ctx context.Context, params acp.UnstableAcceptNesNotification) error {
	w.logger.Debug().Str("method", "UnstableAcceptNes").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableAcceptNes").Msg("notification failed")
	}
	return err
}

// UnstableCloseNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableCloseNes(ctx context.Context, params acp.UnstableCloseNesRequest) (acp.UnstableCloseNesResponse, error) {
	w.logger.Debug().Str("method", "UnstableCloseNes").Msg("handling request")
	resp, err := roundTrip[acp.UnstableCloseNesRequest, acp.UnstableCloseNesResponse](w, ctx, params, wsmessage.UnstableCloseNesResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableCloseNes").Msg("request failed")
	}
	return resp, err
}

// UnstableRejectNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableRejectNes(ctx context.Context, params acp.UnstableRejectNesNotification) error {
	w.logger.Debug().Str("method", "UnstableRejectNes").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableRejectNes").Msg("notification failed")
	}
	return err
}

// UnstableStartNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableStartNes(ctx context.Context, params acp.UnstableStartNesRequest) (acp.UnstableStartNesResponse, error) {
	w.logger.Debug().Str("method", "UnstableStartNes").Msg("handling request")
	resp, err := roundTrip[acp.UnstableStartNesRequest, acp.UnstableStartNesResponse](w, ctx, params, wsmessage.UnstableStartNesResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableStartNes").Msg("request failed")
	}
	return resp, err
}

// UnstableSuggestNes implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableSuggestNes(ctx context.Context, params acp.UnstableSuggestNesRequest) (acp.UnstableSuggestNesResponse, error) {
	w.logger.Debug().Str("method", "UnstableSuggestNes").Msg("handling request")
	resp, err := roundTrip[acp.UnstableSuggestNesRequest, acp.UnstableSuggestNesResponse](w, ctx, params, wsmessage.UnstableSuggestNesResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableSuggestNes").Msg("request failed")
	}
	return resp, err
}

// UnstableDisableProvider implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDisableProvider(ctx context.Context, params acp.UnstableDisableProviderRequest) (acp.UnstableDisableProviderResponse, error) {
	w.logger.Debug().Str("method", "UnstableDisableProvider").Msg("handling request")
	resp, err := roundTrip[acp.UnstableDisableProviderRequest, acp.UnstableDisableProviderResponse](w, ctx, params, wsmessage.UnstableDisableProviderResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDisableProvider").Msg("request failed")
	}
	return resp, err
}

// UnstableListProviders implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableListProviders(ctx context.Context, params acp.UnstableListProvidersRequest) (acp.UnstableListProvidersResponse, error) {
	w.logger.Debug().Str("method", "UnstableListProviders").Msg("handling request")
	resp, err := roundTrip[acp.UnstableListProvidersRequest, acp.UnstableListProvidersResponse](w, ctx, params, wsmessage.UnstableListProvidersResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableListProviders").Msg("request failed")
	}
	return resp, err
}

// UnstableSetProvider implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableSetProvider(ctx context.Context, params acp.UnstableSetProviderRequest) (acp.UnstableSetProviderResponse, error) {
	w.logger.Debug().Str("method", "UnstableSetProvider").Msg("handling request")
	resp, err := roundTrip[acp.UnstableSetProviderRequest, acp.UnstableSetProviderResponse](w, ctx, params, wsmessage.UnstableSetProviderResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableSetProvider").Msg("request failed")
	}
	return resp, err
}

// UnstableDeleteSession implements acp.AgentExperimental.
func (w *WebsocketAgent) UnstableDeleteSession(ctx context.Context, params acp.UnstableDeleteSessionRequest) (acp.UnstableDeleteSessionResponse, error) {
	w.logger.Debug().Str("method", "UnstableDeleteSession").Msg("handling request")
	resp, err := roundTrip[acp.UnstableDeleteSessionRequest, acp.UnstableDeleteSessionResponse](w, ctx, params, wsmessage.UnstableDeleteSessionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "UnstableDeleteSession").Msg("request failed")
	}
	return resp, err
}

// CloseSession implements acp.Agent.
func (w *WebsocketAgent) CloseSession(ctx context.Context, params acp.CloseSessionRequest) (acp.CloseSessionResponse, error) {
	w.logger.Debug().Str("method", "CloseSession").Msg("handling request")
	resp, err := roundTrip[acp.CloseSessionRequest, acp.CloseSessionResponse](w, ctx, params, wsmessage.CloseSessionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "CloseSession").Msg("request failed")
	}
	return resp, err
}

// SetAgentConnection stores the ACP agent-side connection used to forward
// client capability requests received over the websocket to the client over
// stdio.
func (w *WebsocketAgent) SetAgentConnection(conn *acp.AgentSideConnection) {
	w.logger.Debug().Msg("agent connection set")
	w.conn = conn
}

func (w *WebsocketAgent) Initialize(ctx context.Context, params acp.InitializeRequest) (acp.InitializeResponse, error) {
	w.logger.Debug().Str("method", "Initialize").Msg("handling request")
	resp, err := roundTrip[acp.InitializeRequest, acp.InitializeResponse](w, ctx, params, wsmessage.InitializeResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "Initialize").Msg("request failed")
	}
	return resp, err
}

func (w *WebsocketAgent) NewSession(ctx context.Context, params acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	w.logger.Debug().Str("method", "NewSession").Msg("handling request")
	resp, err := roundTrip[acp.NewSessionRequest, acp.NewSessionResponse](w, ctx, params, wsmessage.NewSessionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "NewSession").Msg("request failed")
	}
	return resp, err
}

func (w *WebsocketAgent) Authenticate(ctx context.Context, params acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	w.logger.Debug().Str("method", "Authenticate").Msg("handling request")
	resp, err := roundTrip[acp.AuthenticateRequest, acp.AuthenticateResponse](w, ctx, params, wsmessage.AuthenticateResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "Authenticate").Msg("request failed")
	}
	return resp, err
}

func (w *WebsocketAgent) LoadSession(ctx context.Context, params acp.LoadSessionRequest) (acp.LoadSessionResponse, error) {
	w.logger.Debug().Str("method", "LoadSession").Msg("handling request")
	resp, err := roundTrip[acp.LoadSessionRequest, acp.LoadSessionResponse](w, ctx, params, wsmessage.LoadSessionResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "LoadSession").Msg("request failed")
	}
	return resp, err
}

func (w *WebsocketAgent) Cancel(ctx context.Context, params acp.CancelNotification) error {
	w.logger.Debug().Str("method", "Cancel").Msg("handling notification")
	err := notify(w, ctx, params)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "Cancel").Msg("notification failed")
	}
	return err
}

func (w *WebsocketAgent) Prompt(ctx context.Context, params acp.PromptRequest) (acp.PromptResponse, error) {
	w.logger.Debug().Str("method", "Prompt").Msg("handling request")
	resp, err := roundTrip[acp.PromptRequest, acp.PromptResponse](w, ctx, params, wsmessage.PromptResponseType)
	if err != nil {
		w.logger.Error().Err(err).Str("method", "Prompt").Msg("request failed")
	}
	return resp, err
}

// Serve runs the demultiplexing read loop, forwarding client capability
// requests received over the websocket to the client over stdio.
func (w *WebsocketAgent) Serve(ctx context.Context) error {
	return w.ws.serve(ctx, w.handleIncoming)
}

// handleIncoming forwards a client capability request or notification received
// over the websocket to the client over stdio.
func (w *WebsocketAgent) handleIncoming(ctx context.Context, mt wsmessage.MessageType, id uint32, data []byte) error {
	switch mt {
	case wsmessage.RequestPermissionRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.RequestPermissionResponseType, w.conn.RequestPermission)
	case wsmessage.WriteTextFileRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.WriteTextFileResponseType, w.conn.WriteTextFile)
	case wsmessage.ReadTextFileRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.ReadTextFileResponseType, w.conn.ReadTextFile)
	case wsmessage.CreateTerminalRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.CreateTerminalResponseType, w.conn.CreateTerminal)
	case wsmessage.TerminalOutputRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.TerminalOutputResponseType, w.conn.TerminalOutput)
	case wsmessage.ReleaseTerminalRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.ReleaseTerminalResponseType, w.conn.ReleaseTerminal)
	case wsmessage.WaitForTerminalExitRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.WaitForTerminalExitResponseType, w.conn.WaitForTerminalExit)
	case wsmessage.KillTerminalRequestType:
		return serveRequest(w.ws, ctx, id, data, wsmessage.KillTerminalResponseType, w.conn.KillTerminal)
	case wsmessage.SessionNotificationType:
		return serveNotification(ctx, data, w.conn.SessionUpdate)
	default:
		return fmt.Errorf("unexpected message type on websocket: %d", mt)
	}
}
