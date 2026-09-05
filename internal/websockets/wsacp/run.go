package wsacp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/SethCurry/abyss/internal/acp/termacp"
	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Oneshot(ctx context.Context, prompt string, wsURL string, tlsConfig *tls.Config, logger zerolog.Logger) error {
	dialer := websocket.DefaultDialer
	if tlsConfig != nil {
		dialer = &websocket.Dialer{TLSClientConfig: tlsConfig}
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial Docker websocket: %w", err)
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Warn().Err(closeErr).Msg("failed to close client websocket connection")
		}
	}()
	socket := wsrouter.NewProtoRouter()
	router := wsrouter.NewACPRouter()
	acpConn := wsrouter.NewACPConn(socket, router.ServeMessage)
	router.SetConn(acpConn)
	socket.Handle(1, acpConn.Handle)

	termACPClient := termacp.NewTermACPClient()
	router.SetClient(termACPClient)
	proxiedAgent := NewProxiedACPAgent(router)

	go func() {
		socket.Serve(conn)
	}()

	cwd, err := os.Getwd()
	if err != nil {
		logger.Error().Err(err).Msg("failed to get current working directory")
		return err
	}

	newSession, err := proxiedAgent.NewSession(ctx, acp.NewSessionRequest{
		AdditionalDirectories: []string{},
		McpServers:            []acp.McpServer{},
		Cwd:                   cwd,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create new session")
	}

	_, err = proxiedAgent.Prompt(ctx, acp.PromptRequest{
		SessionId: newSession.SessionId,
		Prompt: []acp.ContentBlock{
			acp.TextBlock(prompt),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// RunClient dials the websocket server at wsURL and bridges it to a client
// (typically an editor) over stdio. A non-nil tlsConfig enables TLS for the
// connection.
func RunClient(ctx context.Context, wsURL string, tlsConfig *tls.Config, logger zerolog.Logger) error {
	dialer := websocket.DefaultDialer
	if tlsConfig != nil {
		dialer = &websocket.Dialer{TLSClientConfig: tlsConfig}
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial Docker websocket: %w", err)
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Warn().Err(closeErr).Msg("failed to close client websocket connection")
		}
	}()
	socket := wsrouter.NewProtoRouter()
	router := wsrouter.NewACPRouter()
	acpConn := wsrouter.NewACPConn(socket, router.ServeMessage)
	router.SetConn(acpConn)
	socket.Handle(1, acpConn.Handle)

	proxiedAgent := NewProxiedACPAgent(router)

	agent := NewWebsocketAgent(proxiedAgent, router, logger)
	asc := acp.NewAgentSideConnection(agent, os.Stdout, os.Stdin)
	asc.SetLogger(slog.Default())
	agent.SetAgentConnection(asc)

	// Run the demultiplexing read loop, forwarding client capability requests
	// from the websocket to the client over stdio.
	go func() {
		//nolint:staticcheck
		socket.Serve(conn)
	}()

	<-asc.Done()
	logger.Info().Msg("agent closed connection")
	return nil
}

func GenerateClientRoutes(proxiedAgent *WebsocketAgent) []wsrouter.MessageType {
	return []wsrouter.MessageType{
		{
			ID:   int32(wsmessage.InitializeRequestType),
			Type: reflect.TypeOf(acp.InitializeRequest{}),
			Handler: func(router *wsrouter.ACPRouter, msg *protobyss.ACPContainer) any {
				creator := wsmessage.MessageTypeToMessage[wsmessage.MessageType(msg.TypeId)]
				newValue := creator()

				err := json.Unmarshal(msg.Content, &newValue)
				if err != nil {
					log.Warn().Err(err).Msg("failed to unmarshal message")
				}

				if asInit, ok := newValue.(acp.InitializeRequest); ok {
					_, _ = proxiedAgent.Initialize(context.Background(), asInit)
				}
				return nil
			},
			IsRPC: true,
		},
	}
}
