package wsrouter

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newID() (string, error) {
	gotUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return gotUUID.String(), nil
}

func NewConn(wsConn *websocket.Conn, recvChan chan *protobyss.Container) *Conn {
	return &Conn{
		socket:   NewSocket(wsConn),
		sendChan: make(chan *protobyss.Container, 100),
		recvChan: recvChan,
		logger:   log.Logger,
	}
}

type Conn struct {
	socket   *Socket
	sendChan chan *protobyss.Container
	recvChan chan *protobyss.Container
	logger   zerolog.Logger
}

func (c *Conn) writeLoop() {
	for msg := range c.sendChan {
		err := c.socket.WriteMessage(msg)
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to write protobuf websocket")
		}
	}
}

func (c *Conn) readLoop() {
	for {
		msg, err := c.socket.ReadMessage()
		if err != nil {
			c.logger.Error().Err(err).Msg("failed to read websocket message")
			continue
		}

		c.recvChan <- msg
	}
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
