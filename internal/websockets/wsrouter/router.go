package wsrouter

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

func NewACPConn(conn IProtoRouter, handler func(*protobyss.Container)) *ACPConn {
	return &ACPConn{
		logger:    log.Logger,
		protoConn: conn,
		handler:   handler,
	}
}

type ACPConn struct {
	logger    zerolog.Logger
	protoConn IProtoRouter
	handler   func(*protobyss.Container)
}

func (c *ACPConn) Handle(mt int, content []byte) {
	protoMsg := &protobyss.Container{}
	err := proto.Unmarshal(content, protoMsg)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to unmarshal proto message")
		return
	}

	c.handler(protoMsg)
}

func (c *ACPConn) Send(msg *protobyss.Container) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal proto message")
		return err
	}

	err = c.protoConn.WriteMessage(1, data)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to send proto message")
		return err
	}

	return nil
}

type MessageType struct {
	ID      int32
	Type    reflect.Type
	Handler func(*ACPRouter, *protobyss.Container) any
	IsRPC   bool
}

func NewACPRouter() *ACPRouter {
	return &ACPRouter{
		logger:          log.Logger,
		responseWatcher: NewResponseWatcher(),
	}
}

type ACPRouter struct {
	messageTypes    []MessageType
	conn            *ACPConn
	logger          zerolog.Logger
	responseWatcher *ResponseWatcher
}

func (r *ACPRouter) SetConn(conn *ACPConn) {
	r.conn = conn
}

func (r *ACPRouter) RouteMessage(mt int, content []byte) {
	protoMsg := &protobyss.Container{}
	err := proto.Unmarshal(content, protoMsg)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal proto message")
		return
	}

	msgType, err := r.getMessageTypeByID(protoMsg.TypeId)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get message type")
		return
	}

	msgType.Handler(r, protoMsg)
}

func (r *ACPRouter) Handle(id int32, messageType any, handler func(*ACPRouter, *protobyss.Container) any, isRPC bool) {
	r.messageTypes = append(r.messageTypes, MessageType{
		ID:      id,
		Type:    reflect.TypeOf(messageType),
		Handler: handler,
		IsRPC:   isRPC,
	})
}

func (r *ACPRouter) Send(message any) error {
	return r.Respond("", message)
}

func (r *ACPRouter) Request(message any) (*Promise[*protobyss.Container], error) {
	msgID, err := newID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := r.getMessageTypeByType(message)
	if err != nil {
		return nil, err
	}

	prom := r.responseWatcher.Register(msgID)

	jsonMarshalled, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.conn.Send(&protobyss.Container{
		MessageId: msgID,
		TypeId:    msgType.ID,
		Content:   jsonMarshalled})
	if err != nil {
		return nil, fmt.Errorf("failed to send RPC request: %w", err)
	}

	return prom, nil
}

func (r *ACPRouter) Respond(requestID string, message any) error {
	msgID, err := newID()
	if err != nil {
		return fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := r.getMessageTypeByType(message)
	if err != nil {
		return err
	}

	marshalled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.conn.Send(&protobyss.Container{
		MessageId: msgID,
		TypeId:    msgType.ID,
		Content:   marshalled,
	})
}

func (r *ACPRouter) getMessageTypeByType(msg any) (MessageType, error) {
	msgType := reflect.TypeOf(msg)
	for _, v := range r.messageTypes {
		if msgType == v.Type {
			return v, nil
		}
	}

	return MessageType{}, fmt.Errorf("unknown router message type %T", msg)
}

func (r *ACPRouter) getMessageTypeByID(msgTypeID int32) (MessageType, error) {
	for _, v := range r.messageTypes {
		if v.ID == msgTypeID {
			return v, nil
		}
	}

	return MessageType{}, fmt.Errorf("no message type with ID %d", msgTypeID)
}

func (r *ACPRouter) ServeMessage(msg *protobyss.Container) {
	if msg.ResponseFor != "" {
		r.responseWatcher.Handle(r, msg)
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
