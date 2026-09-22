package abyss

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SethCurry/abyss/internal/fp"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

// MessageTypeID is an enum used to identify which ACP struct the contents
// of a protobyss.ACPContainer should be unmarshalled into.
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

func newMessageTypeT[T any](msgTypeID MessageTypeID, direction MessageDirection, isResponse bool) *MessageTypeT[T] {
	return &MessageTypeT[T]{
		typeID:     msgTypeID,
		direction:  direction,
		isResponse: isResponse,
	}
}

// MessageTypeT is a generic struct that holds a message type's ID and direction.
// It is primarily used to unmarshal ACP messages by their MessageTypeID
type MessageTypeT[T any] struct {
	typeID     MessageTypeID
	direction  MessageDirection
	isResponse bool
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

func (m *MessageTypeT[T]) IsResponse() bool {
	return m.isResponse
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
	IsResponse() bool
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
	return fp.MapE(ACPContainer[T], msgs)
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
	RequestPermissionRequestMsg = newMessageTypeT[acp.RequestPermissionRequest](
		RequestPermissionRequestType,
		ToACPClient,
		false)

	RequestPermissionResponseMsg = newMessageTypeT[acp.RequestPermissionResponse](
		RequestPermissionResponseType, ToAgent, true)

	WriteTextFileRequestMsg = newMessageTypeT[acp.WriteTextFileRequest](
		WriteTextFileRequestType, ToACPClient, false)

	WriteTextFileResponseMsg = newMessageTypeT[acp.WriteTextFileResponse](
		WriteTextFileResponseType, ToAgent, true)

	ReadTextFileRequestMsg = newMessageTypeT[acp.ReadTextFileRequest](
		ReadTextFileRequestType, ToACPClient, false)

	ReadTextFileResponseMsg = newMessageTypeT[acp.ReadTextFileResponse](
		ReadTextFileResponseType, ToAgent, true)

	CreateTerminalRequestMsg = newMessageTypeT[acp.CreateTerminalRequest](
		CreateTerminalRequestType, ToACPClient, false)

	CreateTerminalResponseMsg = newMessageTypeT[acp.CreateTerminalResponse](
		CreateTerminalResponseType, ToAgent, true)

	TerminalOutputRequestMsg = newMessageTypeT[acp.TerminalOutputRequest](
		TerminalOutputRequestType, ToACPClient, false)

	TerminalOutputResponseMsg = newMessageTypeT[acp.TerminalOutputResponse](
		TerminalOutputResponseType, ToAgent, true)

	ReleaseTerminalRequestMsg = newMessageTypeT[acp.ReleaseTerminalRequest](
		ReleaseTerminalRequestType, ToACPClient, false)

	ReleaseTerminalResponseMsg = newMessageTypeT[acp.ReleaseTerminalResponse](
		ReleaseTerminalResponseType, ToAgent, true)

	WaitForTerminalExitRequestMsg = newMessageTypeT[acp.WaitForTerminalExitRequest](
		WaitForTerminalExitRequestType, ToACPClient, false)

	WaitForTerminalExitResponseMsg = newMessageTypeT[acp.WaitForTerminalExitResponse](
		WaitForTerminalExitResponseType, ToAgent, true)

	KillTerminalRequestMsg = newMessageTypeT[acp.KillTerminalRequest](
		KillTerminalRequestType, ToACPClient, false)

	KillTerminalResponseMsg = newMessageTypeT[acp.KillTerminalResponse](
		KillTerminalResponseType, ToAgent, true)

	SessionNotificationMsg = newMessageTypeT[acp.SessionNotification](
		SessionNotificationType, ToACPClient, false)

	SetSessionModeRequestMsg = newMessageTypeT[acp.SetSessionModeRequest](
		SetSessionModeRequestType, ToAgent, false)

	SetSessionModeResponseMsg = newMessageTypeT[acp.SetSessionModeResponse](
		SetSessionModeResponseType, ToACPClient, true)

	UnstableForkSessionRequestMsg = newMessageTypeT[acp.UnstableForkSessionRequest](
		UnstableForkSessionRequestType, ToAgent, false)

	UnstableForkSessionResponseMsg = newMessageTypeT[acp.UnstableForkSessionResponse](
		UnstableForkSessionResponseType, ToACPClient, true)

	ListSessionsRequestMsg = newMessageTypeT[acp.ListSessionsRequest](ListSessionsRequestType, ToAgent, false)

	ListSessionsResponseMsg = newMessageTypeT[acp.ListSessionsResponse](
		ListSessionsResponseType, ToACPClient, true)

	ResumeSessionRequestMsg = newMessageTypeT[acp.ResumeSessionRequest](
		ResumeSessionRequestType, ToAgent, false)

	ResumeSessionResponseMsg = newMessageTypeT[acp.ResumeSessionResponse](
		ResumeSessionResponseType, ToACPClient, true)

	SetSessionConfigOptionRequestMsg = newMessageTypeT[acp.SetSessionConfigOptionRequest](
		SetSessionConfigOptionRequestType, ToAgent, false)

	SetSessionConfigOptionResponseMsg = newMessageTypeT[acp.SetSessionConfigOptionResponse](
		SetSessionConfigOptionResponseType, ToACPClient, true)

	LogoutRequestMsg = newMessageTypeT[acp.LogoutRequest](
		LogoutRequestType, ToAgent, false)

	LogoutResponseMsg = newMessageTypeT[acp.LogoutResponse](LogoutResponseType, ToACPClient, true)

	UnstableCloseNesRequestMsg = newMessageTypeT[acp.UnstableCloseNesRequest](
		UnstableCloseNesRequestType, ToAgent, false)

	UnstableCloseNesResponseMsg = newMessageTypeT[acp.UnstableCloseNesResponse](
		UnstableCloseNesResponseType, ToACPClient, true)

	UnstableStartNesRequestMsg = newMessageTypeT[acp.UnstableStartNesRequest](
		UnstableStartNesRequestType, ToAgent, false)

	UnstableStartNesResponseMsg = newMessageTypeT[acp.UnstableStartNesResponse](
		UnstableStartNesResponseType, ToACPClient, true)

	UnstableSuggestNesRequestMsg = newMessageTypeT[acp.UnstableSuggestNesRequest](
		UnstableSuggestNesRequestType, ToAgent, false)

	UnstableSuggestNesResponseMsg = newMessageTypeT[acp.UnstableSuggestNesResponse](
		UnstableSuggestNesResponseType, ToACPClient, true)

	UnstableAcceptNesNotificationMsg = newMessageTypeT[acp.UnstableAcceptNesNotification](
		UnstableAcceptNesNotificationType, ToAgent, false)

	UnstableRejectNesNotificationMsg = newMessageTypeT[acp.UnstableRejectNesNotification](
		UnstableRejectNesNotificationType, ToAgent, false)

	UnstableDidChangeDocumentNotifMsg = newMessageTypeT[acp.UnstableDidChangeDocumentNotification](
		UnstableDidChangeDocumentNotificationType, ToAgent, false)

	UnstableDidCloseDocumentNotifMsg = newMessageTypeT[acp.UnstableDidCloseDocumentNotification](
		UnstableDidCloseDocumentNotificationType, ToAgent, false)

	UnstableDidFocusDocumentNotifMsg = newMessageTypeT[acp.UnstableDidFocusDocumentNotification](
		UnstableDidFocusDocumentNotificationType, ToAgent, false)

	UnstableDidOpenDocumentNotifMsg = newMessageTypeT[acp.UnstableDidOpenDocumentNotification](
		UnstableDidOpenDocumentNotificationType, ToAgent, false)

	UnstableDidSaveDocumentNotifMsg = newMessageTypeT[acp.UnstableDidSaveDocumentNotification](
		UnstableDidSaveDocumentNotificationType, ToAgent, false)

	UnstableDisableProviderRequestMsg = newMessageTypeT[acp.UnstableDisableProviderRequest](
		UnstableDisableProviderRequestType, ToAgent, false)

	UnstableDisableProviderResponseMsg = newMessageTypeT[acp.UnstableDisableProviderResponse](
		UnstableDisableProviderResponseType, ToACPClient, true)

	UnstableListProvidersRequestMsg = newMessageTypeT[acp.UnstableListProvidersRequest](
		UnstableListProvidersRequestType, ToAgent, false)

	UnstableListProvidersResponseMsg = newMessageTypeT[acp.UnstableListProvidersResponse](
		UnstableListProvidersResponseType, ToACPClient, true)

	UnstableSetProviderRequestMsg = newMessageTypeT[acp.UnstableSetProviderRequest](
		UnstableSetProviderRequestType, ToAgent, false)

	UnstableSetProviderResponseMsg = newMessageTypeT[acp.UnstableSetProviderResponse](
		UnstableSetProviderResponseType, ToACPClient, true)

	UnstableDeleteSessionRequestMsg = newMessageTypeT[acp.UnstableDeleteSessionRequest](
		UnstableDeleteSessionRequestType, ToAgent, false)

	UnstableDeleteSessionResponseMsg = newMessageTypeT[acp.UnstableDeleteSessionResponse](
		UnstableDeleteSessionResponseType, ToACPClient, true)

	CloseSessionRequestMsg = newMessageTypeT[acp.CloseSessionRequest](
		CloseSessionRequestType, ToAgent, false)

	CloseSessionResponseMsg = newMessageTypeT[acp.CloseSessionResponse](
		CloseSessionResponseType, ToACPClient, true)

	InitializeRequestMsg = newMessageTypeT[acp.InitializeRequest](
		InitializeRequestType, ToAgent, false)

	InitializeResponseMsg = newMessageTypeT[acp.InitializeResponse](
		InitializeResponseType, ToACPClient, true)

	NewSessionRequestMsg = newMessageTypeT[acp.NewSessionRequest](
		NewSessionRequestType, ToAgent, false)

	NewSessionResponseMsg = newMessageTypeT[acp.NewSessionResponse](
		NewSessionResponseType, ToACPClient, true)

	AuthenticateRequestMsg = newMessageTypeT[acp.AuthenticateRequest](
		AuthenticateRequestType, ToAgent, false)

	AuthenticateResponseMsg = newMessageTypeT[acp.AuthenticateResponse](
		AuthenticateResponseType, ToACPClient, true)

	LoadSessionRequestMsg = newMessageTypeT[acp.LoadSessionRequest](
		LoadSessionRequestType, ToAgent, false)

	LoadSessionResponseMsg = newMessageTypeT[acp.LoadSessionResponse](
		LoadSessionResponseType, ToACPClient, true)

	PromptRequestMsg = newMessageTypeT[acp.PromptRequest](
		PromptRequestType, ToAgent, false)

	PromptResponseMsg = newMessageTypeT[acp.PromptResponse](
		PromptResponseType, ToACPClient, true)

	CancelNotificationMsg = newMessageTypeT[acp.CancelNotification](
		CancelNotificationType, ToAgent, false)
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
