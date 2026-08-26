package agentapi

import (
	"log/slog"
	"net/http"
	"os"
	"os/exec"

	"github.com/SethCurry/abyss/internal/api/pacific"
	"github.com/SethCurry/abyss/internal/websockets/wsacp"
	"github.com/coder/acp-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RequestContext struct {
	Logger   zerolog.Logger
	Request  *http.Request
	Response http.ResponseWriter
}

func NewServer(agentCommand []string) *Server {
	return &Server{
		httpServer: pacific.NewServer(func(w http.ResponseWriter, r *http.Request) *RequestContext {
			return &RequestContext{
				Logger:   log.Logger.With().Str("method", r.Method).Str("path", r.URL.Path).Logger(),
				Request:  r,
				Response: w,
			}
		}),
		upgrader: websocket.Upgrader{
			// TODO this should probably actually check
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		agentCommand: agentCommand,
	}
}

type Server struct {
	httpServer   *pacific.Server[*RequestContext]
	upgrader     websocket.Upgrader
	agentCommand []string
}

// Serve listens on addr and bridges each websocket connection to an agent
// process spawned from agentCommand.
func (s *Server) Serve(addr string) error {
	s.httpServer.AddRoute("GET", "/ws", s.handleWebsocket)
	return s.httpServer.Serve(addr)
}

func (s *Server) handleWebsocket(req *RequestContext) {
	conn, err := s.upgrader.Upgrade(req.Response, req.Request, nil)
	if err != nil {
		req.Logger.Error().Err(err).Msg("failed to upgrade to websocket")
		return
	}
	defer conn.Close()

	if len(s.agentCommand) == 0 {
		req.Logger.Error().Msg("no agent command configured")
		return
	}

	req.Logger.Info().Strs("agent_command", s.agentCommand).Msg("spawning underlying agent")

	args := []string{}
	if len(s.agentCommand) > 1 {
		args = s.agentCommand[1:]
	}

	cmd := exec.CommandContext(req.Request.Context(), s.agentCommand[0], args...)

	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		req.Logger.Error().Err(err).Msg("failed to open agent stdin")
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		req.Logger.Error().Err(err).Msg("failed to open agent stdout")
		return
	}
	if err := cmd.Start(); err != nil {
		req.Logger.Error().Err(err).Msg("failed to start agent")
		return
	}
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	client := wsacp.NewWebsocketAgentClient(conn, req.Logger)
	csc := acp.NewClientSideConnection(client, stdin, stdout)
	csc.SetLogger(slog.Default())
	client.SetClientConnection(csc)

	go func() {
		if err := client.Serve(req.Request.Context()); err != nil {
			req.Logger.Error().Err(err).Msg("websocket bridge failed")
		}
		_ = conn.Close()
	}()

	<-csc.Done()
}
