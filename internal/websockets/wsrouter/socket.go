package wsrouter

import (
	"fmt"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func NewSocket(conn *websocket.Conn) *Socket {
	return &Socket{
		conn:   conn,
		logger: log.Logger,
	}
}

type Socket struct {
	conn   *websocket.Conn
	logger zerolog.Logger
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

func (s *Socket) WriteMessage(data *protobyss.Container) error {
	marshalled, err := proto.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal protobyss.Container: %w", err)
	}

	err = s.conn.WriteMessage(1, marshalled)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
