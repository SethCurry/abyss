package abyss

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/coder/acp-go-sdk"
)

type MessageDirection string

const (
	ToACPClient MessageDirection = "to_acp_client"
	ToAgent     MessageDirection = "to_agent"
)

func NewMessage[T any](msgTypeID MessageType, direction MessageDirection) *Message[T] {
	return &Message[T]{
		typeID: msgTypeID,
	}
}

type Message[T any] struct {
	typeID    MessageType
	direction MessageDirection
}

func (m *Message[T]) Unmarshal(content []byte) (T, error) {
	var resp T
	err := json.Unmarshal(content, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// TypeID returns the identifier for this message type.
func (m *Message[T]) TypeID() MessageType { return m.typeID }

// Type returns the concrete Go type this message unmarshals into.
func (m *Message[T]) Type() reflect.Type { return reflect.TypeOf((*T)(nil)).Elem() }

// UnmarshalAny unmarshals content into the message's concrete type, returned as any.
func (m *Message[T]) UnmarshalAny(content []byte) (any, error) {
	return m.Unmarshal(content)
}

func (m *Message[T]) Direction() MessageDirection { return m.direction }

// TypedMessage is the common interface for generic Message[T] values so they can
// be stored together in a single heterogeneous registry.
type TypedMessage interface {
	TypeID() MessageType
	Type() reflect.Type
	UnmarshalAny(content []byte) (any, error)
	Direction() MessageDirection
}

func GetMessageTypeByID(msgTypeID int32) (TypedMessage, error) {
	for _, v := range AllMessages {
		if int32(v.TypeID()) == msgTypeID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no ACP message type with ID %d", msgTypeID)
}

func GetMessageTypeByType(msg any) (TypedMessage, error) {
	msgType := reflect.TypeOf(msg)
	for _, v := range AllMessages {
		if msgType == v.Type() {
			return v, nil
		}
	}

	return nil, fmt.Errorf("unknown router message type %T", msg)
}

var AllMessages = []TypedMessage{
	NewMessage[acp.RequestPermissionRequest](RequestPermissionRequestType, ToACPClient),
	NewMessage[acp.RequestPermissionResponse](RequestPermissionResponseType, ToAgent),
	NewMessage[acp.WriteTextFileRequest](WriteTextFileRequestType, ToACPClient),
	NewMessage[acp.WriteTextFileResponse](WriteTextFileResponseType, ToAgent),
	NewMessage[acp.ReadTextFileRequest](ReadTextFileRequestType, ToACPClient),
	NewMessage[acp.ReadTextFileResponse](ReadTextFileResponseType, ToAgent),
	NewMessage[acp.CreateTerminalRequest](CreateTerminalRequestType, ToACPClient),
	NewMessage[acp.CreateTerminalResponse](CreateTerminalResponseType, ToAgent),
	NewMessage[acp.TerminalOutputRequest](TerminalOutputRequestType, ToACPClient),
	NewMessage[acp.TerminalOutputResponse](TerminalOutputResponseType, ToAgent),
	NewMessage[acp.ReleaseTerminalRequest](ReleaseTerminalRequestType, ToACPClient),
	NewMessage[acp.ReleaseTerminalResponse](ReleaseTerminalResponseType, ToAgent),
	NewMessage[acp.WaitForTerminalExitRequest](WaitForTerminalExitRequestType, ToACPClient),
	NewMessage[acp.WaitForTerminalExitResponse](WaitForTerminalExitResponseType, ToAgent),
	NewMessage[acp.KillTerminalRequest](KillTerminalRequestType, ToACPClient),
	NewMessage[acp.KillTerminalResponse](KillTerminalResponseType, ToAgent),
	NewMessage[acp.SessionNotification](SessionNotificationType, ToACPClient),
	NewMessage[acp.SetSessionModeRequest](SetSessionModeRequestType, ToAgent),
	NewMessage[acp.SetSessionModeResponse](SetSessionModeResponseType, ToACPClient),
	NewMessage[acp.UnstableForkSessionRequest](UnstableForkSessionRequestType, ToAgent),
	NewMessage[acp.UnstableForkSessionResponse](UnstableForkSessionResponseType, ToACPClient),
	NewMessage[acp.ListSessionsRequest](ListSessionsRequestType, ToAgent),
	NewMessage[acp.ListSessionsResponse](ListSessionsResponseType, ToACPClient),
	NewMessage[acp.ResumeSessionRequest](ResumeSessionRequestType, ToAgent),
	NewMessage[acp.ResumeSessionResponse](ResumeSessionResponseType, ToACPClient),
	NewMessage[acp.SetSessionConfigOptionRequest](SetSessionConfigOptionRequestType, ToAgent),
	NewMessage[acp.SetSessionConfigOptionResponse](SetSessionConfigOptionResponseType, ToACPClient),
	NewMessage[acp.LogoutRequest](LogoutRequestType, ToAgent),
	NewMessage[acp.LogoutResponse](LogoutResponseType, ToACPClient),
	NewMessage[acp.UnstableCloseNesRequest](UnstableCloseNesRequestType, ToAgent),
	NewMessage[acp.UnstableCloseNesResponse](UnstableCloseNesResponseType, ToACPClient),
	NewMessage[acp.UnstableStartNesRequest](UnstableStartNesRequestType, ToAgent),
	NewMessage[acp.UnstableStartNesResponse](UnstableStartNesResponseType, ToACPClient),
	NewMessage[acp.UnstableSuggestNesRequest](UnstableSuggestNesRequestType, ToAgent),
	NewMessage[acp.UnstableSuggestNesResponse](UnstableSuggestNesResponseType, ToACPClient),
	NewMessage[acp.UnstableAcceptNesNotification](UnstableAcceptNesNotificationType, ToAgent),
	NewMessage[acp.UnstableRejectNesNotification](UnstableRejectNesNotificationType, ToAgent),
	NewMessage[acp.UnstableDidChangeDocumentNotification](UnstableDidChangeDocumentNotificationType, ToAgent),
	NewMessage[acp.UnstableDidCloseDocumentNotification](UnstableDidCloseDocumentNotificationType, ToAgent),
	NewMessage[acp.UnstableDidFocusDocumentNotification](UnstableDidFocusDocumentNotificationType, ToAgent),
	NewMessage[acp.UnstableDidOpenDocumentNotification](UnstableDidOpenDocumentNotificationType, ToAgent),
	NewMessage[acp.UnstableDidSaveDocumentNotification](UnstableDidSaveDocumentNotificationType, ToAgent),
	NewMessage[acp.UnstableDisableProviderRequest](UnstableDisableProviderRequestType, ToAgent),
	NewMessage[acp.UnstableDisableProviderResponse](UnstableDisableProviderResponseType, ToACPClient),
	NewMessage[acp.UnstableListProvidersRequest](UnstableListProvidersRequestType, ToAgent),
	NewMessage[acp.UnstableListProvidersResponse](UnstableListProvidersResponseType, ToACPClient),
	NewMessage[acp.UnstableSetProviderRequest](UnstableSetProviderRequestType, ToAgent),
	NewMessage[acp.UnstableSetProviderResponse](UnstableSetProviderResponseType, ToACPClient),
	NewMessage[acp.UnstableDeleteSessionRequest](UnstableDeleteSessionRequestType, ToAgent),
	NewMessage[acp.UnstableDeleteSessionResponse](UnstableDeleteSessionResponseType, ToACPClient),
	NewMessage[acp.CloseSessionRequest](CloseSessionRequestType, ToAgent),
	NewMessage[acp.CloseSessionResponse](CloseSessionResponseType, ToACPClient),
	NewMessage[acp.InitializeRequest](InitializeRequestType, ToAgent),
	NewMessage[acp.InitializeResponse](InitializeResponseType, ToACPClient),
	NewMessage[acp.NewSessionRequest](NewSessionRequestType, ToAgent),
	NewMessage[acp.NewSessionResponse](NewSessionResponseType, ToACPClient),
	NewMessage[acp.AuthenticateRequest](AuthenticateRequestType, ToAgent),
	NewMessage[acp.AuthenticateResponse](AuthenticateResponseType, ToACPClient),
	NewMessage[acp.LoadSessionRequest](LoadSessionRequestType, ToAgent),
	NewMessage[acp.LoadSessionResponse](LoadSessionResponseType, ToACPClient),
	NewMessage[acp.PromptRequest](PromptRequestType, ToAgent),
	NewMessage[acp.PromptResponse](PromptResponseType, ToACPClient),
	NewMessage[acp.CancelNotification](CancelNotificationType, ToAgent),
}
