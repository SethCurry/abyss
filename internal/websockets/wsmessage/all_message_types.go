package wsmessage

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/coder/acp-go-sdk"
)

func NewMessage[T any](msgTypeID MessageType) *Message[T] {
	return &Message[T]{
		TypeID: msgTypeID,
	}
}

type Message[T any] struct {
	TypeID MessageType
}

func (m *Message[T]) Unmarshal(content []byte) (T, error) {
	var resp T
	err := json.Unmarshal(content, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetTypeID returns the identifier for this message type.
func (m *Message[T]) GetTypeID() MessageType { return m.TypeID }

// ReflectType returns the concrete Go type this message unmarshals into.
func (m *Message[T]) ReflectType() reflect.Type { return reflect.TypeOf((*T)(nil)).Elem() }

// UnmarshalAny unmarshals content into the message's concrete type, returned as any.
func (m *Message[T]) UnmarshalAny(content []byte) (any, error) {
	return m.Unmarshal(content)
}

// TypedMessage is the common interface for generic Message[T] values so they can
// be stored together in a single heterogeneous registry.
type TypedMessage interface {
	GetTypeID() MessageType
	ReflectType() reflect.Type
	UnmarshalAny(content []byte) (any, error)
}

func GetMessageTypeByID(msgTypeID int32) (TypedMessage, error) {
	for _, v := range AllMessages {
		if int32(v.GetTypeID()) == msgTypeID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no ACP message type with ID %d", msgTypeID)
}

func GetMessageTypeByType(msg any) (TypedMessage, error) {
	msgType := reflect.TypeOf(msg)
	for _, v := range AllMessages {
		if msgType == v.ReflectType() {
			return v, nil
		}
	}

	return nil, fmt.Errorf("unknown router message type %T", msg)
}

var AllMessages = []TypedMessage{
	NewMessage[acp.RequestPermissionRequest](RequestPermissionRequestType),
	NewMessage[acp.RequestPermissionResponse](RequestPermissionResponseType),
	NewMessage[acp.WriteTextFileRequest](WriteTextFileRequestType),
	NewMessage[acp.WriteTextFileResponse](WriteTextFileResponseType),
	NewMessage[acp.ReadTextFileRequest](ReadTextFileRequestType),
	NewMessage[acp.ReadTextFileResponse](ReadTextFileResponseType),
	NewMessage[acp.CreateTerminalRequest](CreateTerminalRequestType),
	NewMessage[acp.CreateTerminalResponse](CreateTerminalResponseType),
	NewMessage[acp.TerminalOutputRequest](TerminalOutputRequestType),
	NewMessage[acp.TerminalOutputResponse](TerminalOutputResponseType),
	NewMessage[acp.ReleaseTerminalRequest](ReleaseTerminalRequestType),
	NewMessage[acp.ReleaseTerminalResponse](ReleaseTerminalResponseType),
	NewMessage[acp.WaitForTerminalExitRequest](WaitForTerminalExitRequestType),
	NewMessage[acp.WaitForTerminalExitResponse](WaitForTerminalExitResponseType),
	NewMessage[acp.KillTerminalRequest](KillTerminalRequestType),
	NewMessage[acp.KillTerminalResponse](KillTerminalResponseType),
	NewMessage[acp.SessionNotification](SessionNotificationType),
	NewMessage[acp.SetSessionModeRequest](SetSessionModeRequestType),
	NewMessage[acp.SetSessionModeResponse](SetSessionModeResponseType),
	NewMessage[acp.UnstableForkSessionRequest](UnstableForkSessionRequestType),
	NewMessage[acp.UnstableForkSessionResponse](UnstableForkSessionResponseType),
	NewMessage[acp.ListSessionsRequest](ListSessionsRequestType),
	NewMessage[acp.ListSessionsResponse](ListSessionsResponseType),
	NewMessage[acp.ResumeSessionRequest](ResumeSessionRequestType),
	NewMessage[acp.ResumeSessionResponse](ResumeSessionResponseType),
	NewMessage[acp.SetSessionConfigOptionRequest](SetSessionConfigOptionRequestType),
	NewMessage[acp.SetSessionConfigOptionResponse](SetSessionConfigOptionResponseType),
	NewMessage[acp.LogoutRequest](LogoutRequestType),
	NewMessage[acp.LogoutResponse](LogoutResponseType),
	NewMessage[acp.UnstableCloseNesRequest](UnstableCloseNesRequestType),
	NewMessage[acp.UnstableCloseNesResponse](UnstableCloseNesResponseType),
	NewMessage[acp.UnstableStartNesRequest](UnstableStartNesRequestType),
	NewMessage[acp.UnstableStartNesResponse](UnstableStartNesResponseType),
	NewMessage[acp.UnstableSuggestNesRequest](UnstableSuggestNesRequestType),
	NewMessage[acp.UnstableSuggestNesResponse](UnstableSuggestNesResponseType),
	NewMessage[acp.UnstableAcceptNesNotification](UnstableAcceptNesNotificationType),
	NewMessage[acp.UnstableRejectNesNotification](UnstableRejectNesNotificationType),
	NewMessage[acp.UnstableDidChangeDocumentNotification](UnstableDidChangeDocumentNotificationType),
	NewMessage[acp.UnstableDidCloseDocumentNotification](UnstableDidCloseDocumentNotificationType),
	NewMessage[acp.UnstableDidFocusDocumentNotification](UnstableDidFocusDocumentNotificationType),
	NewMessage[acp.UnstableDidOpenDocumentNotification](UnstableDidOpenDocumentNotificationType),
	NewMessage[acp.UnstableDidSaveDocumentNotification](UnstableDidSaveDocumentNotificationType),
	NewMessage[acp.UnstableDisableProviderRequest](UnstableDisableProviderRequestType),
	NewMessage[acp.UnstableDisableProviderResponse](UnstableDisableProviderResponseType),
	NewMessage[acp.UnstableListProvidersRequest](UnstableListProvidersRequestType),
	NewMessage[acp.UnstableListProvidersResponse](UnstableListProvidersResponseType),
	NewMessage[acp.UnstableSetProviderRequest](UnstableSetProviderRequestType),
	NewMessage[acp.UnstableSetProviderResponse](UnstableSetProviderResponseType),
	NewMessage[acp.UnstableDeleteSessionRequest](UnstableDeleteSessionRequestType),
	NewMessage[acp.UnstableDeleteSessionResponse](UnstableDeleteSessionResponseType),
	NewMessage[acp.CloseSessionRequest](CloseSessionRequestType),
	NewMessage[acp.CloseSessionResponse](CloseSessionResponseType),
	NewMessage[acp.InitializeRequest](InitializeRequestType),
	NewMessage[acp.InitializeResponse](InitializeResponseType),
	NewMessage[acp.NewSessionRequest](NewSessionRequestType),
	NewMessage[acp.NewSessionResponse](NewSessionResponseType),
	NewMessage[acp.AuthenticateRequest](AuthenticateRequestType),
	NewMessage[acp.AuthenticateResponse](AuthenticateResponseType),
	NewMessage[acp.LoadSessionRequest](LoadSessionRequestType),
	NewMessage[acp.LoadSessionResponse](LoadSessionResponseType),
	NewMessage[acp.PromptRequest](PromptRequestType),
	NewMessage[acp.PromptResponse](PromptResponseType),
	NewMessage[acp.CancelNotification](CancelNotificationType),
}
