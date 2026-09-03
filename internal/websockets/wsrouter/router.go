package wsrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/coder/acp-go-sdk"
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

// Agent is the set of ACP agent interfaces the router can dispatch to.
type Agent interface {
	acp.Agent
	acp.AgentLoader
	acp.AgentExperimental
}

type ACPRouter struct {
	messageTypes    []MessageType
	conn            *ACPConn
	logger          zerolog.Logger
	responseWatcher *ResponseWatcher
	client          acp.Client
	agent           Agent
}

func (r *ACPRouter) SetConn(conn *ACPConn) {
	r.conn = conn
}

func (r *ACPRouter) SetClient(client acp.Client) {
	r.client = client
}

func (r *ACPRouter) SetAgent(agent Agent) {
	r.agent = agent
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

	msgType, err := wsmessage.GetMessageTypeByType(message)
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
		TypeId:    int32(msgType.TypeID),
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

	msgType, err := wsmessage.GetMessageTypeByType(message)
	if err != nil {
		return err
	}

	marshalled, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return r.conn.Send(&protobyss.Container{
		MessageId:   msgID,
		ResponseFor: requestID,
		TypeId:      int32(msgType.TypeID),
		Content:     marshalled,
	})
}

func (r *ACPRouter) ServeMessage(msg *protobyss.Container) {
	if msg.ResponseFor != "" {
		r.responseWatcher.Handle(r, msg)
		return
	}

	switch wsmessage.MessageType(msg.TypeId) {
	// Client capability requests (agent -> client).
	case wsmessage.RequestPermissionRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.RequestPermission)
	case wsmessage.WriteTextFileRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.WriteTextFile)
	case wsmessage.ReadTextFileRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.ReadTextFile)
	case wsmessage.CreateTerminalRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.CreateTerminal)
	case wsmessage.TerminalOutputRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.TerminalOutput)
	case wsmessage.ReleaseTerminalRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.ReleaseTerminal)
	case wsmessage.WaitForTerminalExitRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.WaitForTerminalExit)
	case wsmessage.KillTerminalRequestType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleRequest(r, msg, r.client.KillTerminal)
	case wsmessage.SessionNotificationType:
		if r.client == nil {
			r.logger.Warn().Msg("no client configured")
			return
		}
		handleNotification(r, msg, r.client.SessionUpdate)

	// Agent requests (client -> agent).
	case wsmessage.AuthenticateRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Authenticate)
	case wsmessage.InitializeRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Initialize)
	case wsmessage.LogoutRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Logout)
	case wsmessage.CancelNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.Cancel)
	case wsmessage.CloseSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.CloseSession)
	case wsmessage.ListSessionsRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.ListSessions)
	case wsmessage.NewSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.NewSession)
	case wsmessage.PromptRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.Prompt)
	case wsmessage.ResumeSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.ResumeSession)
	case wsmessage.SetSessionConfigOptionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.SetSessionConfigOption)
	case wsmessage.SetSessionModeRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.SetSessionMode)
	case wsmessage.LoadSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.LoadSession)

	// Experimental agent requests (client -> agent).
	case wsmessage.UnstableDidChangeDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidChangeDocument)
	case wsmessage.UnstableDidCloseDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidCloseDocument)
	case wsmessage.UnstableDidFocusDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidFocusDocument)
	case wsmessage.UnstableDidOpenDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidOpenDocument)
	case wsmessage.UnstableDidSaveDocumentNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableDidSaveDocument)
	case wsmessage.UnstableAcceptNesNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableAcceptNes)
	case wsmessage.UnstableCloseNesRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableCloseNes)
	case wsmessage.UnstableRejectNesNotificationType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleNotification(r, msg, r.agent.UnstableRejectNes)
	case wsmessage.UnstableStartNesRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableStartNes)
	case wsmessage.UnstableSuggestNesRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableSuggestNes)
	case wsmessage.UnstableDisableProviderRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableDisableProvider)
	case wsmessage.UnstableListProvidersRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableListProviders)
	case wsmessage.UnstableSetProviderRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableSetProvider)
	case wsmessage.UnstableDeleteSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableDeleteSession)
	case wsmessage.UnstableForkSessionRequestType:
		if r.agent == nil {
			r.logger.Warn().Msg("no agent configured")
			return
		}
		handleRequest(r, msg, r.agent.UnstableForkSession)

	default:
		r.logger.Warn().Int32("type_id", msg.TypeId).Msg("unhandled message type")
	}
}

func handleRequest[T, R any](r *ACPRouter, msg *protobyss.Container, fn func(context.Context, T) (R, error)) {
	var params T
	if err := json.Unmarshal(msg.Content, &params); err != nil {
		r.logger.Warn().Int32("type_id", msg.TypeId).Err(err).Msg("failed to unmarshal message")
		return
	}

	resp, err := fn(context.Background(), params)
	if err != nil {
		r.logger.Warn().Int32("type_id", msg.TypeId).Err(err).Msg("failed to handle request")
		return
	}

	if err := r.Respond(msg.MessageId, resp); err != nil {
		r.logger.Warn().Err(err).Msg("failed to send response")
	}
}

func handleNotification[T any](r *ACPRouter, msg *protobyss.Container, fn func(context.Context, T) error) {
	var params T
	if err := json.Unmarshal(msg.Content, &params); err != nil {
		r.logger.Warn().Int32("type_id", msg.TypeId).Err(err).Msg("failed to unmarshal message")
		return
	}

	if err := fn(context.Background(), params); err != nil {
		r.logger.Warn().Int32("type_id", msg.TypeId).Err(err).Msg("failed to handle notification")
	}
}
