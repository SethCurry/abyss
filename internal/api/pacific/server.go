package pacific

import (
	"context"
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

func Logger(ctx context.Context) zerolog.Logger {
	return getFromContext(ctx, keyLogger, log.Logger)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLogger := log.With().Str("url", r.URL.Path).Logger()
		ctx := context.WithValue(r.Context(), keyLogger, requestLogger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewServer[T any](contextCreator func(w http.ResponseWriter, r *http.Request) T) *Server[T] {
	router := chi.NewRouter()
	return &Server[T]{
		router:     router,
		getContext: contextCreator,
	}
}

type Server[T any] struct {
	router     *chi.Mux
	getContext func(http.ResponseWriter, *http.Request) T
}

func (s *Server[T]) AddRoute(method string, pattern string, handler func(T)) {
	s.router.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		madeCtx := s.getContext(w, r)
		handler(madeCtx)
	})
}

func (s *Server[T]) Serve(listenAddr string) error {
	return http.ListenAndServe(listenAddr, s.router)
}
