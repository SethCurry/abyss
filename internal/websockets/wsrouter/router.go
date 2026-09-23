package wsrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

// newID generates a new UUID to use for identifying a particular message.
func newID() (string, error) {
	gotUUID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return gotUUID.String(), nil
}

// NewACPConn creates an ACPConn that demuxes incoming proto messages from conn
// into handler.
func NewACPConn(conn IProtoRouter, handler func(*protobyss.ACPContainer)) *ACPConn {
	return &ACPConn{
		logger:    log.Logger.With().Str("from", "ACPConn").Logger(),
		protoConn: conn,
		handler:   handler,
	}
}

// ACPConn demuxes incoming proto messages into a handler and sends outgoing
// proto messages over an IProtoRouter.
type ACPConn struct {
	logger    zerolog.Logger
	protoConn IProtoRouter
	handler   func(*protobyss.ACPContainer)
}

// Handle unmarshals an incoming proto message and dispatches it to the handler.
func (c *ACPConn) Handle(msg ProtoMessage) {
	protoMsg := &protobyss.ACPContainer{}
	err := proto.Unmarshal(msg.Content, protoMsg)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to unmarshal proto message")
		return
	}

	c.handler(protoMsg)
}

// Send marshals and writes an outgoing proto message over the connection.
func (c *ACPConn) Send(msg *protobyss.ACPContainer) error {
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

// MessageType describes a registered ACP message, pairing its numeric ID
// with the concrete payload type and the handler that processes it.
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

// SetConn sets the connection used to send outgoing messages.
func (r *ACPRouter) SetConn(conn *ACPConn) {
	r.conn = conn
}

// SetClient sets the client that handles agent-to-client requests.
func (r *ACPRouter) SetClient(client acp.Client) {
	r.client = client
}

// SetAgent sets the agent that handles client-to-agent requests.
func (r *ACPRouter) SetAgent(agent Agent) {
	r.agent = agent
}

// Handle registers a message type and its handler with the router.
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

// Send sends a message that is not a response to a prior request.
func (r *ACPRouter) Send(message any) error {
	return r.Respond("", message)
}

// Request sends an RPC message and returns a promise for its response.
func (r *ACPRouter) Request(message any) (*Promise[*protobyss.ACPContainer], error) {
	msgID, err := newID()
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

// Respond sends a message, optionally as a response to a prior request.
func (r *ACPRouter) Respond(requestID string, message any) error {
	msgID, err := newID()
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

// ServeMessage dispatches an incoming message to the appropriate handler.
func (r *ACPRouter) ServeMessage(msg *protobyss.ACPContainer) {
	if msg.GetResponseFor() != "" {
		r.responseWatcher.Handle(r, msg)
		return
	}

	switch abyss.MessageTypeID(msg.GetTypeId()) {
	// Client capability requests (agent -> client).
	case abyss.RequestPermissionRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.RequestPermission)
	case abyss.WriteTextFileRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.WriteTextFile)
	case abyss.ReadTextFileRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.ReadTextFile)
	case abyss.CreateTerminalRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.CreateTerminal)
	case abyss.TerminalOutputRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.TerminalOutput)
	case abyss.ReleaseTerminalRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.ReleaseTerminal)
	case abyss.WaitForTerminalExitRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.WaitForTerminalExit)
	case abyss.KillTerminalRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.KillTerminal)
	case abyss.SessionNotificationType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleNotification(r, msg, r.client.SessionUpdate)

	// Agent requests (client -> agent).
	case abyss.AuthenticateRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Authenticate)
	case abyss.InitializeRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Initialize)
	case abyss.LogoutRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Logout)
	case abyss.CancelNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.Cancel)
	case abyss.CloseSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.CloseSession)
	case abyss.ListSessionsRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.ListSessions)
	case abyss.NewSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.NewSession)
	case abyss.PromptRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Prompt)
	case abyss.ResumeSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.ResumeSession)
	case abyss.SetSessionConfigOptionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.SetSessionConfigOption)
	case abyss.SetSessionModeRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.SetSessionMode)
	case abyss.LoadSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.LoadSession)

	// Experimental agent requests (client -> agent).
	case abyss.UnstableDidChangeDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidChangeDocument)
	case abyss.UnstableDidCloseDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidCloseDocument)
	case abyss.UnstableDidFocusDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidFocusDocument)
	case abyss.UnstableDidOpenDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidOpenDocument)
	case abyss.UnstableDidSaveDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidSaveDocument)
	case abyss.UnstableAcceptNesNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableAcceptNes)
	case abyss.UnstableCloseNesRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableCloseNes)
	case abyss.UnstableRejectNesNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableRejectNes)
	case abyss.UnstableStartNesRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableStartNes)
	case abyss.UnstableSuggestNesRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableSuggestNes)
	case abyss.UnstableDisableProviderRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableDisableProvider)
	case abyss.UnstableListProvidersRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableListProviders)
	case abyss.UnstableSetProviderRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableSetProvider)
	case abyss.UnstableDeleteSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableDeleteSession)
	case abyss.UnstableForkSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
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
