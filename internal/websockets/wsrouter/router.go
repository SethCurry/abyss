package wsrouter

import (
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type MessageType struct {
	ID      int32
	Type    reflect.Type
	Handler func(*Router, *protobyss.Container) any
	IsRPC   bool
}

func NewRouter(conn *Conn) *Router {
	return &Router{
		conn:            conn,
		logger:          log.Logger,
		responseWatcher: NewResponseWatcher(),
	}
}

type Router struct {
	messageTypes    []MessageType
	conn            *Conn
	logger          zerolog.Logger
	responseWatcher *ResponseWatcher
}

func (r *Router) RegisterMessage(id int32, messageType any, handler func(*Router, *protobyss.Container) any, isRPC bool) {
	r.messageTypes = append(r.messageTypes, MessageType{
		ID:      id,
		Type:    reflect.TypeOf(messageType),
		Handler: handler,
		IsRPC:   isRPC,
	})
}

func (r *Router) Send(message any) error {
	return r.Respond("", message)
}

func (r *Router) Request(message any) (*Promise[*protobyss.Container], error) {
	msgID, err := newID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := r.getMessageTypeByType(message)
	if err != nil {
		return nil, err
	}

	prom := r.responseWatcher.Register(msgID)

	err = r.conn.Notify(msgID, msgType.ID, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send RPC request: %w", err)
	}

	return prom, nil
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

	return MessageType{}, fmt.Errorf("unknown router message type %T", msg)
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
	for {
		msg, err := r.conn.Read()
		if err != nil {
			r.logger.Warn().Err(err).Msg("failed to read message from conn")
		}

		if msg.ResponseFor != "" {
			r.responseWatcher.Handle(r, msg)
			continue
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
