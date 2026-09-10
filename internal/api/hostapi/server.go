package hostapi

import (
	"crypto/tls"
	"net/http"

	"github.com/SethCurry/abyss/internal/api/pacific"
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

func NewServer() *Server {
	httpSrv := pacific.NewServer(contextCreator)
	return &Server{
		httpServer: httpSrv,
	}
}

type Server struct {
	httpServer *pacific.Server[*RequestContext]
}

// Serve listens on addr and bridges each websocket connection to an agent
// process spawned from agentCommand.
func (s *Server) Serve(addr string) error {
	return s.httpServer.Serve(addr)
}

// ServeTLS listens on addr and serves over TLS using the provided config,
// bridging each websocket connection to an agent process.
func (s *Server) ServeTLS(addr string, tlsConfig *tls.Config) error {
	return s.httpServer.ServeTLS(addr, tlsConfig)
}
