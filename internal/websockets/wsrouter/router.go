package wsrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/plugin"
	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// NewID generates a new UUID to use for identifying a particular message.
func NewID() (string, error) {
	gotUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return gotUUID.String(), nil
}

// NewACPConn creates an ACPConn that demuxes incoming proto messages from conn
// into handler.
func NewACPConn(conn IProtoRouter, pluginMgr *plugin.ACPManager, location abyss.ProxyLocation, handler func(*protobyss.ACPContainer)) *ACPConn {
	return &ACPConn{
		logger:    log.Logger.With().Str("from", "ACPConn").Logger(),
		protoConn: conn,
		handler:   handler,
		location:  location,
		plugins:   pluginMgr,
	}
}

// ACPConn demuxes incoming proto messages into a handler and sends outgoing
// proto messages over an IProtoRouter.
type ACPConn struct {
	logger    zerolog.Logger
	protoConn IProtoRouter
	handler   func(*protobyss.ACPContainer)
	plugins   *plugin.ACPManager
	location  abyss.ProxyLocation
}

// Handle unmarshals an incoming proto message and dispatches it to the handler.
func (c *ACPConn) Handle(msg ProtoMessage) {
	protoMsg := &protobyss.ACPContainer{}
	err := proto.Unmarshal(msg.Content, protoMsg)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to unmarshal proto message")
		return
	}

	newMsgs, err := c.plugins.HandleMessage(context.Background(), protoMsg)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to execute plugins in ACPRouter")
	}

	for _, v := range newMsgs {
		msgType, err := abyss.GetMessageTypeByID(v.TypeId)
		if err != nil {
			c.logger.Error().
				Err(err).
				Int32("acp_message_type_id", v.TypeId).
				Msg("failed to get message type in ACPConn.Handle")
		}

		isRemote := false

		if (c.location == abyss.LocationHost && msgType.Direction() == abyss.ToAgent) ||
			(c.location == abyss.LocationContainer && msgType.Direction() == abyss.ToACPClient) {
			isRemote = true
		}

		if isRemote {
			err = c.Send(v)
			if err != nil {
				c.logger.Error().
					Err(err).
					Msg("failed to send message")
			}
		} else {
			c.handler(v)
		}
	}
}

// Send marshals and writes an outgoing proto message over the connection.
func (c *ACPConn) Send(msg *protobyss.ACPContainer) error {
	newMsgs, err := c.plugins.HandleMessage(context.Background(), msg)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to execute plugins in ACPRouter")
	}

	for _, v := range newMsgs {
		msgType, err := abyss.GetMessageTypeByID(v.TypeId)
		if err != nil {
			c.logger.Error().
				Err(err).
				Int32("acp_message_type_id", v.TypeId).
				Msg("failed to get message type in ACPConn.Handle")
		}

		isRemote := false

		if (c.location == abyss.LocationHost && msgType.Direction() == abyss.ToAgent) ||
			(c.location == abyss.LocationContainer && msgType.Direction() == abyss.ToACPClient) {
			isRemote = true
		}

		if v.GetMessageId() == "" {
			newID, err := NewID()
			if err != nil {
				return err
			}
			v.MessageId = newID
		}

		if isRemote {
			marshalled, err := proto.Marshal(v)
			if err != nil {
				c.logger.Error().
					Err(err).
					Msg("failed to marshal proto message in ACPConn")
			}

			err = c.Send(v)
			if err != nil {
				c.logger.Error().
					Err(err).
					Msg("failed to send message")
			}

			err = c.protoConn.WriteMessage(1, marshalled)
			if err != nil {
				c.logger.Error().Err(err).Msg("")
			}
		} else {
			c.handler(v)
		}
	}

	return nil
}

// MessageType describes a registered ACP message,
// pairing its numeric ID with the concrete payload type
// and the handler that processes it.
type MessageType struct {
	ID      int32
	Type    reflect.Type
	Handler func(*ACPRouter, *protobyss.ACPContainer) any
	IsRPC   bool
}

// NewACPRouter creates a new ACPRouter with a logger and response watcher.
func NewACPRouter() *ACPRouter {
	return &ACPRouter{
		logger:          log.Logger.With().Str("from", "ACPRouter").Logger(),
		responseWatcher: NewResponseWatcher(),
	}
}

// Agent is the set of ACP agent interfaces the router can dispatch to.
type Agent interface {
	acp.Agent
	acp.AgentLoader
	acp.AgentExperimental
}

// ACPRouter is a router for the ACP websockets protocol.
type ACPRouter struct {
	messageTypes    []MessageType
	conn            *ACPConn
	logger          zerolog.Logger
	responseWatcher *ResponseWatcher
	client          acp.Client
	agent           Agent
}

// SetConn sets the ACP connection for the router.
func (r *ACPRouter) SetConn(conn *ACPConn) {
	r.conn = conn
}

// SetClient sets the ACP client for the router.
func (r *ACPRouter) SetClient(client acp.Client) {
	r.client = client
}

// SetAgent sets the ACP agent for the router.
func (r *ACPRouter) SetAgent(agent Agent) {
	r.agent = agent
}

// Handle registers a function to be a handler for a single
// type of ACP message.
func (r *ACPRouter) Handle(
	id int32,
	messageType any,
	handler func(*ACPRouter, *protobyss.ACPContainer) any,
	isRPC bool) {
	r.messageTypes = append(r.messageTypes, MessageType{
		ID:      id,
		Type:    reflect.TypeOf(messageType),
		Handler: handler,
		IsRPC:   isRPC,
	})
}

// Send sends a single ACP message, but does not set
// the ResponseFor field.
func (r *ACPRouter) Send(message any) error {
	return r.Respond("", message)
}

// Request starts an ACP RPC interaction and returns
// a promise that resolves to the response for that request.
func (r *ACPRouter) Request(message any) (*Promise[*protobyss.ACPContainer], error) {
	msgID, err := NewID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := abyss.GetMessageTypeByType(message)
	if err != nil {
		return nil, err
	}

	prom := r.responseWatcher.Register(msgID)

	jsonMarshalled, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.conn.Send(&protobyss.ACPContainer{
		MessageId: msgID,
		TypeId:    int32(msgType.TypeID()),
		Content:   jsonMarshalled})
	if err != nil {
		return nil, fmt.Errorf("failed to send RPC request: %w", err)
	}

	return prom, nil
}

// Respond generates an RPC response, including setting the
// ResponseFor field
func (r *ACPRouter) Respond(requestID string, message any) error {
	msgID, err := NewID()
	if err != nil {
		return fmt.Errorf("failed to generate message UUID: %w", err)
	}

	msgType, err := abyss.GetMessageTypeByType(message)
	if err != nil {
		return err
	}

	marshalled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.conn.Send(&protobyss.ACPContainer{
		MessageId:   msgID,
		ResponseFor: requestID,
		TypeId:      int32(msgType.TypeID()),
		Content:     marshalled,
	})
}

// ServeMessage dispatches an incoming ACP message to the appropriate handler.
func (r *ACPRouter) ServeMessage(msg *protobyss.ACPContainer) {
	if r.client == nil {
		r.logger.Warn().Msg("no client configured")
	}

	switch abyss.MessageTypeID(msg.GetTypeId()) {
	// Client capability requests (agent -> client).
	case abyss.RequestPermissionRequestType:
		handleRequest(r, msg, r.client.RequestPermission)
	case abyss.WriteTextFileRequestType:
		handleRequest(r, msg, r.client.WriteTextFile)
	case abyss.ReadTextFileRequestType:
		handleRequest(r, msg, r.client.ReadTextFile)
	case abyss.CreateTerminalRequestType:
		handleRequest(r, msg, r.client.CreateTerminal)
	case abyss.TerminalOutputRequestType:
		handleRequest(r, msg, r.client.TerminalOutput)
	case abyss.ReleaseTerminalRequestType:
		handleRequest(r, msg, r.client.ReleaseTerminal)
	case abyss.WaitForTerminalExitRequestType:
		handleRequest(r, msg, r.client.WaitForTerminalExit)
	case abyss.KillTerminalRequestType:
		handleRequest(r, msg, r.client.KillTerminal)
	case abyss.SessionNotificationType:
		handleNotification(r, msg, r.client.SessionUpdate)

	// Agent requests (client -> agent).
	case abyss.AuthenticateRequestType:
		handleRequest(r, msg, r.agent.Authenticate)
	case abyss.InitializeRequestType:
		handleRequest(r, msg, r.agent.Initialize)
	case abyss.LogoutRequestType:
		handleRequest(r, msg, r.agent.Logout)
	case abyss.CancelNotificationType:
		handleNotification(r, msg, r.agent.Cancel)
	case abyss.CloseSessionRequestType:
		handleRequest(r, msg, r.agent.CloseSession)
	case abyss.ListSessionsRequestType:
		handleRequest(r, msg, r.agent.ListSessions)
	case abyss.NewSessionRequestType:
		handleRequest(r, msg, r.agent.NewSession)
	case abyss.PromptRequestType:
		handleRequest(r, msg, r.agent.Prompt)
	case abyss.ResumeSessionRequestType:
		handleRequest(r, msg, r.agent.ResumeSession)
	case abyss.SetSessionConfigOptionRequestType:
		handleRequest(r, msg, r.agent.SetSessionConfigOption)
	case abyss.SetSessionModeRequestType:
		handleRequest(r, msg, r.agent.SetSessionMode)
	case abyss.LoadSessionRequestType:
		handleRequest(r, msg, r.agent.LoadSession)

	// Experimental agent requests (client -> agent).
	case abyss.UnstableDidChangeDocumentNotificationType:
		handleNotification(r, msg, r.agent.UnstableDidChangeDocument)
	case abyss.UnstableDidCloseDocumentNotificationType:
		handleNotification(r, msg, r.agent.UnstableDidCloseDocument)
	case abyss.UnstableDidFocusDocumentNotificationType:
		handleNotification(r, msg, r.agent.UnstableDidFocusDocument)
	case abyss.UnstableDidOpenDocumentNotificationType:
		handleNotification(r, msg, r.agent.UnstableDidOpenDocument)
	case abyss.UnstableDidSaveDocumentNotificationType:
		handleNotification(r, msg, r.agent.UnstableDidSaveDocument)
	case abyss.UnstableAcceptNesNotificationType:
		handleNotification(r, msg, r.agent.UnstableAcceptNes)
	case abyss.UnstableCloseNesRequestType:
		handleRequest(r, msg, r.agent.UnstableCloseNes)
	case abyss.UnstableRejectNesNotificationType:
		handleNotification(r, msg, r.agent.UnstableRejectNes)
	case abyss.UnstableStartNesRequestType:
		handleRequest(r, msg, r.agent.UnstableStartNes)
	case abyss.UnstableSuggestNesRequestType:
		handleRequest(r, msg, r.agent.UnstableSuggestNes)
	case abyss.UnstableDisableProviderRequestType:
		handleRequest(r, msg, r.agent.UnstableDisableProvider)
	case abyss.UnstableListProvidersRequestType:
		handleRequest(r, msg, r.agent.UnstableListProviders)
	case abyss.UnstableSetProviderRequestType:
		handleRequest(r, msg, r.agent.UnstableSetProvider)
	case abyss.UnstableDeleteSessionRequestType:
		handleRequest(r, msg, r.agent.UnstableDeleteSession)
	case abyss.UnstableForkSessionRequestType:
		handleRequest(r, msg, r.agent.UnstableForkSession)
	default:
		r.logger.Warn().Int32("type_id", msg.GetTypeId()).Msg("unhandled message type")
	}
}

func handleRequest[T, R any](r *ACPRouter, msg *protobyss.ACPContainer, fn func(context.Context, T) (R, error)) {
	var params T
	if err := json.Unmarshal(msg.GetContent(), &params); err != nil {
		r.logger.Warn().Int32("type_id", msg.GetTypeId()).Err(err).Msg("failed to unmarshal message")
		return
	}

	resp, err := fn(context.Background(), params)
	if err != nil {
		r.logger.Warn().Int32("type_id", msg.GetTypeId()).Err(err).Msg("failed to handle request")
		return
	}

	if err := r.Respond(msg.GetMessageId(), resp); err != nil {
		r.logger.Warn().Err(err).Msg("failed to send response")
	}
}

func handleNotification[T any](r *ACPRouter, msg *protobyss.ACPContainer, fn func(context.Context, T) error) {
	var params T
	if err := json.Unmarshal(msg.GetContent(), &params); err != nil {
		r.logger.Warn().Int32("type_id", msg.GetTypeId()).Err(err).Msg("failed to unmarshal message")
		return
	}

	if err := fn(context.Background(), params); err != nil {
		r.logger.Warn().Int32("type_id", msg.GetTypeId()).Err(err).Msg("failed to handle notification")
	}
}
