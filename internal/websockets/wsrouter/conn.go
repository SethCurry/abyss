package wsrouter

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func newID() (string, error) {
	gotUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return gotUUID.String(), nil
}

func NewConn(socket *Socket, recvChan chan SocketMessage) *Conn {
	conn := &Conn{
		socket:   socket,
		sendChan: make(chan *protobyss.Container, 100),
		recvChan: recvChan,
		logger:   log.Logger,
	}

	go conn.writeLoop()

	return conn
}

type Conn struct {
	socket   *Socket
	sendChan chan *protobyss.Container
	recvChan chan SocketMessage
	logger   zerolog.Logger
}

func (c *Conn) writeLoop() {
	for msg := range c.sendChan {
		msgMarshalled, err := proto.Marshal(msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to marshal proto message")
		}
		err = c.socket.WriteMessage(1, msgMarshalled)
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to write protobuf websocket")
		}
	}
}

func (c *Conn) Read() (*protobyss.Container, error) {
	socketMsg := <-c.recvChan
	msg := &protobyss.Container{}

	err := proto.Unmarshal(socketMsg.Content, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return msg, nil
}

func (c *Conn) Notify(messageID string, messageTypeID int32, message any) error {
	return c.Respond(messageID, "", messageTypeID, message)
}

func (c *Conn) Respond(msgID string, requestID string, msgTypeID int32, message any) error {
	msgContent, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON for message: %w", err)
	}

	msg := &protobyss.Container{
		TypeId:      msgTypeID,
		MessageId:   msgID,
		Content:     msgContent,
		ResponseFor: requestID,
	}

	c.sendChan <- msg

	return nil
}

type MessageType struct {
	ID      int32
	Type    reflect.Type
	Handler func(*Router, *protobyss.Container) any
	IsRPC   bool
}
