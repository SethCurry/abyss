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

func NewRouter(conn *Conn) *Router {
	return &Router{
		conn:         conn,
		messageTypes: []MessageType{},
		logger:       log.Logger,
	}
}

type Router struct {
	messageTypes       []MessageType
	conn               *Conn
	logger             zerolog.Logger
	waitingForResponse map[string]func(*protobyss.Container)
}

func (r *Router) RegisterMessage(msgType MessageType) {
	r.messageTypes = append(r.messageTypes, msgType)
}

func (r *Router) Send(message any) error {
	return r.Respond("", message)
}

func (r *Router) Request(message any, handler func(*protobyss.Container)) error {
	msgID, err := newID()
	if err != nil {
		return fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := r.getMessageTypeByType(message)
	if err != nil {
		return err
	}

	r.waitingForResponse[msgID] = handler

	err = r.conn.Notify(msgID, msgType.ID, message)
	if err != nil {
		return fmt.Errorf("failed to send RPC request: %w", err)
	}

	return nil
}

func (r *Router) Respond(requestID string, message any) error {
	msgID, err := newID()
	if err != nil {
		return fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := r.getMessageTypeByType(message)
	if err != nil {
		return err
	}

	return r.conn.Respond(msgID, requestID, msgType.ID, message)
}

func (r *Router) getMessageTypeByType(msg any) (MessageType, error) {
	msgType := reflect.TypeOf(msg)
	for _, v := range r.messageTypes {
		if msgType == v.Type {
			return v, nil
		}
	}

	return MessageType{}, fmt.Errorf("unknown message type %T", msg)
}

func (r *Router) getMessageTypeByID(msgTypeID int32) (MessageType, error) {
	for _, v := range r.messageTypes {
		if v.ID == msgTypeID {
			return v, nil
		}
	}

	return MessageType{}, fmt.Errorf("no message type with ID %d", msgTypeID)
}

func (r *Router) Serve() {
	for msg := range r.conn.recvChan {
		if msg.ResponseFor != "" {
			handler, ok := r.waitingForResponse[msg.ResponseFor]
			if !ok {
				r.logger.Error().Str("response_for", msg.ResponseFor).Msg("no handler found for response")
				continue
			}

			handler(msg)
		}

		messageType, err := r.getMessageTypeByID(msg.TypeId)
		if err != nil {
			r.logger.Warn().Int32("type_id", msg.TypeId).Err(err).Msg("failed to find message by type ID")
		}

		resp := messageType.Handler(r, msg)
		if messageType.IsRPC && resp != nil {
			err := r.Respond(msg.MessageId, resp)
			if err != nil {
				r.logger.Warn().Err(err).Msg("failed to send response")
			}
		}
	}
}
