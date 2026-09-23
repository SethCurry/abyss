// Package abyss defines the core domain types used across the Abyss project.
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
	// MessageTypeNotExist is the zero value for MessageTypeID and does not
	// correspond to a real message.
	MessageTypeNotExist MessageTypeID = iota

	// RequestPermissionRequestType identifies an ACP request-permission request.
	RequestPermissionRequestType

	// RequestPermissionResponseType identifies an ACP request-permission response.
	RequestPermissionResponseType

	// WriteTextFileRequestType identifies an ACP write-text-file request.
	WriteTextFileRequestType

	// WriteTextFileResponseType identifies an ACP write-text-file response.
	WriteTextFileResponseType

	// ReadTextFileRequestType identifies an ACP read-text-file request.
	ReadTextFileRequestType

	// ReadTextFileResponseType identifies an ACP read-text-file response.
	ReadTextFileResponseType

	// CreateTerminalRequestType identifies an ACP create-terminal request.
	CreateTerminalRequestType

	// CreateTerminalResponseType identifies an ACP create-terminal response.
	CreateTerminalResponseType

	// TerminalOutputRequestType identifies an ACP terminal-output request.
	TerminalOutputRequestType

	// TerminalOutputResponseType identifies an ACP terminal-output response.
	TerminalOutputResponseType

	// ReleaseTerminalRequestType identifies an ACP release-terminal request.
	ReleaseTerminalRequestType

	// ReleaseTerminalResponseType identifies an ACP release-terminal response.
	ReleaseTerminalResponseType

	// WaitForTerminalExitRequestType identifies a terminal-exit-wait request.
	WaitForTerminalExitRequestType

	// WaitForTerminalExitResponseType identifies an ACP wait-for-terminal-exit
	// response.
	WaitForTerminalExitResponseType

	// KillTerminalRequestType identifies an ACP kill-terminal request.
	KillTerminalRequestType

	// KillTerminalResponseType identifies an ACP kill-terminal response.
	KillTerminalResponseType

	// SessionNotificationType identifies an ACP session notification.
	SessionNotificationType

	// SetSessionModeRequestType identifies an ACP set-session-mode request.
	SetSessionModeRequestType

	// SetSessionModeResponseType identifies an ACP set-session-mode response.
	SetSessionModeResponseType

	// UnstableForkSessionRequestType identifies an ACP fork-session request.
	UnstableForkSessionRequestType

	// UnstableForkSessionResponseType identifies an ACP fork-session response.
	UnstableForkSessionResponseType

	// ListSessionsRequestType identifies an ACP list-sessions request.
	ListSessionsRequestType

	// ListSessionsResponseType identifies an ACP list-sessions response.
	ListSessionsResponseType

	// ResumeSessionRequestType identifies an ACP resume-session request.
	ResumeSessionRequestType

	// ResumeSessionResponseType identifies an ACP resume-session response.
	ResumeSessionResponseType

	// SetSessionConfigOptionRequestType identifies an ACP
	// set-session-config-option request.
	SetSessionConfigOptionRequestType

	// SetSessionConfigOptionResponseType identifies an ACP
	// set-session-config-option response.
	SetSessionConfigOptionResponseType

	// LogoutRequestType identifies an ACP logout request.
	LogoutRequestType

	// LogoutResponseType identifies an ACP logout response.
	LogoutResponseType

	// UnstableCloseNesRequestType identifies an ACP unstable close-nes request.
	UnstableCloseNesRequestType

	// UnstableCloseNesResponseType identifies an ACP unstable close-nes response.
	UnstableCloseNesResponseType

	// UnstableStartNesRequestType identifies an ACP unstable start-nes request.
	UnstableStartNesRequestType

	// UnstableStartNesResponseType identifies an ACP unstable start-nes response.
	UnstableStartNesResponseType

	// UnstableSuggestNesRequestType identifies an ACP
	// suggest-nes request.
	UnstableSuggestNesRequestType

	// UnstableSuggestNesResponseType identifies an ACP suggest-nes response.
	UnstableSuggestNesResponseType

	// UnstableAcceptNesNotificationType identifies an ACP accept-nes notification.
	UnstableAcceptNesNotificationType

	// UnstableRejectNesNotificationType identifies an ACP reject-nes notification.
	UnstableRejectNesNotificationType

	// UnstableDidChangeDocumentNotificationType identifies an ACP
	// did-change-document notification.
	UnstableDidChangeDocumentNotificationType

	// UnstableDidCloseDocumentNotificationType identifies an ACP
	// did-close-document notification.
	UnstableDidCloseDocumentNotificationType

	// UnstableDidFocusDocumentNotificationType identifies an ACP
	// did-focus-document notification.
	UnstableDidFocusDocumentNotificationType

	// UnstableDidOpenDocumentNotificationType identifies an ACP
	// did-open-document notification.
	UnstableDidOpenDocumentNotificationType

	// UnstableDidSaveDocumentNotificationType identifies an ACP
	// did-save-document notification.
	UnstableDidSaveDocumentNotificationType

	// UnstableDisableProviderRequestType identifies an ACP
	// disable-provider request.
	UnstableDisableProviderRequestType

	// UnstableDisableProviderResponseType identifies an ACP
	// disable-provider response.
	UnstableDisableProviderResponseType

	// UnstableListProvidersRequestType identifies an ACP list-providers request.
	UnstableListProvidersRequestType

	// UnstableListProvidersResponseType identifies an ACP list-providers response.
	UnstableListProvidersResponseType

	// UnstableSetProviderRequestType identifies an ACP set-provider request.
	UnstableSetProviderRequestType

	// UnstableSetProviderResponseType identifies an ACP set-provider response.
	UnstableSetProviderResponseType

	// UnstableDeleteSessionRequestType identifies an ACP delete-session request.
	UnstableDeleteSessionRequestType

	// UnstableDeleteSessionResponseType identifies an ACP delete-session response.
	UnstableDeleteSessionResponseType

	// CloseSessionRequestType identifies an ACP close-session request.
	CloseSessionRequestType

	// CloseSessionResponseType identifies an ACP close-session response.
	CloseSessionResponseType

	// InitializeRequestType identifies an ACP initialize request.
	InitializeRequestType

	// InitializeResponseType identifies an ACP initialize response.
	InitializeResponseType

	// NewSessionRequestType identifies an ACP new-session request.
	NewSessionRequestType

	// NewSessionResponseType identifies an ACP new-session response.
	NewSessionResponseType

	// AuthenticateRequestType identifies an ACP authenticate request.
	AuthenticateRequestType

	// AuthenticateResponseType identifies an ACP authenticate response.
	AuthenticateResponseType

	// LoadSessionRequestType identifies an ACP load-session request.
	LoadSessionRequestType

	// LoadSessionResponseType identifies an ACP load-session response.
	LoadSessionResponseType

	// PromptRequestType identifies an ACP prompt request.
	PromptRequestType

	// PromptResponseType identifies an ACP prompt response.
	PromptResponseType

	// CancelNotificationType identifies an ACP cancel notification.
	CancelNotificationType
)

// MessageDirection represents the direction of a message,
// either "to_acp_client" or "to_agent".
type MessageDirection string

const (
	// ToACPClient is the direction of a message sent from the agent
	// to the ACP client.
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

// MessageTypeT is a generic struct that holds a message type's ID
// and direction. It is primarily used to unmarshal ACP messages by
// their MessageTypeID.
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

// IsResponse reports whether this message type is a response message.
func (m *MessageTypeT[T]) IsResponse() bool {
	return m.isResponse
}

// TypeID returns the identifier for this message type.
func (m *MessageTypeT[T]) TypeID() MessageTypeID { return m.typeID }

// Type returns the concrete Go type this message unmarshals into.
func (m *MessageTypeT[T]) Type() reflect.Type { return reflect.TypeOf((*T)(nil)).Elem() }

// UnmarshalAny unmarshals content into the message's concrete type,
// returned as any.
func (m *MessageTypeT[T]) UnmarshalAny(content []byte) (any, error) {
	return m.Unmarshal(content)
}

// Direction returns the direction of this message type.
func (m *MessageTypeT[T]) Direction() MessageDirection { return m.direction }

// TypedMessage is the common interface for generic Message[T] values so
// they can be stored together in a single heterogeneous registry.
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

// ACPContainer wraps an ACP message of any supported type into a
// protobyss.ACPContainer, recording its type id and marshalling the
// payload as JSON.
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

// ACPContainers converts a variadic list of ACP messages into a slice of
// ACPContainers.
func ACPContainers[T acpMessageTypes](msgs ...T) ([]*protobyss.ACPContainer, error) {
	return fp.MapE(ACPContainer[T], msgs)
}

// ACPContainerList converts a variadic list of ACP messages into an
// ACPContainerList.
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
	// RequestPermissionRequestMsg is the message type for requesting permission
	// to perform an action.
	RequestPermissionRequestMsg = newMessageTypeT[acp.RequestPermissionRequest](
		RequestPermissionRequestType,
		ToACPClient,
		false,
	)

	// RequestPermissionResponseMsg is the message type for the response to a
	// request permission request.
	RequestPermissionResponseMsg = newMessageTypeT[acp.RequestPermissionResponse](
		RequestPermissionResponseType, ToAgent, true,
	)

	// WriteTextFileRequestMsg is the message type for requesting to write a text
	// file.
	WriteTextFileRequestMsg = newMessageTypeT[acp.WriteTextFileRequest](
		WriteTextFileRequestType, ToACPClient, false,
	)

	// WriteTextFileResponseMsg is the message type for the response to a write
	// text file request.
	WriteTextFileResponseMsg = newMessageTypeT[acp.WriteTextFileResponse](
		WriteTextFileResponseType, ToAgent, true,
	)

	// ReadTextFileRequestMsg is the message type for requesting to read a text
	// file.
	ReadTextFileRequestMsg = newMessageTypeT[acp.ReadTextFileRequest](
		ReadTextFileRequestType, ToACPClient, false,
	)

	// ReadTextFileResponseMsg is the message type for the response to a read
	// text file request.
	ReadTextFileResponseMsg = newMessageTypeT[acp.ReadTextFileResponse](
		ReadTextFileResponseType, ToAgent, true,
	)

	// CreateTerminalRequestMsg is the message type for requesting to create a
	// terminal.
	CreateTerminalRequestMsg = newMessageTypeT[acp.CreateTerminalRequest](
		CreateTerminalRequestType, ToACPClient, false,
	)

	// CreateTerminalResponseMsg is the message type for the response to a
	// create-terminal request.
	CreateTerminalResponseMsg = newMessageTypeT[acp.CreateTerminalResponse](
		CreateTerminalResponseType, ToAgent, true,
	)

	// TerminalOutputRequestMsg is the message type for requesting terminal
	// output.
	TerminalOutputRequestMsg = newMessageTypeT[acp.TerminalOutputRequest](
		TerminalOutputRequestType, ToACPClient, false,
	)

	// TerminalOutputResponseMsg is the message type for the response to a
	// terminal-output request.
	TerminalOutputResponseMsg = newMessageTypeT[acp.TerminalOutputResponse](
		TerminalOutputResponseType, ToAgent, true,
	)

	// ReleaseTerminalRequestMsg is the message type for requesting to release a
	// terminal.
	ReleaseTerminalRequestMsg = newMessageTypeT[acp.ReleaseTerminalRequest](
		ReleaseTerminalRequestType, ToACPClient, false,
	)

	// ReleaseTerminalResponseMsg is the message type for the response to a
	// release-terminal request.
	ReleaseTerminalResponseMsg = newMessageTypeT[acp.ReleaseTerminalResponse](
		ReleaseTerminalResponseType, ToAgent, true,
	)

	// WaitForTerminalExitRequestMsg is the message type for requesting to wait
	// for a terminal to exit.
	WaitForTerminalExitRequestMsg = newMessageTypeT[acp.WaitForTerminalExitRequest](
		WaitForTerminalExitRequestType, ToACPClient, false,
	)

	// WaitForTerminalExitResponseMsg is the message type for the response to a
	// wait-for-terminal-exit request.
	WaitForTerminalExitResponseMsg = newMessageTypeT[acp.WaitForTerminalExitResponse](
		WaitForTerminalExitResponseType, ToAgent, true,
	)

	// KillTerminalRequestMsg is the message type for requesting to kill a
	// terminal.
	KillTerminalRequestMsg = newMessageTypeT[acp.KillTerminalRequest](
		KillTerminalRequestType, ToACPClient, false,
	)

	// KillTerminalResponseMsg is the message type for the response to a
	// kill-terminal request.
	KillTerminalResponseMsg = newMessageTypeT[acp.KillTerminalResponse](
		KillTerminalResponseType, ToAgent, true,
	)

	// SessionNotificationMsg is the message type for a session notification.
	SessionNotificationMsg = newMessageTypeT[acp.SessionNotification](
		SessionNotificationType, ToACPClient, false,
	)

	// SetSessionModeRequestMsg is the message type for requesting to set the
	// session mode.
	SetSessionModeRequestMsg = newMessageTypeT[acp.SetSessionModeRequest](
		SetSessionModeRequestType, ToAgent, false,
	)

	// SetSessionModeResponseMsg is the message type for the response to a
	// set-session-mode request.
	SetSessionModeResponseMsg = newMessageTypeT[acp.SetSessionModeResponse](
		SetSessionModeResponseType, ToACPClient, true,
	)

	// UnstableForkSessionRequestMsg is the message type for requesting to fork
	// a session.
	UnstableForkSessionRequestMsg = newMessageTypeT[acp.UnstableForkSessionRequest](
		UnstableForkSessionRequestType, ToAgent, false,
	)

	// UnstableForkSessionResponseMsg is the message type for the response to a
	// fork-session request.
	UnstableForkSessionResponseMsg = newMessageTypeT[acp.UnstableForkSessionResponse](
		UnstableForkSessionResponseType, ToACPClient, true,
	)

	// ListSessionsRequestMsg is the message type for requesting to list
	// sessions.
	ListSessionsRequestMsg = newMessageTypeT[acp.ListSessionsRequest](ListSessionsRequestType, ToAgent, false)

	// ListSessionsResponseMsg is the message type for the response to a
	// list-sessions request.
	ListSessionsResponseMsg = newMessageTypeT[acp.ListSessionsResponse](
		ListSessionsResponseType, ToACPClient, true,
	)

	// ResumeSessionRequestMsg is the message type for requesting to resume a
	// session.
	ResumeSessionRequestMsg = newMessageTypeT[acp.ResumeSessionRequest](
		ResumeSessionRequestType, ToAgent, false,
	)

	// ResumeSessionResponseMsg is the message type for the response to a
	// resume-session request.
	ResumeSessionResponseMsg = newMessageTypeT[acp.ResumeSessionResponse](
		ResumeSessionResponseType, ToACPClient, true,
	)

	// SetSessionConfigOptionRequestMsg is the message type for requesting to
	// set a session config option.
	SetSessionConfigOptionRequestMsg = newMessageTypeT[acp.SetSessionConfigOptionRequest](
		SetSessionConfigOptionRequestType, ToAgent, false,
	)

	// SetSessionConfigOptionResponseMsg is the message type for the response
	// to a set-session-config-option request.
	SetSessionConfigOptionResponseMsg = newMessageTypeT[acp.SetSessionConfigOptionResponse](
		SetSessionConfigOptionResponseType, ToACPClient, true,
	)

	// LogoutRequestMsg is the message type for requesting to log out.
	LogoutRequestMsg = newMessageTypeT[acp.LogoutRequest](
		LogoutRequestType, ToAgent, false,
	)

	// LogoutResponseMsg is the message type for the response to a logout
	// request.
	LogoutResponseMsg = newMessageTypeT[acp.LogoutResponse](LogoutResponseType, ToACPClient, true)

	// UnstableCloseNesRequestMsg is the message type for requesting to close
	// an NES.
	UnstableCloseNesRequestMsg = newMessageTypeT[acp.UnstableCloseNesRequest](
		UnstableCloseNesRequestType, ToAgent, false,
	)

	// UnstableCloseNesResponseMsg is the message type for the response to a
	// close-NES request.
	UnstableCloseNesResponseMsg = newMessageTypeT[acp.UnstableCloseNesResponse](
		UnstableCloseNesResponseType, ToACPClient, true,
	)

	// UnstableStartNesRequestMsg is the message type for requesting to start
	// an NES.
	UnstableStartNesRequestMsg = newMessageTypeT[acp.UnstableStartNesRequest](
		UnstableStartNesRequestType, ToAgent, false,
	)

	// UnstableStartNesResponseMsg is the message type for the response to a
	// start-NES request.
	UnstableStartNesResponseMsg = newMessageTypeT[acp.UnstableStartNesResponse](
		UnstableStartNesResponseType, ToACPClient, true,
	)

	// UnstableSuggestNesRequestMsg is the message type for requesting NES
	// suggestions.
	UnstableSuggestNesRequestMsg = newMessageTypeT[acp.UnstableSuggestNesRequest](
		UnstableSuggestNesRequestType, ToAgent, false,
	)

	// UnstableSuggestNesResponseMsg is the message type for the response to an
	// NES-suggest request.
	UnstableSuggestNesResponseMsg = newMessageTypeT[acp.UnstableSuggestNesResponse](
		UnstableSuggestNesResponseType, ToACPClient, true,
	)

	// UnstableAcceptNesNotificationMsg is the message type for an accept-NES
	// notification.
	UnstableAcceptNesNotificationMsg = newMessageTypeT[acp.UnstableAcceptNesNotification](
		UnstableAcceptNesNotificationType, ToAgent, false,
	)

	// UnstableRejectNesNotificationMsg is the message type for a reject-NES
	// notification.
	UnstableRejectNesNotificationMsg = newMessageTypeT[acp.UnstableRejectNesNotification](
		UnstableRejectNesNotificationType, ToAgent, false,
	)

	// UnstableDidChangeDocumentNotifMsg is the message type for a
	// did-change-document notification.
	UnstableDidChangeDocumentNotifMsg = newMessageTypeT[acp.UnstableDidChangeDocumentNotification](
		UnstableDidChangeDocumentNotificationType, ToAgent, false,
	)

	// UnstableDidCloseDocumentNotifMsg is the message type for a
	// did-close-document notification.
	UnstableDidCloseDocumentNotifMsg = newMessageTypeT[acp.UnstableDidCloseDocumentNotification](
		UnstableDidCloseDocumentNotificationType, ToAgent, false,
	)

	// UnstableDidFocusDocumentNotifMsg is the message type for a
	// did-focus-document notification.
	UnstableDidFocusDocumentNotifMsg = newMessageTypeT[acp.UnstableDidFocusDocumentNotification](
		UnstableDidFocusDocumentNotificationType, ToAgent, false,
	)

	// UnstableDidOpenDocumentNotifMsg is the message type for a
	// did-open-document notification.
	UnstableDidOpenDocumentNotifMsg = newMessageTypeT[acp.UnstableDidOpenDocumentNotification](
		UnstableDidOpenDocumentNotificationType, ToAgent, false,
	)

	// UnstableDidSaveDocumentNotifMsg is the message type for a
	// did-save-document notification.
	UnstableDidSaveDocumentNotifMsg = newMessageTypeT[acp.UnstableDidSaveDocumentNotification](
		UnstableDidSaveDocumentNotificationType, ToAgent, false,
	)

	// UnstableDisableProviderRequestMsg is the message type for requesting to
	// disable a provider.
	UnstableDisableProviderRequestMsg = newMessageTypeT[acp.UnstableDisableProviderRequest](
		UnstableDisableProviderRequestType, ToAgent, false,
	)

	// UnstableDisableProviderResponseMsg is the message type for the response
	// to a disable-provider request.
	UnstableDisableProviderResponseMsg = newMessageTypeT[acp.UnstableDisableProviderResponse](
		UnstableDisableProviderResponseType, ToACPClient, true,
	)

	// UnstableListProvidersRequestMsg is the message type for requesting to
	// list providers.
	UnstableListProvidersRequestMsg = newMessageTypeT[acp.UnstableListProvidersRequest](
		UnstableListProvidersRequestType, ToAgent, false,
	)

	// UnstableListProvidersResponseMsg is the message type for the response to
	// a list-providers request.
	UnstableListProvidersResponseMsg = newMessageTypeT[acp.UnstableListProvidersResponse](
		UnstableListProvidersResponseType, ToACPClient, true,
	)

	// UnstableSetProviderRequestMsg is the message type for requesting to set
	// a provider.
	UnstableSetProviderRequestMsg = newMessageTypeT[acp.UnstableSetProviderRequest](
		UnstableSetProviderRequestType, ToAgent, false,
	)

	// UnstableSetProviderResponseMsg is the message type for the response to a
	// set-provider request.
	UnstableSetProviderResponseMsg = newMessageTypeT[acp.UnstableSetProviderResponse](
		UnstableSetProviderResponseType, ToACPClient, true,
	)

	// UnstableDeleteSessionRequestMsg is the message type for requesting to
	// delete a session.
	UnstableDeleteSessionRequestMsg = newMessageTypeT[acp.UnstableDeleteSessionRequest](
		UnstableDeleteSessionRequestType, ToAgent, false,
	)

	// UnstableDeleteSessionResponseMsg is the message type for the response to
	// a delete-session request.
	UnstableDeleteSessionResponseMsg = newMessageTypeT[acp.UnstableDeleteSessionResponse](
		UnstableDeleteSessionResponseType, ToACPClient, true,
	)

	// CloseSessionRequestMsg is the message type for requesting to close a
	// session.
	CloseSessionRequestMsg = newMessageTypeT[acp.CloseSessionRequest](
		CloseSessionRequestType, ToAgent, false,
	)

	// CloseSessionResponseMsg is the message type for the response to a
	// close-session request.
	CloseSessionResponseMsg = newMessageTypeT[acp.CloseSessionResponse](
		CloseSessionResponseType, ToACPClient, true,
	)

	// InitializeRequestMsg is the message type for an initialize request.
	InitializeRequestMsg = newMessageTypeT[acp.InitializeRequest](
		InitializeRequestType, ToAgent, false,
	)

	// InitializeResponseMsg is the message type for the response to an
	// initialize request.
	InitializeResponseMsg = newMessageTypeT[acp.InitializeResponse](
		InitializeResponseType, ToACPClient, true,
	)

	// NewSessionRequestMsg is the message type for requesting a new session.
	NewSessionRequestMsg = newMessageTypeT[acp.NewSessionRequest](
		NewSessionRequestType, ToAgent, false,
	)

	// NewSessionResponseMsg is the message type for the response to a
	// new-session request.
	NewSessionResponseMsg = newMessageTypeT[acp.NewSessionResponse](
		NewSessionResponseType, ToACPClient, true,
	)

	// AuthenticateRequestMsg is the message type for an authenticate request.
	AuthenticateRequestMsg = newMessageTypeT[acp.AuthenticateRequest](
		AuthenticateRequestType, ToAgent, false,
	)

	// AuthenticateResponseMsg is the message type for the response to an
	// authenticate request.
	AuthenticateResponseMsg = newMessageTypeT[acp.AuthenticateResponse](
		AuthenticateResponseType, ToACPClient, true,
	)

	// LoadSessionRequestMsg is the message type for requesting to load a
	// session.
	LoadSessionRequestMsg = newMessageTypeT[acp.LoadSessionRequest](
		LoadSessionRequestType, ToAgent, false,
	)

	// LoadSessionResponseMsg is the message type for the response to a
	// load-session request.
	LoadSessionResponseMsg = newMessageTypeT[acp.LoadSessionResponse](
		LoadSessionResponseType, ToACPClient, true,
	)

	// PromptRequestMsg is the message type for a prompt request.
	PromptRequestMsg = newMessageTypeT[acp.PromptRequest](
		PromptRequestType, ToAgent, false,
	)

	// PromptResponseMsg is the message type for the response to a prompt.
	PromptResponseMsg = newMessageTypeT[acp.PromptResponse](
		PromptResponseType, ToACPClient, true,
	)

	// CancelNotificationMsg is the message type for a cancel notification.
	CancelNotificationMsg = newMessageTypeT[acp.CancelNotification](
		CancelNotificationType, ToAgent, false,
	)
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
