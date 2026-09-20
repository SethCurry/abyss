package abyss

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/fp"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

type MessageTypeID int32

const (
	MessageTypeNotExist MessageTypeID = iota

	RequestPermissionRequestType
	RequestPermissionResponseType
	WriteTextFileRequestType
	WriteTextFileResponseType
	ReadTextFileRequestType
	ReadTextFileResponseType
	CreateTerminalRequestType
	CreateTerminalResponseType
	TerminalOutputRequestType
	TerminalOutputResponseType
	ReleaseTerminalRequestType
	ReleaseTerminalResponseType
	WaitForTerminalExitRequestType
	WaitForTerminalExitResponseType
	KillTerminalRequestType
	KillTerminalResponseType
	SessionNotificationType

	SetSessionModeRequestType
	SetSessionModeResponseType
	UnstableForkSessionRequestType
	UnstableForkSessionResponseType
	ListSessionsRequestType
	ListSessionsResponseType
	ResumeSessionRequestType
	ResumeSessionResponseType
	SetSessionConfigOptionRequestType
	SetSessionConfigOptionResponseType
	LogoutRequestType
	LogoutResponseType
	UnstableCloseNesRequestType
	UnstableCloseNesResponseType
	UnstableStartNesRequestType
	UnstableStartNesResponseType
	UnstableSuggestNesRequestType
	UnstableSuggestNesResponseType
	UnstableAcceptNesNotificationType
	UnstableRejectNesNotificationType
	UnstableDidChangeDocumentNotificationType
	UnstableDidCloseDocumentNotificationType
	UnstableDidFocusDocumentNotificationType
	UnstableDidOpenDocumentNotificationType
	UnstableDidSaveDocumentNotificationType
	UnstableDisableProviderRequestType
	UnstableDisableProviderResponseType
	UnstableListProvidersRequestType
	UnstableListProvidersResponseType
	UnstableSetProviderRequestType
	UnstableSetProviderResponseType
	UnstableDeleteSessionRequestType
	UnstableDeleteSessionResponseType
	CloseSessionRequestType
	CloseSessionResponseType
	InitializeRequestType
	InitializeResponseType
	NewSessionRequestType
	NewSessionResponseType
	AuthenticateRequestType
	AuthenticateResponseType
	LoadSessionRequestType
	LoadSessionResponseType
	PromptRequestType
	PromptResponseType
	CancelNotificationType
)

// MessageDirection represents the direction of a message,
// either "to_acp_client" or "to_agent".
type MessageDirection string

const (
	// ToACPClient is the direction of a message sent from the agent to the ACP client.
	ToACPClient MessageDirection = "to_acp_client"

	// ToAgent is the direction of a message sent from the ACP client to the agent.
	ToAgent MessageDirection = "to_agent"
)

func newMessageTypeT[T any](msgTypeID MessageTypeID, direction MessageDirection) *MessageTypeT[T] {
	return &MessageTypeT[T]{
		typeID:    msgTypeID,
		direction: direction,
	}
}

// MessageTypeT is a generic struct that holds a message type's ID and direction.
// It is primarily used to unmarshal ACP messages by their MessageTypeID
type MessageTypeT[T any] struct {
	typeID    MessageTypeID
	direction MessageDirection
}

// Unmarshal unmarshals the content into the message's concrete type.
func (m *MessageTypeT[T]) Unmarshal(content []byte) (T, error) {
	var resp T
	err := json.Unmarshal(content, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// TypeID returns the identifier for this message type.
func (m *MessageTypeT[T]) TypeID() MessageTypeID { return m.typeID }

// Type returns the concrete Go type this message unmarshals into.
func (m *MessageTypeT[T]) Type() reflect.Type { return reflect.TypeOf((*T)(nil)).Elem() }

// UnmarshalAny unmarshals content into the message's concrete type, returned as any.
func (m *MessageTypeT[T]) UnmarshalAny(content []byte) (any, error) {
	return m.Unmarshal(content)
}

// Direction returns the direction of this message type.
func (m *MessageTypeT[T]) Direction() MessageDirection { return m.direction }

// TypedMessage is the common interface for generic Message[T] values so they can
// be stored together in a single heterogeneous registry.
type TypedMessage interface {
	TypeID() MessageTypeID
	Type() reflect.Type
	UnmarshalAny(content []byte) (any, error)
	Direction() MessageDirection
}

// GetMessageTypeByID returns the TypedMessage for the given message type ID,
// or an error if not found.
func GetMessageTypeByID(msgTypeID int32) (TypedMessage, error) {
	for _, v := range AllMessageTypes {
		if int32(v.TypeID()) == msgTypeID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no ACP message type with ID %d", msgTypeID)
}

// GetMessageTypeByType returns the TypedMessage for the given message,
// or an error if not found (or if the message type is not an ACP message).
func GetMessageTypeByType(msg any) (TypedMessage, error) {
	msgType := reflect.TypeOf(msg)
	for _, v := range AllMessageTypes {
		if msgType == v.Type() {
			return v, nil
		}
	}

	return nil, fmt.Errorf("unknown router message type %T", msg)
}

type acpMessageTypes interface {
	acp.RequestPermissionRequest | acp.RequestPermissionResponse |
		acp.WriteTextFileRequest | acp.WriteTextFileResponse |
		acp.ReadTextFileRequest | acp.ReadTextFileResponse |
		acp.CreateTerminalRequest | acp.CreateTerminalResponse |
		acp.TerminalOutputRequest | acp.TerminalOutputResponse |
		acp.ReleaseTerminalRequest | acp.ReleaseTerminalResponse |
		acp.WaitForTerminalExitRequest | acp.WaitForTerminalExitResponse |
		acp.KillTerminalRequest | acp.KillTerminalResponse |
		acp.SessionNotification |
		acp.SetSessionModeRequest | acp.SetSessionModeResponse |
		acp.UnstableForkSessionRequest | acp.UnstableForkSessionResponse |
		acp.ListSessionsRequest | acp.ListSessionsResponse |
		acp.ResumeSessionRequest | acp.ResumeSessionResponse |
		acp.SetSessionConfigOptionRequest | acp.SetSessionConfigOptionResponse |
		acp.LogoutRequest | acp.LogoutResponse |
		acp.UnstableCloseNesRequest | acp.UnstableCloseNesResponse |
		acp.UnstableStartNesRequest | acp.UnstableStartNesResponse |
		acp.UnstableSuggestNesRequest | acp.UnstableSuggestNesResponse |
		acp.UnstableAcceptNesNotification | acp.UnstableRejectNesNotification |
		acp.UnstableDidChangeDocumentNotification | acp.UnstableDidCloseDocumentNotification |
		acp.UnstableDidFocusDocumentNotification | acp.UnstableDidOpenDocumentNotification |
		acp.UnstableDidSaveDocumentNotification |
		acp.UnstableDisableProviderRequest | acp.UnstableDisableProviderResponse |
		acp.UnstableListProvidersRequest | acp.UnstableListProvidersResponse |
		acp.UnstableSetProviderRequest | acp.UnstableSetProviderResponse |
		acp.UnstableDeleteSessionRequest | acp.UnstableDeleteSessionResponse |
		acp.CloseSessionRequest | acp.CloseSessionResponse |
		acp.InitializeRequest | acp.InitializeResponse |
		acp.NewSessionRequest | acp.NewSessionResponse |
		acp.AuthenticateRequest | acp.AuthenticateResponse |
		acp.LoadSessionRequest | acp.LoadSessionResponse |
		acp.PromptRequest | acp.PromptResponse |
		acp.CancelNotification
}

func ACPContainer[T acpMessageTypes](msg T) (*protobyss.ACPContainer, error) {
	msgType, err := GetMessageTypeByType(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to find message type for %T: %w", msg, err)
	}

	typeId := msgType.TypeID()

	marshalled, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ACP message of type %T: %w", msg, err)
	}

	return &protobyss.ACPContainer{
		TypeId:  int32(typeId),
		Content: marshalled,
	}, nil
}

func ACPContainers[T acpMessageTypes](msgs ...T) ([]*protobyss.ACPContainer, error) {
	return fp.MapE(func(msg T) (*protobyss.ACPContainer, error) {
		return ACPContainer(msg)
	}, msgs)
}

func ACPContainerList[T acpMessageTypes](msgs ...T) (*protobyss.ACPContainerList, error) {
	newMsgs, err := ACPContainers(msgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to create ACP containers from messages: %w", err)
	}

	return &protobyss.ACPContainerList{
		Containers: newMsgs,
	}, nil
}

var (
	RequestPermissionRequestMsg        = newMessageTypeT[acp.RequestPermissionRequest](RequestPermissionRequestType, ToACPClient)
	RequestPermissionResponseMsg       = newMessageTypeT[acp.RequestPermissionResponse](RequestPermissionResponseType, ToAgent)
	WriteTextFileRequestMsg            = newMessageTypeT[acp.WriteTextFileRequest](WriteTextFileRequestType, ToACPClient)
	WriteTextFileResponseMsg           = newMessageTypeT[acp.WriteTextFileResponse](WriteTextFileResponseType, ToAgent)
	ReadTextFileRequestMsg             = newMessageTypeT[acp.ReadTextFileRequest](ReadTextFileRequestType, ToACPClient)
	ReadTextFileResponseMsg            = newMessageTypeT[acp.ReadTextFileResponse](ReadTextFileResponseType, ToAgent)
	CreateTerminalRequestMsg           = newMessageTypeT[acp.CreateTerminalRequest](CreateTerminalRequestType, ToACPClient)
	CreateTerminalResponseMsg          = newMessageTypeT[acp.CreateTerminalResponse](CreateTerminalResponseType, ToAgent)
	TerminalOutputRequestMsg           = newMessageTypeT[acp.TerminalOutputRequest](TerminalOutputRequestType, ToACPClient)
	TerminalOutputResponseMsg          = newMessageTypeT[acp.TerminalOutputResponse](TerminalOutputResponseType, ToAgent)
	ReleaseTerminalRequestMsg          = newMessageTypeT[acp.ReleaseTerminalRequest](ReleaseTerminalRequestType, ToACPClient)
	ReleaseTerminalResponseMsg         = newMessageTypeT[acp.ReleaseTerminalResponse](ReleaseTerminalResponseType, ToAgent)
	WaitForTerminalExitRequestMsg      = newMessageTypeT[acp.WaitForTerminalExitRequest](WaitForTerminalExitRequestType, ToACPClient)
	WaitForTerminalExitResponseMsg     = newMessageTypeT[acp.WaitForTerminalExitResponse](WaitForTerminalExitResponseType, ToAgent)
	KillTerminalRequestMsg             = newMessageTypeT[acp.KillTerminalRequest](KillTerminalRequestType, ToACPClient)
	KillTerminalResponseMsg            = newMessageTypeT[acp.KillTerminalResponse](KillTerminalResponseType, ToAgent)
	SessionNotificationMsg             = newMessageTypeT[acp.SessionNotification](SessionNotificationType, ToACPClient)
	SetSessionModeRequestMsg           = newMessageTypeT[acp.SetSessionModeRequest](SetSessionModeRequestType, ToAgent)
	SetSessionModeResponseMsg          = newMessageTypeT[acp.SetSessionModeResponse](SetSessionModeResponseType, ToACPClient)
	UnstableForkSessionRequestMsg      = newMessageTypeT[acp.UnstableForkSessionRequest](UnstableForkSessionRequestType, ToAgent)
	UnstableForkSessionResponseMsg     = newMessageTypeT[acp.UnstableForkSessionResponse](UnstableForkSessionResponseType, ToACPClient)
	ListSessionsRequestMsg             = newMessageTypeT[acp.ListSessionsRequest](ListSessionsRequestType, ToAgent)
	ListSessionsResponseMsg            = newMessageTypeT[acp.ListSessionsResponse](ListSessionsResponseType, ToACPClient)
	ResumeSessionRequestMsg            = newMessageTypeT[acp.ResumeSessionRequest](ResumeSessionRequestType, ToAgent)
	ResumeSessionResponseMsg           = newMessageTypeT[acp.ResumeSessionResponse](ResumeSessionResponseType, ToACPClient)
	SetSessionConfigOptionRequestMsg   = newMessageTypeT[acp.SetSessionConfigOptionRequest](SetSessionConfigOptionRequestType, ToAgent)
	SetSessionConfigOptionResponseMsg  = newMessageTypeT[acp.SetSessionConfigOptionResponse](SetSessionConfigOptionResponseType, ToACPClient)
	LogoutRequestMsg                   = newMessageTypeT[acp.LogoutRequest](LogoutRequestType, ToAgent)
	LogoutResponseMsg                  = newMessageTypeT[acp.LogoutResponse](LogoutResponseType, ToACPClient)
	UnstableCloseNesRequestMsg         = newMessageTypeT[acp.UnstableCloseNesRequest](UnstableCloseNesRequestType, ToAgent)
	UnstableCloseNesResponseMsg        = newMessageTypeT[acp.UnstableCloseNesResponse](UnstableCloseNesResponseType, ToACPClient)
	UnstableStartNesRequestMsg         = newMessageTypeT[acp.UnstableStartNesRequest](UnstableStartNesRequestType, ToAgent)
	UnstableStartNesResponseMsg        = newMessageTypeT[acp.UnstableStartNesResponse](UnstableStartNesResponseType, ToACPClient)
	UnstableSuggestNesRequestMsg       = newMessageTypeT[acp.UnstableSuggestNesRequest](UnstableSuggestNesRequestType, ToAgent)
	UnstableSuggestNesResponseMsg      = newMessageTypeT[acp.UnstableSuggestNesResponse](UnstableSuggestNesResponseType, ToACPClient)
	UnstableAcceptNesNotificationMsg   = newMessageTypeT[acp.UnstableAcceptNesNotification](UnstableAcceptNesNotificationType, ToAgent)
	UnstableRejectNesNotificationMsg   = newMessageTypeT[acp.UnstableRejectNesNotification](UnstableRejectNesNotificationType, ToAgent)
	UnstableDidChangeDocumentNotifMsg  = newMessageTypeT[acp.UnstableDidChangeDocumentNotification](UnstableDidChangeDocumentNotificationType, ToAgent)
	UnstableDidCloseDocumentNotifMsg   = newMessageTypeT[acp.UnstableDidCloseDocumentNotification](UnstableDidCloseDocumentNotificationType, ToAgent)
	UnstableDidFocusDocumentNotifMsg   = newMessageTypeT[acp.UnstableDidFocusDocumentNotification](UnstableDidFocusDocumentNotificationType, ToAgent)
	UnstableDidOpenDocumentNotifMsg    = newMessageTypeT[acp.UnstableDidOpenDocumentNotification](UnstableDidOpenDocumentNotificationType, ToAgent)
	UnstableDidSaveDocumentNotifMsg    = newMessageTypeT[acp.UnstableDidSaveDocumentNotification](UnstableDidSaveDocumentNotificationType, ToAgent)
	UnstableDisableProviderRequestMsg  = newMessageTypeT[acp.UnstableDisableProviderRequest](UnstableDisableProviderRequestType, ToAgent)
	UnstableDisableProviderResponseMsg = newMessageTypeT[acp.UnstableDisableProviderResponse](UnstableDisableProviderResponseType, ToACPClient)
	UnstableListProvidersRequestMsg    = newMessageTypeT[acp.UnstableListProvidersRequest](UnstableListProvidersRequestType, ToAgent)
	UnstableListProvidersResponseMsg   = newMessageTypeT[acp.UnstableListProvidersResponse](UnstableListProvidersResponseType, ToACPClient)
	UnstableSetProviderRequestMsg      = newMessageTypeT[acp.UnstableSetProviderRequest](UnstableSetProviderRequestType, ToAgent)
	UnstableSetProviderResponseMsg     = newMessageTypeT[acp.UnstableSetProviderResponse](UnstableSetProviderResponseType, ToACPClient)
	UnstableDeleteSessionRequestMsg    = newMessageTypeT[acp.UnstableDeleteSessionRequest](UnstableDeleteSessionRequestType, ToAgent)
	UnstableDeleteSessionResponseMsg   = newMessageTypeT[acp.UnstableDeleteSessionResponse](UnstableDeleteSessionResponseType, ToACPClient)
	CloseSessionRequestMsg             = newMessageTypeT[acp.CloseSessionRequest](CloseSessionRequestType, ToAgent)
	CloseSessionResponseMsg            = newMessageTypeT[acp.CloseSessionResponse](CloseSessionResponseType, ToACPClient)
	InitializeRequestMsg               = newMessageTypeT[acp.InitializeRequest](InitializeRequestType, ToAgent)
	InitializeResponseMsg              = newMessageTypeT[acp.InitializeResponse](InitializeResponseType, ToACPClient)
	NewSessionRequestMsg               = newMessageTypeT[acp.NewSessionRequest](NewSessionRequestType, ToAgent)
	NewSessionResponseMsg              = newMessageTypeT[acp.NewSessionResponse](NewSessionResponseType, ToACPClient)
	AuthenticateRequestMsg             = newMessageTypeT[acp.AuthenticateRequest](AuthenticateRequestType, ToAgent)
	AuthenticateResponseMsg            = newMessageTypeT[acp.AuthenticateResponse](AuthenticateResponseType, ToACPClient)
	LoadSessionRequestMsg              = newMessageTypeT[acp.LoadSessionRequest](LoadSessionRequestType, ToAgent)
	LoadSessionResponseMsg             = newMessageTypeT[acp.LoadSessionResponse](LoadSessionResponseType, ToACPClient)
	PromptRequestMsg                   = newMessageTypeT[acp.PromptRequest](PromptRequestType, ToAgent)
	PromptResponseMsg                  = newMessageTypeT[acp.PromptResponse](PromptResponseType, ToACPClient)
	CancelNotificationMsg              = newMessageTypeT[acp.CancelNotification](CancelNotificationType, ToAgent)
)

// AllMessageTypes is a list of all ACP message types.
// Used for looking up message types by ID or type.
var AllMessageTypes = []TypedMessage{
	RequestPermissionRequestMsg,
	RequestPermissionResponseMsg,
	WriteTextFileRequestMsg,
	WriteTextFileResponseMsg,
	ReadTextFileRequestMsg,
	ReadTextFileResponseMsg,
	CreateTerminalRequestMsg,
	CreateTerminalResponseMsg,
	TerminalOutputRequestMsg,
	TerminalOutputResponseMsg,
	ReleaseTerminalRequestMsg,
	ReleaseTerminalResponseMsg,
	WaitForTerminalExitRequestMsg,
	WaitForTerminalExitResponseMsg,
	KillTerminalRequestMsg,
	KillTerminalResponseMsg,
	SessionNotificationMsg,
	SetSessionModeRequestMsg,
	SetSessionModeResponseMsg,
	UnstableForkSessionRequestMsg,
	UnstableForkSessionResponseMsg,
	ListSessionsRequestMsg,
	ListSessionsResponseMsg,
	ResumeSessionRequestMsg,
	ResumeSessionResponseMsg,
	SetSessionConfigOptionRequestMsg,
	SetSessionConfigOptionResponseMsg,
	LogoutRequestMsg,
	LogoutResponseMsg,
	UnstableCloseNesRequestMsg,
	UnstableCloseNesResponseMsg,
	UnstableStartNesRequestMsg,
	UnstableStartNesResponseMsg,
	UnstableSuggestNesRequestMsg,
	UnstableSuggestNesResponseMsg,
	UnstableAcceptNesNotificationMsg,
	UnstableRejectNesNotificationMsg,
	UnstableDidChangeDocumentNotifMsg,
	UnstableDidCloseDocumentNotifMsg,
	UnstableDidFocusDocumentNotifMsg,
	UnstableDidOpenDocumentNotifMsg,
	UnstableDidSaveDocumentNotifMsg,
	UnstableDisableProviderRequestMsg,
	UnstableDisableProviderResponseMsg,
	UnstableListProvidersRequestMsg,
	UnstableListProvidersResponseMsg,
	UnstableSetProviderRequestMsg,
	UnstableSetProviderResponseMsg,
	UnstableDeleteSessionRequestMsg,
	UnstableDeleteSessionResponseMsg,
	CloseSessionRequestMsg,
	CloseSessionResponseMsg,
	InitializeRequestMsg,
	InitializeResponseMsg,
	NewSessionRequestMsg,
	NewSessionResponseMsg,
	AuthenticateRequestMsg,
	AuthenticateResponseMsg,
	LoadSessionRequestMsg,
	LoadSessionResponseMsg,
	PromptRequestMsg,
	PromptResponseMsg,
	CancelNotificationMsg,
}
