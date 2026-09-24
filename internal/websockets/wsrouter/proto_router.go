// Package wsrouter routes websocket protocol messages to handlers.
package wsrouter

import (
	"errors"
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
		logger:        log.Logger.With().Str("from", "ProtoRouter").Logger(),
		handlers:      make(map[int]func(ProtoMessage)),
		writeHandlers: make(map[int]func(ProtoMessage)),
	}

	return sock
}

var _ IProtoRouter = &ProtoRouter{}

// ProtoMessage encapsulates the information from a protobuf message
// transmitted via websocket.
type ProtoMessage struct {
	TypeID  int
	Content []byte
}

// ProtoRouter maps protobuf schema numbers to handlers by schema number.
// It's semantically similar to HTTP routing by path.
// It also contains read handlers that are called for incoming messages,
// and write handlers that are called for outgoing messages.
type ProtoRouter struct {
	conn          *websocket.Conn
	logger        zerolog.Logger
	handlers      map[int]func(ProtoMessage)
	writeHandlers map[int]func(ProtoMessage)
	writeMut      sync.Mutex
}

// Serve reads messages and synchronously dispatches them to handlers.
// It is goroutine-safe, so run this in a goroutine if you want async.
func (s *ProtoRouter) Serve(ws *websocket.Conn) {
	s.conn = ws
	for {

		mt, content, err := s.conn.ReadMessage()
		if err != nil {
			// Any read error leaves the connection unusable; reading again panics
			// with "repeated read on failed websocket connection".
			if !errors.Is(err, websocket.ErrCloseSent) {
				s.logger.Error().Err(err).Msg("failed to read raw websocket message")
			}
			return
		}

		sendTo, ok := s.handlers[mt]
		if ok {
			sendTo(ProtoMessage{
				TypeID:  mt,
				Content: content,
			})
		} else {
			s.logger.Debug().
				Int("message_type_id", mt).
				Msg("no receiving channel for message type")
		}
	}
}

// Handle configures a function to handle protobuf schemas with the given
// number. It does not check if there is an existing handler; existing
// handlers are overwritten.
func (s *ProtoRouter) Handle(mt int, handler func(ProtoMessage)) {
	s.handlers[mt] = handler
}

// WriteMessage writes a message to the websocket. This method is
// protected by a mutex and is thread-safe.
func (s *ProtoRouter) WriteMessage(mt int, data []byte) error {
	s.writeMut.Lock()
	defer s.writeMut.Unlock()

	sendTo, ok := s.writeHandlers[mt]
	if ok {
		sendTo(ProtoMessage{
			TypeID:  mt,
			Content: data,
		})
	}
	return s.conn.WriteMessage(mt, data)
}

// WriteHandler registers a write handler for the given message type.
func (s *ProtoRouter) WriteHandler(mt int, handler func(ProtoMessage)) {
	s.writeHandlers[mt] = handler
}
