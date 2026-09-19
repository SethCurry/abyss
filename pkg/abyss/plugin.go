package abyss

import (
	"encoding/json"

	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

type ACPPluginRouter struct {
	// Client capability requests (agent -> client).
	OnRequestPermissionRequest   func(acp.RequestPermissionRequest) ([]*protobyss.ACPContainer, error)
	OnWriteTextFileRequest       func(acp.WriteTextFileRequest) ([]*protobyss.ACPContainer, error)
	OnReadTextFileRequest        func(acp.ReadTextFileRequest) ([]*protobyss.ACPContainer, error)
	OnCreateTerminalRequest      func(acp.CreateTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnTerminalOutputRequest      func(acp.TerminalOutputRequest) ([]*protobyss.ACPContainer, error)
	OnReleaseTerminalRequest     func(acp.ReleaseTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnWaitForTerminalExitRequest func(acp.WaitForTerminalExitRequest) ([]*protobyss.ACPContainer, error)
	OnKillTerminalRequest        func(acp.KillTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnSessionNotification        func(acp.SessionNotification) ([]*protobyss.ACPContainer, error)

	// Agent requests (client -> agent).
	OnAuthenticateRequest           func(acp.AuthenticateRequest) ([]*protobyss.ACPContainer, error)
	OnInitializeRequest             func(acp.InitializeRequest) ([]*protobyss.ACPContainer, error)
	OnLogoutRequest                 func(acp.LogoutRequest) ([]*protobyss.ACPContainer, error)
	OnCancelNotification            func(acp.CancelNotification) ([]*protobyss.ACPContainer, error)
	OnCloseSessionRequest           func(acp.CloseSessionRequest) ([]*protobyss.ACPContainer, error)
	OnListSessionsRequest           func(acp.ListSessionsRequest) ([]*protobyss.ACPContainer, error)
	OnNewSessionRequest             func(acp.NewSessionRequest) ([]*protobyss.ACPContainer, error)
	OnPromptRequest                 func(acp.PromptRequest) ([]*protobyss.ACPContainer, error)
	OnResumeSessionRequest          func(acp.ResumeSessionRequest) ([]*protobyss.ACPContainer, error)
	OnSetSessionConfigOptionRequest func(acp.SetSessionConfigOptionRequest) ([]*protobyss.ACPContainer, error)
	OnSetSessionModeRequest         func(acp.SetSessionModeRequest) ([]*protobyss.ACPContainer, error)
	OnLoadSessionRequest            func(acp.LoadSessionRequest) ([]*protobyss.ACPContainer, error)

	// Experimental agent requests (client -> agent).
	OnUnstableDidChangeDocumentNotification func(acp.UnstableDidChangeDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidCloseDocumentNotification  func(acp.UnstableDidCloseDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidFocusDocumentNotification  func(acp.UnstableDidFocusDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidOpenDocumentNotification   func(acp.UnstableDidOpenDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidSaveDocumentNotification   func(acp.UnstableDidSaveDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableAcceptNesNotification         func(acp.UnstableAcceptNesNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableCloseNesRequest               func(acp.UnstableCloseNesRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableRejectNesNotification         func(acp.UnstableRejectNesNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableStartNesRequest               func(acp.UnstableStartNesRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableSuggestNesRequest             func(acp.UnstableSuggestNesRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableDisableProviderRequest        func(acp.UnstableDisableProviderRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableListProvidersRequest          func(acp.UnstableListProvidersRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableSetProviderRequest            func(acp.UnstableSetProviderRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableDeleteSessionRequest          func(acp.UnstableDeleteSessionRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableForkSessionRequest            func(acp.UnstableForkSessionRequest) ([]*protobyss.ACPContainer, error)
}

func (r *ACPPluginRouter) Handle(msg *protobyss.ACPContainer) ([]*protobyss.ACPContainer, error) {
	switch MessageTypeID(msg.TypeId) {
	// Client capability requests (agent -> client).
	case RequestPermissionRequestType:
		return handleRequest(msg, r.OnRequestPermissionRequest)
	case WriteTextFileRequestType:
		return handleRequest(msg, r.OnWriteTextFileRequest)
	case ReadTextFileRequestType:
		return handleRequest(msg, r.OnReadTextFileRequest)
	case CreateTerminalRequestType:
		return handleRequest(msg, r.OnCreateTerminalRequest)
	case TerminalOutputRequestType:
		return handleRequest(msg, r.OnTerminalOutputRequest)
	case ReleaseTerminalRequestType:
		return handleRequest(msg, r.OnReleaseTerminalRequest)
	case WaitForTerminalExitRequestType:
		return handleRequest(msg, r.OnWaitForTerminalExitRequest)
	case KillTerminalRequestType:
		return handleRequest(msg, r.OnKillTerminalRequest)
	case SessionNotificationType:
		return handleRequest(msg, r.OnSessionNotification)

	// Agent requests (client -> agent).
	case AuthenticateRequestType:
		return handleRequest(msg, r.OnAuthenticateRequest)
	case InitializeRequestType:
		return handleRequest(msg, r.OnInitializeRequest)
	case LogoutRequestType:
		return handleRequest(msg, r.OnLogoutRequest)
	case CancelNotificationType:
		return handleRequest(msg, r.OnCancelNotification)
	case CloseSessionRequestType:
		return handleRequest(msg, r.OnCloseSessionRequest)
	case ListSessionsRequestType:
		return handleRequest(msg, r.OnListSessionsRequest)
	case NewSessionRequestType:
		return handleRequest(msg, r.OnNewSessionRequest)
	case PromptRequestType:
		return handleRequest(msg, r.OnPromptRequest)
	case ResumeSessionRequestType:
		return handleRequest(msg, r.OnResumeSessionRequest)
	case SetSessionConfigOptionRequestType:
		return handleRequest(msg, r.OnSetSessionConfigOptionRequest)
	case SetSessionModeRequestType:
		return handleRequest(msg, r.OnSetSessionModeRequest)
	case LoadSessionRequestType:
		return handleRequest(msg, r.OnLoadSessionRequest)

	// Experimental agent requests (client -> agent).
	case UnstableDidChangeDocumentNotificationType:
		return handleRequest(msg, r.OnUnstableDidChangeDocumentNotification)
	case UnstableDidCloseDocumentNotificationType:
		return handleRequest(msg, r.OnUnstableDidCloseDocumentNotification)
	case UnstableDidFocusDocumentNotificationType:
		return handleRequest(msg, r.OnUnstableDidFocusDocumentNotification)
	case UnstableDidOpenDocumentNotificationType:
		return handleRequest(msg, r.OnUnstableDidOpenDocumentNotification)
	case UnstableDidSaveDocumentNotificationType:
		return handleRequest(msg, r.OnUnstableDidSaveDocumentNotification)
	case UnstableAcceptNesNotificationType:
		return handleRequest(msg, r.OnUnstableAcceptNesNotification)
	case UnstableCloseNesRequestType:
		return handleRequest(msg, r.OnUnstableCloseNesRequest)
	case UnstableRejectNesNotificationType:
		return handleRequest(msg, r.OnUnstableRejectNesNotification)
	case UnstableStartNesRequestType:
		return handleRequest(msg, r.OnUnstableStartNesRequest)
	case UnstableSuggestNesRequestType:
		return handleRequest(msg, r.OnUnstableSuggestNesRequest)
	case UnstableDisableProviderRequestType:
		return handleRequest(msg, r.OnUnstableDisableProviderRequest)
	case UnstableListProvidersRequestType:
		return handleRequest(msg, r.OnUnstableListProvidersRequest)
	case UnstableSetProviderRequestType:
		return handleRequest(msg, r.OnUnstableSetProviderRequest)
	case UnstableDeleteSessionRequestType:
		return handleRequest(msg, r.OnUnstableDeleteSessionRequest)
	case UnstableForkSessionRequestType:
		return handleRequest(msg, r.OnUnstableForkSessionRequest)

	default:
		// Pass unhandled messages through unchanged.
		return []*protobyss.ACPContainer{msg}, nil
	}
}

// handleRequest unmarshals msg into T, invokes fn, and returns its response.
func handleRequest[T any](msg *protobyss.ACPContainer, fn func(T) ([]*protobyss.ACPContainer, error)) ([]*protobyss.ACPContainer, error) {
	if fn == nil {
		return []*protobyss.ACPContainer{msg}, nil
	}
	var params T
	if err := json.Unmarshal(msg.Content, &params); err != nil {
		return nil, err
	}

	return fn(params)
}
