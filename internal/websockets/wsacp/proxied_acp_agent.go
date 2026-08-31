package wsacp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
)

func proxyRoundTrip[T, P any](router *wsrouter.Router, params T) (P, error) {
	var ret P

	prom, err := router.Request(params)
	if err != nil {
		return ret, fmt.Errorf("failed to make websocket RPC request: %s", err)
	}

	container := prom.Wait()

	err = json.Unmarshal(container.Content, &ret)
	if err != nil {
		return ret, fmt.Errorf("failed to unmarshal response content: %w", err)
	}

	return ret, nil
}

func NewProxiedACPAgent(router *wsrouter.Router) *ProxiedACPAgent {
	return &ProxiedACPAgent{
		router: router,
	}
}

type ProxiedACPAgent struct {
	router *wsrouter.Router
}

// SetSessionMode implements acp.Agent.
func (w *ProxiedACPAgent) SetSessionMode(ctx context.Context, params acp.SetSessionModeRequest) (acp.SetSessionModeResponse, error) {
	return proxyRoundTrip[acp.SetSessionModeRequest, acp.SetSessionModeResponse](w.router, params)
}

// UnstableForkSession implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableForkSession(ctx context.Context, params acp.UnstableForkSessionRequest) (acp.UnstableForkSessionResponse, error) {
	return proxyRoundTrip[acp.UnstableForkSessionRequest, acp.UnstableForkSessionResponse](w.router, params)
}

// ListSessions implements acp.Agent.
func (w *ProxiedACPAgent) ListSessions(ctx context.Context, params acp.ListSessionsRequest) (acp.ListSessionsResponse, error) {
	return proxyRoundTrip[acp.ListSessionsRequest, acp.ListSessionsResponse](w.router, params)
}

// ResumeSession implements acp.Agent.
func (w *ProxiedACPAgent) ResumeSession(ctx context.Context, params acp.ResumeSessionRequest) (acp.ResumeSessionResponse, error) {
	return proxyRoundTrip[acp.ResumeSessionRequest, acp.ResumeSessionResponse](w.router, params)
}

// SetSessionConfigOption implements acp.Agent.
func (w *ProxiedACPAgent) SetSessionConfigOption(ctx context.Context, params acp.SetSessionConfigOptionRequest) (acp.SetSessionConfigOptionResponse, error) {
	return proxyRoundTrip[acp.SetSessionConfigOptionRequest, acp.SetSessionConfigOptionResponse](w.router, params)
}

// UnstableDidChangeDocument implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDidChangeDocument(ctx context.Context, params acp.UnstableDidChangeDocumentNotification) error {
	return w.router.Send(params)
}

// UnstableDidCloseDocument implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDidCloseDocument(ctx context.Context, params acp.UnstableDidCloseDocumentNotification) error {
	return w.router.Send(params)
}

// UnstableDidFocusDocument implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDidFocusDocument(ctx context.Context, params acp.UnstableDidFocusDocumentNotification) error {
	return w.router.Send(params)
}

// UnstableDidOpenDocument implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDidOpenDocument(ctx context.Context, params acp.UnstableDidOpenDocumentNotification) error {
	return w.router.Send(params)
}

// UnstableDidSaveDocument implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDidSaveDocument(ctx context.Context, params acp.UnstableDidSaveDocumentNotification) error {
	return w.router.Send(params)
}

// Logout implements acp.Agent.
func (w *ProxiedACPAgent) Logout(ctx context.Context, params acp.LogoutRequest) (acp.LogoutResponse, error) {
	return proxyRoundTrip[acp.LogoutRequest, acp.LogoutResponse](w.router, params)
}

// UnstableAcceptNes implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableAcceptNes(ctx context.Context, params acp.UnstableAcceptNesNotification) error {
	return w.router.Send(params)
}

// UnstableCloseNes implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableCloseNes(ctx context.Context, params acp.UnstableCloseNesRequest) (acp.UnstableCloseNesResponse, error) {
	return proxyRoundTrip[acp.UnstableCloseNesRequest, acp.UnstableCloseNesResponse](w.router, params)
}

// UnstableRejectNes implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableRejectNes(ctx context.Context, params acp.UnstableRejectNesNotification) error {
	return w.router.Send(params)
}

// UnstableStartNes implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableStartNes(ctx context.Context, params acp.UnstableStartNesRequest) (acp.UnstableStartNesResponse, error) {
	return proxyRoundTrip[acp.UnstableStartNesRequest, acp.UnstableStartNesResponse](w.router, params)
}

// UnstableSuggestNes implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableSuggestNes(ctx context.Context, params acp.UnstableSuggestNesRequest) (acp.UnstableSuggestNesResponse, error) {
	return proxyRoundTrip[acp.UnstableSuggestNesRequest, acp.UnstableSuggestNesResponse](w.router, params)
}

// UnstableDisableProvider implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDisableProvider(ctx context.Context, params acp.UnstableDisableProviderRequest) (acp.UnstableDisableProviderResponse, error) {
	return proxyRoundTrip[acp.UnstableDisableProviderRequest, acp.UnstableDisableProviderResponse](w.router, params)
}

// UnstableListProviders implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableListProviders(ctx context.Context, params acp.UnstableListProvidersRequest) (acp.UnstableListProvidersResponse, error) {
	return proxyRoundTrip[acp.UnstableListProvidersRequest, acp.UnstableListProvidersResponse](w.router, params)
}

// UnstableSetProvider implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableSetProvider(ctx context.Context, params acp.UnstableSetProviderRequest) (acp.UnstableSetProviderResponse, error) {
	return proxyRoundTrip[acp.UnstableSetProviderRequest, acp.UnstableSetProviderResponse](w.router, params)
}

// UnstableDeleteSession implements acp.AgentExperimental.
func (w *ProxiedACPAgent) UnstableDeleteSession(ctx context.Context, params acp.UnstableDeleteSessionRequest) (acp.UnstableDeleteSessionResponse, error) {
	return proxyRoundTrip[acp.UnstableDeleteSessionRequest, acp.UnstableDeleteSessionResponse](w.router, params)
}

// CloseSession implements acp.Agent.
func (w *ProxiedACPAgent) CloseSession(ctx context.Context, params acp.CloseSessionRequest) (acp.CloseSessionResponse, error) {
	return proxyRoundTrip[acp.CloseSessionRequest, acp.CloseSessionResponse](w.router, params)
}

func (w *ProxiedACPAgent) Initialize(ctx context.Context, params acp.InitializeRequest) (acp.InitializeResponse, error) {
	return proxyRoundTrip[acp.InitializeRequest, acp.InitializeResponse](w.router, params)
}

func (w *ProxiedACPAgent) NewSession(ctx context.Context, params acp.NewSessionRequest) (acp.NewSessionResponse, error) {
	return proxyRoundTrip[acp.NewSessionRequest, acp.NewSessionResponse](w.router, params)
}

func (w *ProxiedACPAgent) Authenticate(ctx context.Context, params acp.AuthenticateRequest) (acp.AuthenticateResponse, error) {
	return proxyRoundTrip[acp.AuthenticateRequest, acp.AuthenticateResponse](w.router, params)
}

func (w *ProxiedACPAgent) LoadSession(ctx context.Context, params acp.LoadSessionRequest) (acp.LoadSessionResponse, error) {
	return proxyRoundTrip[acp.LoadSessionRequest, acp.LoadSessionResponse](w.router, params)
}

func (w *ProxiedACPAgent) Cancel(ctx context.Context, params acp.CancelNotification) error {
	return w.router.Send(params)
}

func (w *ProxiedACPAgent) Prompt(ctx context.Context, params acp.PromptRequest) (acp.PromptResponse, error) {
	return proxyRoundTrip[acp.PromptRequest, acp.PromptResponse](w.router, params)
}
