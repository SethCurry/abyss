package wsacp

import (
	"context"

	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
)

func NewProxiedACPClient(router *wsrouter.ACPRouter) *ProxiedACPClient {
	return &ProxiedACPClient{
		router: router,
	}
}

type ProxiedACPClient struct {
	router *wsrouter.ACPRouter
}

func (e *ProxiedACPClient) RequestPermission(ctx context.Context, params acp.RequestPermissionRequest) (acp.RequestPermissionResponse, error) {
	return proxyRoundTrip[acp.RequestPermissionRequest, acp.RequestPermissionResponse](e.router, params)
}

func (e *ProxiedACPClient) SessionUpdate(ctx context.Context, params acp.SessionNotification) error {
	return e.router.Send(params)
}

func (e *ProxiedACPClient) WriteTextFile(ctx context.Context, params acp.WriteTextFileRequest) (acp.WriteTextFileResponse, error) {
	return proxyRoundTrip[acp.WriteTextFileRequest, acp.WriteTextFileResponse](e.router, params)
}

func (e *ProxiedACPClient) ReadTextFile(ctx context.Context, params acp.ReadTextFileRequest) (acp.ReadTextFileResponse, error) {
	return proxyRoundTrip[acp.ReadTextFileRequest, acp.ReadTextFileResponse](e.router, params)
}

func (e *ProxiedACPClient) CreateTerminal(ctx context.Context, params acp.CreateTerminalRequest) (acp.CreateTerminalResponse, error) {
	return proxyRoundTrip[acp.CreateTerminalRequest, acp.CreateTerminalResponse](e.router, params)
}

func (e *ProxiedACPClient) TerminalOutput(ctx context.Context, params acp.TerminalOutputRequest) (acp.TerminalOutputResponse, error) {
	return proxyRoundTrip[acp.TerminalOutputRequest, acp.TerminalOutputResponse](e.router, params)
}

func (e *ProxiedACPClient) ReleaseTerminal(ctx context.Context, params acp.ReleaseTerminalRequest) (acp.ReleaseTerminalResponse, error) {
	return proxyRoundTrip[acp.ReleaseTerminalRequest, acp.ReleaseTerminalResponse](e.router, params)
}

func (e *ProxiedACPClient) WaitForTerminalExit(ctx context.Context, params acp.WaitForTerminalExitRequest) (acp.WaitForTerminalExitResponse, error) {
	return proxyRoundTrip[acp.WaitForTerminalExitRequest, acp.WaitForTerminalExitResponse](e.router, params)
}

// KillTerminal implements acp.Client.
func (e *ProxiedACPClient) KillTerminal(ctx context.Context, params acp.KillTerminalRequest) (acp.KillTerminalResponse, error) {
	return proxyRoundTrip[acp.KillTerminalRequest, acp.KillTerminalResponse](e.router, params)
}
