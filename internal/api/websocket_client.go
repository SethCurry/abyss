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

type WebsocketAgentClient struct {
	logger  zerolog.Logger
	ws      *wsConn
	acpConn *acp.ClientSideConnection
}

var _ acp.Client = (*WebsocketAgentClient)(nil)

// NewWebsocketAgentClient creates a client-side ACP proxy that bridges a
// websocket connection to an agent over stdio.
func NewWebsocketAgentClient(conn *websocket.Conn, logger zerolog.Logger) *WebsocketAgentClient {
	return &WebsocketAgentClient{
		logger: logger,
		ws:     newWSConn(conn, logger),
	}
}

// SetClientConnection stores the ACP client-side connection used to forward
// agent requests received over the websocket to the agent over stdio.
func (e *WebsocketAgentClient) SetClientConnection(conn *acp.ClientSideConnection) {
	e.acpConn = conn
}

func (e *WebsocketAgentClient) RequestPermission(ctx context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	e.logger.Debug().Str("method", "RequestPermission").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.RequestPermissionResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "RequestPermission").Msg("request failed")
		return acp.RequestPermissionResponse{}, err
	}
	return *resp.(*acp.RequestPermissionResponse), nil
}

func (e *WebsocketAgentClient) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	e.logger.Debug().Str("method", "SessionUpdate").Msg("handling notification")
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	if err := e.ws.notify(wsmessage.SessionNotificationType, data); err != nil {
		e.logger.Error().Err(err).Str("method", "SessionUpdate").Msg("notification failed")
		return err
	}
	return nil
}

func (e *WebsocketAgentClient) WriteTextFile(ctx context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	e.logger.Debug().Str("method", "WriteTextFile").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.WriteTextFileResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "WriteTextFile").Msg("request failed")
		return acp.WriteTextFileResponse{}, err
	}
	return *resp.(*acp.WriteTextFileResponse), nil
}

func (e *WebsocketAgentClient) ReadTextFile(ctx context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	e.logger.Debug().Str("method", "ReadTextFile").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.ReadTextFileResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "ReadTextFile").Msg("request failed")
		return acp.ReadTextFileResponse{}, err
	}
	return *resp.(*acp.ReadTextFileResponse), nil
}

func (e *WebsocketAgentClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	e.logger.Debug().Str("method", "CreateTerminal").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.CreateTerminalResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "CreateTerminal").Msg("request failed")
		return acp.CreateTerminalResponse{}, err
	}
	return *resp.(*acp.CreateTerminalResponse), nil
}

func (e *WebsocketAgentClient) TerminalOutput(ctx context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	e.logger.Debug().Str("method", "TerminalOutput").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.TerminalOutputResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "TerminalOutput").Msg("request failed")
		return acp.TerminalOutputResponse{}, err
	}
	return *resp.(*acp.TerminalOutputResponse), nil
}

func (e *WebsocketAgentClient) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	e.logger.Debug().Str("method", "ReleaseTerminal").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.ReleaseTerminalResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "ReleaseTerminal").Msg("request failed")
		return acp.ReleaseTerminalResponse{}, err
	}
	return *resp.(*acp.ReleaseTerminalResponse), nil
}

func (e *WebsocketAgentClient) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	e.logger.Debug().Str("method", "WaitForTerminalExit").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.WaitForTerminalExitResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "WaitForTerminalExit").Msg("request failed")
		return acp.WaitForTerminalExitResponse{}, err
	}
	return *resp.(*acp.WaitForTerminalExitResponse), nil
}

// KillTerminal implements acp.Client.
func (e *WebsocketAgentClient) KillTerminal(ctx context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	e.logger.Debug().Str("method", "KillTerminal").Msg("handling request")
	resp, err := e.ws.request(ctx, messageTypeFor(params), params, wsmessage.KillTerminalResponseType)
	if err != nil {
		e.logger.Error().Err(err).Str("method", "KillTerminal").Msg("request failed")
		return acp.KillTerminalResponse{}, err
	}
	return *resp.(*acp.KillTerminalResponse), nil
}

// Serve runs the demultiplexing read loop, forwarding agent requests received
// over the websocket to the agent over stdio.
func (e *WebsocketAgentClient) Serve(ctx context.Context) error {
	return e.ws.serve(ctx, e.handleIncoming)
}

// handleIncoming forwards an agent request or notification received over the
// websocket to the agent over stdio.
func (e *WebsocketAgentClient) handleIncoming(ctx context.Context, mt wsmessage.MessageType, id uint32, data []byte) error {
	switch mt {
	case wsmessage.InitializeRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.InitializeResponseType, e.acpConn.Initialize)
	case wsmessage.NewSessionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.NewSessionResponseType, e.acpConn.NewSession)
	case wsmessage.PromptRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.PromptResponseType, e.acpConn.Prompt)
	case wsmessage.AuthenticateRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.AuthenticateResponseType, e.acpConn.Authenticate)
	case wsmessage.LogoutRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.LogoutResponseType, e.acpConn.Logout)
	case wsmessage.CloseSessionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.CloseSessionResponseType, e.acpConn.CloseSession)
	case wsmessage.ListSessionsRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.ListSessionsResponseType, e.acpConn.ListSessions)
	case wsmessage.LoadSessionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.LoadSessionResponseType, e.acpConn.LoadSession)
	case wsmessage.ResumeSessionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.ResumeSessionResponseType, e.acpConn.ResumeSession)
	case wsmessage.SetSessionConfigOptionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.SetSessionConfigOptionResponseType, e.acpConn.SetSessionConfigOption)
	case wsmessage.SetSessionModeRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.SetSessionModeResponseType, e.acpConn.SetSessionMode)
	case wsmessage.UnstableForkSessionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableForkSessionResponseType, e.acpConn.UnstableForkSession)
	case wsmessage.UnstableDeleteSessionRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableDeleteSessionResponseType, e.acpConn.UnstableDeleteSession)
	case wsmessage.UnstableSetProviderRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableSetProviderResponseType, e.acpConn.UnstableSetProvider)
	case wsmessage.UnstableListProvidersRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableListProvidersResponseType, e.acpConn.UnstableListProviders)
	case wsmessage.UnstableDisableProviderRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableDisableProviderResponseType, e.acpConn.UnstableDisableProvider)
	case wsmessage.UnstableSuggestNesRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableSuggestNesResponseType, e.acpConn.UnstableSuggestNes)
	case wsmessage.UnstableStartNesRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableStartNesResponseType, e.acpConn.UnstableStartNes)
	case wsmessage.UnstableCloseNesRequestType:
		return serveRequest(e.ws, ctx, id, data, wsmessage.UnstableCloseNesResponseType, e.acpConn.UnstableCloseNes)
	case wsmessage.CancelNotificationType:
		return serveNotification(ctx, data, e.acpConn.Cancel)
	case wsmessage.UnstableAcceptNesNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableAcceptNes)
	case wsmessage.UnstableRejectNesNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableRejectNes)
	case wsmessage.UnstableDidChangeDocumentNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableDidChangeDocument)
	case wsmessage.UnstableDidCloseDocumentNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableDidCloseDocument)
	case wsmessage.UnstableDidFocusDocumentNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableDidFocusDocument)
	case wsmessage.UnstableDidOpenDocumentNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableDidOpenDocument)
	case wsmessage.UnstableDidSaveDocumentNotificationType:
		return serveNotification(ctx, data, e.acpConn.UnstableDidSaveDocument)
	default:
		return fmt.Errorf("unexpected message type on websocket: %d", mt)
	}
}

// serveRequest unmarshals a request, forwards it to the peer over stdio, and
// writes the response back over the websocket with the original correlation ID.
func serveRequest[T any, R any](w *wsConn, ctx context.Context, id uint32, data []byte, respType wsmessage.MessageType, fn func(context.Context, T) (R, error)) error {
	var req T
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	resp, err := fn(ctx, req)
	if err != nil {
		return err
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	return w.writeMessage(respType, id, respData)
}

// serveNotification unmarshals a notification and forwards it to the peer over
// stdio.
func serveNotification[T any](ctx context.Context, data []byte, fn func(context.Context, T) error) error {
	var req T
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	return fn(ctx, req)
}
