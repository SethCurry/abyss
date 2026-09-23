// Package pacific provides a minimal HTTP API server built on chi.
package pacific

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey int

const (
	keyLogger contextKey = iota
)

func getFromContext[T any](ctx context.Context, key contextKey, defaultValue T) T {
	value := ctx.Value(key)
	if value == nil {
		return defaultValue
	}

	if asVal, ok := value.(T); ok {
		log.Error().Type("expected", defaultValue).Type("actual", value).Msg("getFromContext got value with incorrect type")
		return asVal
	}

	return defaultValue
}

// Logger returns the request-scoped logger from ctx, or the default logger.
func Logger(ctx context.Context) zerolog.Logger {
	return getFromContext(ctx, keyLogger, log.Logger)
}

// LoggerMiddleware attaches a request-scoped logger to the request context.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLogger := log.With().Str("url", r.URL.Path).Logger()
		ctx := context.WithValue(r.Context(), keyLogger, requestLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewServer creates a server that derives per-request context via creator.
func NewServer[T any](contextCreator func(w http.ResponseWriter, r *http.Request) T) *Server[T] {
	router := chi.NewRouter()
	return &Server[T]{
		router:     router,
		getContext: contextCreator,
	}
}

// Server is an HTTP server that builds a request-scoped context per request.
type Server[T any] struct {
	router     *chi.Mux
	getContext func(http.ResponseWriter, *http.Request) T
}

// AddRoute registers a handler for the given method and pattern.
func (s *Server[T]) AddRoute(method string, pattern string, handler func(T)) {
	s.router.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		madeCtx := s.getContext(w, r)
		handler(madeCtx)
	})
}

// Serve listens on listenAddr and serves HTTP requests.
func (s *Server[T]) Serve(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s.router)
}

// ServeTLS listens on addr and serves over TLS using the provided config.
// The config's Certificates and ClientAuth fields must be set for mutual TLS.
func (s *Server[T]) ServeTLS(listenAddr string, tlsConfig *tls.Config) error {
	server := &http.Server{
		Addr:      listenAddr,
		Handler:   s.router,
		TLSConfig: tlsConfig,
	}

	return server.ListenAndServeTLS("", "")
}
