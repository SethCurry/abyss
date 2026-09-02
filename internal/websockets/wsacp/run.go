package wsacp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RunClient dials the websocket server at wsURL and bridges it to a client
// (typically an editor) over stdio.
func RunClient(ctx context.Context, wsURL string, logger zerolog.Logger) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial Docker websocket: %w", err)
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			logger.Warn().Err(closeErr).Msg("failed to close client websocket connection")
		}
	}()
	acpConnChan := make(chan wsrouter.SocketMessage, 10)

	socket := wsrouter.NewSocket(conn, map[int]chan wsrouter.SocketMessage{
		1: acpConnChan,
	})

	wsConn := wsrouter.NewConn(socket, acpConnChan)

	router := wsrouter.NewACPRouter(wsConn)

	proxiedAgent := NewProxiedACPAgent(router)

	agent := NewWebsocketAgent(proxiedAgent, logger)
	asc := acp.NewAgentSideConnection(agent, os.Stdout, os.Stdin)
	asc.SetLogger(slog.Default())
	agent.SetAgentConnection(asc)

	// Run the demultiplexing read loop, forwarding client capability requests
	// from the websocket to the client over stdio.
	go func() {
		//nolint:staticcheck
		router.Serve()
		_ = conn.Close()
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
			Handler: func(router *wsrouter.ACPRouter, msg *protobyss.Container) any {
				creator := wsmessage.MessageTypeToMessage[wsmessage.MessageType(msg.TypeId)]
				newValue := creator()

				err := json.Unmarshal(msg.Content, &newValue)
				if err != nil {
					log.Warn().Err(err).Msg("failed to unmarshal message")
				}

				if asInit, ok := newValue.(acp.InitializeRequest); ok {
					proxiedAgent.Initialize(context.Background(), asInit)
				}
				return nil
			},
			IsRPC: true,
		},
	}
}
