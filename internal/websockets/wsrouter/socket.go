package wsrouter

import (
	"fmt"
	"sync"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type ISocket interface {
	ReadMessage() (*protobyss.ACPContainer, error)
	WriteMessage(int, []byte) error
}

func NewSocket(conn *websocket.Conn, sendChans map[int]chan SocketMessage) *Socket {
	sock := &Socket{
		conn:     conn,
		logger:   log.Logger.With().Str("from", "Socket").Logger(),
		handlers: sendChans,
	}

	go sock.run()

	return sock
}

var _ ISocket = &Socket{}

type SocketMessage struct {
	TypeID  int
	Content []byte
}

type Socket struct {
	conn     IWebsocket
	logger   zerolog.Logger
	handlers map[int]chan SocketMessage
	writeMut sync.Mutex
}

type IWebsocket interface {
	ReadMessage() (int, []byte, error)
	WriteMessage(int, []byte) error
}

func (s *Socket) run() {
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

func (s *Socket) ReadMessage() (*protobyss.ACPContainer, error) {
	req := &protobyss.ACPContainer{}

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
	s.writeMut.Lock()
	defer s.writeMut.Unlock()
	return s.conn.WriteMessage(mt, data)
}
