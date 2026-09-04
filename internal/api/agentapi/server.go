package agentapi

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/SethCurry/abyss/internal/acptools"
	"github.com/SethCurry/abyss/internal/api/pacific"
	"github.com/SethCurry/abyss/internal/websockets/wsacp"
	"github.com/SethCurry/abyss/internal/websockets/wsrouter"
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

func contextCreator(w http.ResponseWriter, r *http.Request) *RequestContext {
	return &RequestContext{
		Logger:   log.Logger.With().Str("method", r.Method).Str("path", r.URL.Path).Logger(),
		Request:  r,
		Response: w,
	}
}

func NewServer(agentCommand []string, terminalTools *acptools.TerminalTools, fileTools *acptools.FilesystemTools) *Server {
	httpSrv := pacific.NewServer(contextCreator)
	return &Server{
		httpServer: httpSrv,
		upgrader: &websocket.Upgrader{
			CheckOrigin:      func(r *http.Request) bool { return true },
			HandshakeTimeout: time.Second * 30,
		},
		agentCommand:  agentCommand,
		terminalTools: terminalTools,
		fileTools:     fileTools,
	}
}

type Server struct {
	httpServer    *pacific.Server[*RequestContext]
	upgrader      *websocket.Upgrader
	agentCommand  []string
	terminalTools *acptools.TerminalTools
	fileTools     *acptools.FilesystemTools
}

// Serve listens on addr and bridges each websocket connection to an agent
// process spawned from agentCommand.
func (s *Server) Serve(addr string) error {
	s.httpServer.AddRoute("GET", "/ws", s.handleWebsocket)
	return s.httpServer.Serve(addr)
}

// ServeTLS listens on addr and serves over TLS using the provided config,
// bridging each websocket connection to an agent process.
func (s *Server) ServeTLS(addr string, tlsConfig *tls.Config) error {
	s.httpServer.AddRoute("GET", "/ws", s.handleWebsocket)
	return s.httpServer.ServeTLS(addr, tlsConfig)
}

func (s *Server) handleWebsocket(req *RequestContext) {
	req.Logger.Error().Msg("handling a websocket conn")
	conn, err := s.upgrader.Upgrade(req.Response, req.Request, nil)
	if err != nil {
		req.Logger.Error().Err(err).Msg("failed to upgrade to websocket")
		return
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			req.Logger.Warn().Err(closeErr).Msg("failed to close websocket connection")
		}
	}()

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

	socket := wsrouter.NewProtoRouter()
	router := wsrouter.NewACPRouter()
	acpConn := wsrouter.NewACPConn(socket, router.ServeMessage)
	socket.Handle(1, acpConn.Handle)
	router.SetConn(acpConn)
	underlying := wsacp.NewProxiedACPClient(router)

	client := wsacp.NewWebsocketAgentClient(underlying, router, s.terminalTools, s.fileTools, req.Logger)
	csc := acp.NewClientSideConnection(client, stdin, stdout)
	csc.SetLogger(slog.Default())
	client.SetClientConnection(csc)

	go func() {
		socket.Serve(conn)
	}()

	<-csc.Done()
}
