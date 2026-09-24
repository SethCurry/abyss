package wsacp

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"

	"github.com/SethCurry/abyss/internal/acptools"
	"github.com/SethCurry/abyss/internal/agentconfig"
	"github.com/SethCurry/abyss/internal/fp"
	"github.com/SethCurry/abyss/internal/plugin"
	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

// dialAndServe dials the websocket server at wsURL and wires up the socket,
// router, and ACP connection. It launches the demultiplexing read loop in a
// goroutine and returns the constructed proxied agent along with the close
// behavior. A non-nil tlsConfig enables TLS for the connection.
func dialAndServe(
	ctx context.Context,
	wsURL string,
	tlsConfig *tls.Config,
	plugins *plugin.ACPManager,
	logger zerolog.Logger) (*websocket.Conn, *wsrouter.ProtoRouter, *ProxiedACPAgent, error) {
	dialer := websocket.DefaultDialer
	if tlsConfig != nil {
		dialer = &websocket.Dialer{TLSClientConfig: tlsConfig}
	}

	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to dial Docker websocket: %w", err)
	}

	socket := wsrouter.NewProtoRouter(true)
	router := wsrouter.NewACPRouter()
	acpConn := wsrouter.NewACPConn(socket, router.ServeMessage)
	router.SetConn(acpConn)
	// socket.Handle(1, acpConn.Handle)
	socket.Handle(1, func(msg wsrouter.ProtoMessage) {
		logger.Info().Msg("reading message")

		var acpMsg protobyss.ACPContainer

		err := proto.Unmarshal(msg.Content, &acpMsg)
		if err != nil {
			logger.Error().Err(err).Msg("invalid proto")
		}

		newMsgs, err := plugins.HandleMessage(context.Background(), &acpMsg)
		if err != nil {
			go acpConn.Handle(msg)
			logger.Error().Err(err).Msg("failed to handle message")
			return
		}

		for _, v := range newMsgs {
			marshalled, err := proto.Marshal(v)
			if err != nil {
				logger.Error().Err(err).Msg("failed to marshal message")
			}
			go acpConn.Handle(wsrouter.ProtoMessage{
				TypeID:  1,
				Content: marshalled,
			})
		}
	})

	socket.WriteHandler(1, func(msg wsrouter.ProtoMessage) []wsrouter.ProtoMessage {
		var acpMsg protobyss.ACPContainer

		err := proto.Unmarshal(msg.Content, &acpMsg)
		if err != nil {
			logger.Error().Err(err).Msg("invalid proto")
		}

		newMsgs, err := plugins.HandleMessage(context.Background(), &acpMsg)
		if err != nil {
			logger.Error().Err(err).Msg("failed to handle message")
			return []wsrouter.ProtoMessage{msg}
		}

		msgs, err := fp.MapE(func(msg *protobyss.ACPContainer) (wsrouter.ProtoMessage, error) {
			marshalled, err := proto.Marshal(msg)
			if err != nil {
				return wsrouter.ProtoMessage{}, err
			}

			return wsrouter.ProtoMessage{
				TypeID:  1,
				Content: marshalled,
			}, nil
		}, newMsgs)
		if err != nil {
			return []wsrouter.ProtoMessage{msg}
		}

		return msgs
	})

	proxiedAgent := NewProxiedACPAgent(router)

	go func() {
		//nolint:staticcheck
		socket.Serve(conn)
	}()

	return conn, socket, proxiedAgent, nil
}

// closeConn closes the websocket connection, logging any error.
func closeConn(conn *websocket.Conn, logger zerolog.Logger) {
	if closeErr := conn.Close(); closeErr != nil {
		logger.Warn().Err(closeErr).Msg("failed to close client websocket connection")
	}
}

// Oneshot runs a single prompt via a disposable agent container.
func Oneshot(ctx context.Context, prompt string, wsURL string, tlsConfig *tls.Config, logger zerolog.Logger) error {
	plugMgr, err := plugin.NewACPManager(ctx)
	if err != nil {
		return err
	}
	conn, _, proxiedAgent, err := dialAndServe(ctx, wsURL, tlsConfig, plugMgr, logger)
	if err != nil {
		return err
	}
	defer closeConn(conn, logger)

	termACPClient := acptools.NewTermACPClient()
	proxiedAgent.router.SetClient(termACPClient)

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

	// TODO close this cleanly so it doesn't spray errors

	return nil
}

// RunClient dials the websocket server at wsURL and bridges it to a client
// (typically an editor) over stdio. A non-nil tlsConfig enables TLS for the
// connection.
func RunClient(
	ctx context.Context,
	agentConfig *agentconfig.AgentConfig,
	wsURL string,
	tlsConfig *tls.Config,
	logger zerolog.Logger) error {
	plugMgr, err := plugin.NewACPManager(ctx)
	if err != nil {
		return err
	}

	for _, v := range agentConfig.Plugins.Client {
		err = plugMgr.Load(ctx, v.Path)
		if err != nil {
			return fmt.Errorf("failed to load plugin %q: %w", v.Path, err)
		}
	}

	conn, _, proxiedAgent, err := dialAndServe(ctx, wsURL, tlsConfig, plugMgr, logger)
	if err != nil {
		return err
	}
	defer closeConn(conn, logger)

	agent := NewHostProxy(proxiedAgent, proxiedAgent.router, logger)
	asc := acp.NewAgentSideConnection(agent, os.Stdout, os.Stdin)
	asc.SetLogger(slog.Default())
	agent.SetAgentConnection(asc)

	<-asc.Done()
	logger.Info().Msg("agent closed connection")
	return nil
}
