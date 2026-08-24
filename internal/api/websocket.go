package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// messageTypeFor returns the MessageType corresponding to the given ACP
// message struct.
func messageTypeFor[T wsmessage.ACPMessageType](msg T) wsmessage.MessageType {
	switch any(msg).(type) {
	case acp.RequestPermissionRequest:
		return wsmessage.RequestPermissionRequestType
	case acp.RequestPermissionResponse:
		return wsmessage.RequestPermissionResponseType
	case acp.WriteTextFileRequest:
		return wsmessage.WriteTextFileRequestType
	case acp.WriteTextFileResponse:
		return wsmessage.WriteTextFileResponseType
	case acp.ReadTextFileRequest:
		return wsmessage.ReadTextFileRequestType
	case acp.ReadTextFileResponse:
		return wsmessage.ReadTextFileResponseType
	case acp.CreateTerminalRequest:
		return wsmessage.CreateTerminalRequestType
	case acp.CreateTerminalResponse:
		return wsmessage.CreateTerminalResponseType
	case acp.TerminalOutputRequest:
		return wsmessage.TerminalOutputRequestType
	case acp.TerminalOutputResponse:
		return wsmessage.TerminalOutputResponseType
	case acp.ReleaseTerminalRequest:
		return wsmessage.ReleaseTerminalRequestType
	case acp.ReleaseTerminalResponse:
		return wsmessage.ReleaseTerminalResponseType
	case acp.WaitForTerminalExitRequest:
		return wsmessage.WaitForTerminalExitRequestType
	case acp.WaitForTerminalExitResponse:
		return wsmessage.WaitForTerminalExitResponseType
	case acp.KillTerminalRequest:
		return wsmessage.KillTerminalRequestType
	case acp.KillTerminalResponse:
		return wsmessage.KillTerminalResponseType
	case acp.SessionNotification:
		return wsmessage.SessionNotificationType
	case acp.SetSessionModeRequest:
		return wsmessage.SetSessionModeRequestType
	case acp.SetSessionModeResponse:
		return wsmessage.SetSessionModeResponseType
	case acp.UnstableForkSessionRequest:
		return wsmessage.UnstableForkSessionRequestType
	case acp.UnstableForkSessionResponse:
		return wsmessage.UnstableForkSessionResponseType
	case acp.ListSessionsRequest:
		return wsmessage.ListSessionsRequestType
	case acp.ListSessionsResponse:
		return wsmessage.ListSessionsResponseType
	case acp.ResumeSessionRequest:
		return wsmessage.ResumeSessionRequestType
	case acp.ResumeSessionResponse:
		return wsmessage.ResumeSessionResponseType
	case acp.SetSessionConfigOptionRequest:
		return wsmessage.SetSessionConfigOptionRequestType
	case acp.SetSessionConfigOptionResponse:
		return wsmessage.SetSessionConfigOptionResponseType
	case acp.LogoutRequest:
		return wsmessage.LogoutRequestType
	case acp.LogoutResponse:
		return wsmessage.LogoutResponseType
	case acp.UnstableCloseNesRequest:
		return wsmessage.UnstableCloseNesRequestType
	case acp.UnstableCloseNesResponse:
		return wsmessage.UnstableCloseNesResponseType
	case acp.UnstableStartNesRequest:
		return wsmessage.UnstableStartNesRequestType
	case acp.UnstableStartNesResponse:
		return wsmessage.UnstableStartNesResponseType
	case acp.UnstableSuggestNesRequest:
		return wsmessage.UnstableSuggestNesRequestType
	case acp.UnstableSuggestNesResponse:
		return wsmessage.UnstableSuggestNesResponseType
	case acp.UnstableAcceptNesNotification:
		return wsmessage.UnstableAcceptNesNotificationType
	case acp.UnstableRejectNesNotification:
		return wsmessage.UnstableRejectNesNotificationType
	case acp.UnstableDidChangeDocumentNotification:
		return wsmessage.UnstableDidChangeDocumentNotificationType
	case acp.UnstableDidCloseDocumentNotification:
		return wsmessage.UnstableDidCloseDocumentNotificationType
	case acp.UnstableDidFocusDocumentNotification:
		return wsmessage.UnstableDidFocusDocumentNotificationType
	case acp.UnstableDidOpenDocumentNotification:
		return wsmessage.UnstableDidOpenDocumentNotificationType
	case acp.UnstableDidSaveDocumentNotification:
		return wsmessage.UnstableDidSaveDocumentNotificationType
	case acp.UnstableDisableProviderRequest:
		return wsmessage.UnstableDisableProviderRequestType
	case acp.UnstableDisableProviderResponse:
		return wsmessage.UnstableDisableProviderResponseType
	case acp.UnstableListProvidersRequest:
		return wsmessage.UnstableListProvidersRequestType
	case acp.UnstableListProvidersResponse:
		return wsmessage.UnstableListProvidersResponseType
	case acp.UnstableSetProviderRequest:
		return wsmessage.UnstableSetProviderRequestType
	case acp.UnstableSetProviderResponse:
		return wsmessage.UnstableSetProviderResponseType
	case acp.UnstableDeleteSessionRequest:
		return wsmessage.UnstableDeleteSessionRequestType
	case acp.UnstableDeleteSessionResponse:
		return wsmessage.UnstableDeleteSessionResponseType
	case acp.CloseSessionRequest:
		return wsmessage.CloseSessionRequestType
	case acp.CloseSessionResponse:
		return wsmessage.CloseSessionResponseType
	case acp.InitializeRequest:
		return wsmessage.InitializeRequestType
	case acp.InitializeResponse:
		return wsmessage.InitializeResponseType
	case acp.NewSessionRequest:
		return wsmessage.NewSessionRequestType
	case acp.NewSessionResponse:
		return wsmessage.NewSessionResponseType
	case acp.AuthenticateRequest:
		return wsmessage.AuthenticateRequestType
	case acp.AuthenticateResponse:
		return wsmessage.AuthenticateResponseType
	case acp.LoadSessionRequest:
		return wsmessage.LoadSessionRequestType
	case acp.LoadSessionResponse:
		return wsmessage.LoadSessionResponseType
	case acp.PromptRequest:
		return wsmessage.PromptRequestType
	case acp.PromptResponse:
		return wsmessage.PromptResponseType
	case acp.CancelNotification:
		return wsmessage.CancelNotificationType
	}

	return wsmessage.MessageTypeNotExist
}

type ACPWebsocketServer struct {
	upgrader websocket.Upgrader
	logger   zerolog.Logger
}

// NewACPWebsocketServer creates a websocket server that bridges incoming
// websocket connections to freshly spawned agent processes over stdio.
func NewACPWebsocketServer(logger zerolog.Logger) *ACPWebsocketServer {
	return &ACPWebsocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		logger: logger,
	}
}

// Serve listens on addr and bridges each websocket connection to an agent
// process spawned from agentCommand.
func (s *ACPWebsocketServer) Serve(ctx context.Context, addr string, agentCommand []string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.handleConnection(ctx, w, r, agentCommand)
	})

	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	s.logger.Info().Str("addr", addr).Msg("websocket server listening")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *ACPWebsocketServer) handleConnection(ctx context.Context, w http.ResponseWriter, r *http.Request, agentCommand []string) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to upgrade to websocket")
		return
	}
	defer conn.Close()

	if len(agentCommand) == 0 {
		s.logger.Error().Msg("no agent command configured")
		return
	}

	s.logger.Info().Strs("agent_command", agentCommand).Msg("spawning underlying agent")

	args := append([]string{"-c"}, agentCommand...)

	cmd := exec.CommandContext(ctx, "bash", args...)

	cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to open agent stdin")
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to open agent stdout")
		return
	}
	if err := cmd.Start(); err != nil {
		s.logger.Error().Err(err).Msg("failed to start agent")
		return
	}
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	client := NewWebsocketAgentClient(conn, s.logger)
	csc := acp.NewClientSideConnection(client, stdin, stdout)
	csc.SetLogger(slog.Default())
	client.SetClientConnection(csc)

	go func() {
		if err := client.Serve(ctx); err != nil {
			s.logger.Error().Err(err).Msg("websocket bridge failed")
		}
		_ = conn.Close()
	}()

	<-csc.Done()
}
