package wsrouter

import (
	"fmt"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func NewSocket(conn *websocket.Conn, sendChans map[int]chan SocketMessage) *Socket {
	return &Socket{
		conn:     conn,
		logger:   log.Logger,
		handlers: sendChans,
	}
}

type SocketMessage struct {
	TypeID  int
	Content []byte
}

type Socket struct {
	conn     *websocket.Conn
	logger   zerolog.Logger
	handlers map[int]chan SocketMessage
}

func (s *Socket) Run() {
	for {
		mt, content, err := s.conn.ReadMessage()
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to read raw websocket message")
		}

		sendTo, ok := s.handlers[mt]
		if ok {
			sendTo <- SocketMessage{
				TypeID:  mt,
				Content: content,
			}
		} else {
			s.logger.Debug().Int("message_type_id", mt).Msg("no receiving channel for message type")
		}
	}
}

// readMessage reads a MessageType- and correlation-ID-tagged JSON payload from
// the websocket.
func (s *Socket) ReadMessage() (*protobyss.Container, error) {
	req := &protobyss.Container{}

	_, data, err := s.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("failed to read message from websocket: %w", err)
	}

	err = proto.Unmarshal(data, req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobyss.Container: %w", err)
	}

	return req, nil
}

func (s *Socket) WriteMessage(mt int, data []byte) error {
	return s.conn.WriteMessage(mt, data)
}
