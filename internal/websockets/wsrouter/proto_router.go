package wsrouter

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// IProtoRouter is the interface for a protobuf router.  Used to allow
// for unit testing.
type IProtoRouter interface {
	WriteMessage(int, []byte) error
	Handle(int, func(ProtoMessage))
}

// NewProtoRouter creates a new *ProtoRouter.
func NewProtoRouter() *ProtoRouter {
	sock := &ProtoRouter{
		logger:   log.Logger.With().Str("from", "ProtoRouter").Logger(),
		handlers: make(map[int]func(ProtoMessage)),
	}

	return sock
}

var _ IProtoRouter = &ProtoRouter{}

// ProtoMessage encapsulates all of the information from a protobuf message transmitted
// via websocket.
type ProtoMessage struct {
	TypeID  int
	Content []byte
}

// ProtoRouter maps protobuf schema numbers to handles by their schema number.
// It's semantically similar to HTTP routing by path.
type ProtoRouter struct {
	conn     *websocket.Conn
	logger   zerolog.Logger
	handlers map[int]func(ProtoMessage)
	writeMut sync.Mutex
}

// Serve runs a loop that reads messages and synchronously dispatches them to handlers.
// It is goroutine-safe, so run this in a goroutine if you want async.
func (s *ProtoRouter) Serve(ws *websocket.Conn) {
	s.conn = ws
	for {
		mt, content, err := s.conn.ReadMessage()
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to read raw websocket message")
			continue
		}

		sendTo, ok := s.handlers[mt]
		if ok {
			sendTo(ProtoMessage{
				TypeID:  mt,
				Content: content,
			})
		} else {
			s.logger.Debug().Int("message_type_id", mt).Msg("no receiving channel for message type")
		}
	}
}

// Handle configures the provided function to handle protobuf schemas with the given number.
// It does not check if there is an existing handler.  Existing handlers are overwritten.
func (s *ProtoRouter) Handle(mt int, handler func(ProtoMessage)) {
	s.handlers[mt] = handler
}

// WriteMessage writes a message to the websocket.  This method is protected by a mutex,
// and is thread-safe.
func (s *ProtoRouter) WriteMessage(mt int, data []byte) error {
	s.writeMut.Lock()
	defer s.writeMut.Unlock()
	return s.conn.WriteMessage(mt, data)
}
