package wsrouter

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type IProtoRouter interface {
	WriteMessage(int, []byte) error
	Handle(int, func(ProtoMessage))
}

func NewProtoRouter() *ProtoRouter {
	sock := &ProtoRouter{
		logger:   log.Logger,
		handlers: make(map[int]func(ProtoMessage)),
	}

	return sock
}

var _ IProtoRouter = &ProtoRouter{}

type ProtoMessage struct {
	TypeID  int
	Content []byte
}

type ProtoRouter struct {
	conn     *websocket.Conn
	logger   zerolog.Logger
	handlers map[int]func(ProtoMessage)
	writeMut sync.Mutex
}

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
			go sendTo(ProtoMessage{
				TypeID:  mt,
				Content: content,
			})
		} else {
			s.logger.Debug().Int("message_type_id", mt).Msg("no receiving channel for message type")
		}
	}
}

func (s *ProtoRouter) Handle(mt int, handler func(ProtoMessage)) {
	s.handlers[mt] = handler
}

func (s *ProtoRouter) WriteMessage(mt int, data []byte) error {
	s.writeMut.Lock()
	defer s.writeMut.Unlock()
	return s.conn.WriteMessage(mt, data)
}
