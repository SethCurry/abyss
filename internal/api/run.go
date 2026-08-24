package api

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// RunClient dials the websocket server at wsURL and bridges it to a client
// (typically an editor) over stdio.
func RunClient(ctx context.Context, wsURL string, logger zerolog.Logger) error {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial Docker websocket: %w", err)
	}
	defer conn.Close()

	agent := NewWebsocketAgent(conn, logger)
	asc := acp.NewAgentSideConnection(agent, os.Stdout, os.Stdin)
	asc.SetLogger(slog.Default())
	agent.SetAgentConnection(asc)

	// Run the demultiplexing read loop, forwarding client capability requests
	// from the websocket to the client over stdio.
	go func() {
		if err := agent.Serve(ctx); err != nil {
			logger.Error().Err(err).Msg("websocket bridge failed")
		}
		_ = conn.Close()
	}()

	<-asc.Done()
	logger.Info().Msg("agent closed connection")
	return nil
}
