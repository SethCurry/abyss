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
					proxiedAgent.Initialize(context.Background(), asInit)
				}
				return nil
			},
			IsRPC: true,
		},
	}
}
