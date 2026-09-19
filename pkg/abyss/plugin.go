package abyss

import (
	"encoding/json"

	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

type RequestPermissionRequestPlugin interface {
	OnRequestPermissionRequest(acp.RequestPermissionRequest) ([]*protobyss.ACPContainer, error)
}

type RequestPermissionResponsePlugin interface {
	OnRequestPermissionResponse(acp.RequestPermissionResponse) ([]*protobyss.ACPContainer, error)
}

type RequestPermissionPlugin interface {
	RequestPermissionRequestPlugin
	RequestPermissionResponsePlugin
}

type WriteTextFileRequestPlugin interface {
	OnWriteTextFileRequest(acp.WriteTextFileRequest) ([]*protobyss.ACPContainer, error)
}

type WriteTextFileResponsePlugin interface {
	OnWriteTextFileResponse(acp.WriteTextFileResponse) ([]*protobyss.ACPContainer, error)
}

type WriteTextFilePlugin interface {
	WriteTextFileRequestPlugin
	WriteTextFileResponsePlugin
}

type ReadTextFileRequestPlugin interface {
	OnReadTextFileRequest(acp.ReadTextFileRequest) ([]*protobyss.ACPContainer, error)
}

type ReadTextFileResponsePlugin interface {
	OnReadTextFileResponse(acp.ReadTextFileResponse) ([]*protobyss.ACPContainer, error)
}

type ReadTextFilePlugin interface {
	ReadTextFileRequestPlugin
	ReadTextFileResponsePlugin
}

type CreateTerminalRequestPlugin interface {
	OnCreateTerminalRequest(acp.CreateTerminalRequest) ([]*protobyss.ACPContainer, error)
}

type CreateTerminalResponsePlugin interface {
	OnCreateTerminalResponse(acp.CreateTerminalResponse) ([]*protobyss.ACPContainer, error)
}

type CreateTerminalPlugin interface {
	CreateTerminalRequestPlugin
	CreateTerminalResponsePlugin
}

type TerminalOutputRequestPlugin interface {
	OnTerminalOutputRequest(acp.TerminalOutputRequest) ([]*protobyss.ACPContainer, error)
}

type TerminalOutputResponsePlugin interface {
	OnTerminalOutputResponse(acp.TerminalOutputResponse) ([]*protobyss.ACPContainer, error)
}

type TerminalOutputPlugin interface {
	TerminalOutputRequestPlugin
	TerminalOutputResponsePlugin
}

type ReleaseTerminalRequestPlugin interface {
	OnReleaseTerminalRequest(acp.ReleaseTerminalRequest) ([]*protobyss.ACPContainer, error)
}

type ReleaseTerminalResponsePlugin interface {
	OnReleaseTerminalResponse(acp.ReleaseTerminalResponse) ([]*protobyss.ACPContainer, error)
}

type ReleaseTerminalPlugin interface {
	ReleaseTerminalRequestPlugin
	ReleaseTerminalResponsePlugin
}

type WaitForTerminalExitRequestPlugin interface {
	OnWaitForTerminalExitRequest(acp.WaitForTerminalExitRequest) ([]*protobyss.ACPContainer, error)
}

type WaitForTerminalExitResponsePlugin interface {
	OnWaitForTerminalExitResponse(acp.WaitForTerminalExitResponse) ([]*protobyss.ACPContainer, error)
}

type WaitForTerminalExitPlugin interface {
	WaitForTerminalExitRequestPlugin
	WaitForTerminalExitResponsePlugin
}

type KillTerminalRequestPlugin interface {
	OnKillTerminalRequest(acp.KillTerminalRequest) ([]*protobyss.ACPContainer, error)
}

type KillTerminalResponsePlugin interface {
	OnKillTerminalResponse(acp.KillTerminalResponse) ([]*protobyss.ACPContainer, error)
}

type KillTerminalPlugin interface {
	KillTerminalRequestPlugin
	KillTerminalResponsePlugin
}

type SessionNotificationPlugin interface {
	OnSessionNotification(acp.SessionNotification) ([]*protobyss.ACPContainer, error)
}

type AuthenticateRequestPlugin interface {
	OnAuthenticateRequest(acp.AuthenticateRequest) ([]*protobyss.ACPContainer, error)
}

type AuthenticateResponsePlugin interface {
	OnAuthenticateResponse(acp.AuthenticateResponse) ([]*protobyss.ACPContainer, error)
}

type AuthenticatePlugin interface {
	AuthenticateRequestPlugin
	AuthenticateResponsePlugin
}

type InitializeRequestPlugin interface {
	OnInitializeRequest(acp.InitializeRequest) ([]*protobyss.ACPContainer, error)
}

type InitializeResponsePlugin interface {
	OnInitializeResponse(acp.InitializeResponse) ([]*protobyss.ACPContainer, error)
}

type InitializePlugin interface {
	InitializeRequestPlugin
	InitializeResponsePlugin
}

type LogoutRequestPlugin interface {
	OnLogoutRequest(acp.LogoutRequest) ([]*protobyss.ACPContainer, error)
}

type LogoutResponsePlugin interface {
	OnLogoutResponse(acp.LogoutResponse) ([]*protobyss.ACPContainer, error)
}

type LogoutPlugin interface {
	LogoutRequestPlugin
	LogoutResponsePlugin
}

type CancelNotificationPlugin interface {
	OnCancelNotification(acp.CancelNotification) ([]*protobyss.ACPContainer, error)
}

type CloseSessionRequestPlugin interface {
	OnCloseSessionRequest(acp.CloseSessionRequest) ([]*protobyss.ACPContainer, error)
}

type CloseSessionResponsePlugin interface {
	OnCloseSessionResponse(acp.CloseSessionResponse) ([]*protobyss.ACPContainer, error)
}

type CloseSessionPlugin interface {
	CloseSessionRequestPlugin
	CloseSessionResponsePlugin
}

type ListSessionsRequestPlugin interface {
	OnListSessionsRequest(acp.ListSessionsRequest) ([]*protobyss.ACPContainer, error)
}

type ListSessionsResponsePlugin interface {
	OnListSessionsResponse(acp.ListSessionsResponse) ([]*protobyss.ACPContainer, error)
}

type ListSessionsPlugin interface {
	ListSessionsRequestPlugin
	ListSessionsResponsePlugin
}

type NewSessionRequestPlugin interface {
	OnNewSessionRequest(acp.NewSessionRequest) ([]*protobyss.ACPContainer, error)
}

type NewSessionResponsePlugin interface {
	OnNewSessionResponse(acp.NewSessionResponse) ([]*protobyss.ACPContainer, error)
}

type NewSessionPlugin interface {
	NewSessionRequestPlugin
	NewSessionResponsePlugin
}

type PromptRequestPlugin interface {
	OnPromptRequest(acp.PromptRequest) ([]*protobyss.ACPContainer, error)
}

type PromptResponsePlugin interface {
	OnPromptResponse(acp.PromptResponse) ([]*protobyss.ACPContainer, error)
}

type PromptPlugin interface {
	PromptRequestPlugin
	PromptResponsePlugin
}

type ResumeSessionRequestPlugin interface {
	OnResumeSessionRequest(acp.ResumeSessionRequest) ([]*protobyss.ACPContainer, error)
}

type ResumeSessionResponsePlugin interface {
	OnResumeSessionResponse(acp.ResumeSessionResponse) ([]*protobyss.ACPContainer, error)
}

type ResumeSessionPlugin interface {
	ResumeSessionRequestPlugin
	ResumeSessionResponsePlugin
}

type SetSessionConfigOptionRequestPlugin interface {
	OnSetSessionConfigOptionRequest(acp.SetSessionConfigOptionRequest) ([]*protobyss.ACPContainer, error)
}

type SetSessionConfigOptionResponsePlugin interface {
	OnSetSessionConfigOptionResponse(acp.SetSessionConfigOptionResponse) ([]*protobyss.ACPContainer, error)
}

type SetSessionConfigOptionPlugin interface {
	SetSessionConfigOptionRequestPlugin
	SetSessionConfigOptionResponsePlugin
}

type SetSessionModeRequestPlugin interface {
	OnSetSessionModeRequest(acp.SetSessionModeRequest) ([]*protobyss.ACPContainer, error)
}

type SetSessionModeResponsePlugin interface {
	OnSetSessionModeResponse(acp.SetSessionModeResponse) ([]*protobyss.ACPContainer, error)
}

type SetSessionModePlugin interface {
	SetSessionModeRequestPlugin
	SetSessionModeResponsePlugin
}

type LoadSessionRequestPlugin interface {
	OnLoadSessionRequest(acp.LoadSessionRequest) ([]*protobyss.ACPContainer, error)
}

type LoadSessionResponsePlugin interface {
	OnLoadSessionResponse(acp.LoadSessionResponse) ([]*protobyss.ACPContainer, error)
}

type LoadSessionPlugin interface {
	LoadSessionRequestPlugin
	LoadSessionResponsePlugin
}

type UnstableDidChangeDocumentNotificationPlugin interface {
	OnUnstableDidChangeDocumentNotification(acp.UnstableDidChangeDocumentNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableDidCloseDocumentNotificationPlugin interface {
	OnUnstableDidCloseDocumentNotification(acp.UnstableDidCloseDocumentNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableDidFocusDocumentNotificationPlugin interface {
	OnUnstableDidFocusDocumentNotification(acp.UnstableDidFocusDocumentNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableDidOpenDocumentNotificationPlugin interface {
	OnUnstableDidOpenDocumentNotification(acp.UnstableDidOpenDocumentNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableDidSaveDocumentNotificationPlugin interface {
	OnUnstableDidSaveDocumentNotification(acp.UnstableDidSaveDocumentNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableAcceptNesNotificationPlugin interface {
	OnUnstableAcceptNesNotification(acp.UnstableAcceptNesNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableCloseNesRequestPlugin interface {
	OnUnstableCloseNesRequest(acp.UnstableCloseNesRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableRejectNesNotificationPlugin interface {
	OnUnstableRejectNesNotification(acp.UnstableRejectNesNotification) ([]*protobyss.ACPContainer, error)
}

type UnstableStartNesRequestPlugin interface {
	OnUnstableStartNesRequest(acp.UnstableStartNesRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableSuggestNesRequestPlugin interface {
	OnUnstableSuggestNesRequest(acp.UnstableSuggestNesRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableDisableProviderRequestPlugin interface {
	OnUnstableDisableProviderRequest(acp.UnstableDisableProviderRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableListProvidersRequestPlugin interface {
	OnUnstableListProvidersRequest(acp.UnstableListProvidersRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableSetProviderRequestPlugin interface {
	OnUnstableSetProviderRequest(acp.UnstableSetProviderRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableDeleteSessionRequestPlugin interface {
	OnUnstableDeleteSessionRequest(acp.UnstableDeleteSessionRequest) ([]*protobyss.ACPContainer, error)
}

type UnstableForkSessionRequestPlugin interface {
	OnUnstableForkSessionRequest(acp.UnstableForkSessionRequest) ([]*protobyss.ACPContainer, error)
}

// ACPPlugin is a typed view of an ACP plugin: each ACP message type is
// dispatched to a dedicated handler so plugin authors implement typed methods
// rather than a single raw HandleMessage switch.
type ACPPlugin interface {
	// TODO this is missing response types
	// Client capability requests (agent -> client).
	// Client capability requests (agent -> client).
	RequestPermissionPlugin
	WriteTextFilePlugin
	ReadTextFilePlugin
	CreateTerminalPlugin
	TerminalOutputPlugin
	ReleaseTerminalPlugin
	WaitForTerminalExitPlugin
	KillTerminalPlugin
	SessionNotificationPlugin

	// Agent requests (client -> agent).
	AuthenticatePlugin
	InitializePlugin
	LogoutPlugin
	CancelNotificationPlugin
	CloseSessionPlugin
	ListSessionsPlugin
	NewSessionPlugin
	PromptPlugin
	ResumeSessionPlugin
	SetSessionConfigOptionPlugin
	SetSessionModePlugin
	LoadSessionPlugin

	// Experimental agent requests (client -> agent).
	UnstableDidChangeDocumentNotificationPlugin
	UnstableDidCloseDocumentNotificationPlugin
	UnstableDidFocusDocumentNotificationPlugin
	UnstableDidOpenDocumentNotificationPlugin
	UnstableDidSaveDocumentNotificationPlugin
	UnstableAcceptNesNotificationPlugin
	UnstableCloseNesRequestPlugin
	UnstableRejectNesNotificationPlugin
	UnstableStartNesRequestPlugin
	UnstableSuggestNesRequestPlugin
	UnstableDisableProviderRequestPlugin
	UnstableListProvidersRequestPlugin
	UnstableSetProviderRequestPlugin
	UnstableDeleteSessionRequestPlugin
	UnstableForkSessionRequestPlugin
}

type ACPPluginRouter struct {
	underlying ACPPlugin
}

func (r *ACPPluginRouter) Handle(msg *protobyss.ACPContainer) ([]*protobyss.ACPContainer, error) {
	switch MessageTypeID(msg.TypeId) {
	// Client capability requests (agent -> client).
	case RequestPermissionRequestType:
		return handleRequest(msg, r.underlying.OnRequestPermissionRequest)
	case WriteTextFileRequestType:
		return handleRequest(msg, r.underlying.OnWriteTextFileRequest)
	case ReadTextFileRequestType:
		return handleRequest(msg, r.underlying.OnReadTextFileRequest)
	case CreateTerminalRequestType:
		return handleRequest(msg, r.underlying.OnCreateTerminalRequest)
	case TerminalOutputRequestType:
		return handleRequest(msg, r.underlying.OnTerminalOutputRequest)
	case ReleaseTerminalRequestType:
		return handleRequest(msg, r.underlying.OnReleaseTerminalRequest)
	case WaitForTerminalExitRequestType:
		return handleRequest(msg, r.underlying.OnWaitForTerminalExitRequest)
	case KillTerminalRequestType:
		return handleRequest(msg, r.underlying.OnKillTerminalRequest)
	case SessionNotificationType:
		return handleRequest(msg, r.underlying.OnSessionNotification)

	// Agent requests (client -> agent).
	case AuthenticateRequestType:
		return handleRequest(msg, r.underlying.OnAuthenticateRequest)
	case InitializeRequestType:
		return handleRequest(msg, r.underlying.OnInitializeRequest)
	case LogoutRequestType:
		return handleRequest(msg, r.underlying.OnLogoutRequest)
	case CancelNotificationType:
		return handleRequest(msg, r.underlying.OnCancelNotification)
	case CloseSessionRequestType:
		return handleRequest(msg, r.underlying.OnCloseSessionRequest)
	case ListSessionsRequestType:
		return handleRequest(msg, r.underlying.OnListSessionsRequest)
	case NewSessionRequestType:
		return handleRequest(msg, r.underlying.OnNewSessionRequest)
	case PromptRequestType:
		return handleRequest(msg, r.underlying.OnPromptRequest)
	case ResumeSessionRequestType:
		return handleRequest(msg, r.underlying.OnResumeSessionRequest)
	case SetSessionConfigOptionRequestType:
		return handleRequest(msg, r.underlying.OnSetSessionConfigOptionRequest)
	case SetSessionModeRequestType:
		return handleRequest(msg, r.underlying.OnSetSessionModeRequest)
	case LoadSessionRequestType:
		return handleRequest(msg, r.underlying.OnLoadSessionRequest)

	// Experimental agent requests (client -> agent).
	case UnstableDidChangeDocumentNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableDidChangeDocumentNotification)
	case UnstableDidCloseDocumentNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableDidCloseDocumentNotification)
	case UnstableDidFocusDocumentNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableDidFocusDocumentNotification)
	case UnstableDidOpenDocumentNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableDidOpenDocumentNotification)
	case UnstableDidSaveDocumentNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableDidSaveDocumentNotification)
	case UnstableAcceptNesNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableAcceptNesNotification)
	case UnstableCloseNesRequestType:
		return handleRequest(msg, r.underlying.OnUnstableCloseNesRequest)
	case UnstableRejectNesNotificationType:
		return handleRequest(msg, r.underlying.OnUnstableRejectNesNotification)
	case UnstableStartNesRequestType:
		return handleRequest(msg, r.underlying.OnUnstableStartNesRequest)
	case UnstableSuggestNesRequestType:
		return handleRequest(msg, r.underlying.OnUnstableSuggestNesRequest)
	case UnstableDisableProviderRequestType:
		return handleRequest(msg, r.underlying.OnUnstableDisableProviderRequest)
	case UnstableListProvidersRequestType:
		return handleRequest(msg, r.underlying.OnUnstableListProvidersRequest)
	case UnstableSetProviderRequestType:
		return handleRequest(msg, r.underlying.OnUnstableSetProviderRequest)
	case UnstableDeleteSessionRequestType:
		return handleRequest(msg, r.underlying.OnUnstableDeleteSessionRequest)
	case UnstableForkSessionRequestType:
		return handleRequest(msg, r.underlying.OnUnstableForkSessionRequest)

	default:
		// Pass unhandled messages through unchanged.
		return []*protobyss.ACPContainer{msg}, nil
	}
}

// handleRequest unmarshals msg into T, invokes fn, and returns its response.
func handleRequest[T any](msg *protobyss.ACPContainer, fn func(T) ([]*protobyss.ACPContainer, error)) ([]*protobyss.ACPContainer, error) {
	var params T
	if err := json.Unmarshal(msg.Content, &params); err != nil {
		return nil, err
	}

	return fn(params)
}
