package abyss

import (
	"encoding/json"

	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

// ACPPlugin is a typed view of an ACP plugin: each ACP message type is
// dispatched to a dedicated handler so plugin authors implement typed methods
// rather than a single raw HandleMessage switch.
type ACPPlugin interface {
	// Client capability requests (agent -> client).
	OnRequestPermissionRequest(acp.RequestPermissionRequest) ([]*protobyss.ACPContainer, error)
	OnWriteTextFileRequest(acp.WriteTextFileRequest) ([]*protobyss.ACPContainer, error)
	OnReadTextFileRequest(acp.ReadTextFileRequest) ([]*protobyss.ACPContainer, error)
	OnCreateTerminalRequest(acp.CreateTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnTerminalOutputRequest(acp.TerminalOutputRequest) ([]*protobyss.ACPContainer, error)
	OnReleaseTerminalRequest(acp.ReleaseTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnWaitForTerminalExitRequest(acp.WaitForTerminalExitRequest) ([]*protobyss.ACPContainer, error)
	OnKillTerminalRequest(acp.KillTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnSessionNotification(acp.SessionNotification) ([]*protobyss.ACPContainer, error)

	// Agent requests (client -> agent).
	OnAuthenticateRequest(acp.AuthenticateRequest) ([]*protobyss.ACPContainer, error)
	OnInitializeRequest(acp.InitializeRequest) ([]*protobyss.ACPContainer, error)
	OnLogoutRequest(acp.LogoutRequest) ([]*protobyss.ACPContainer, error)
	OnCancelNotification(acp.CancelNotification) error
	OnCloseSessionRequest(acp.CloseSessionRequest) ([]*protobyss.ACPContainer, error)
	OnListSessionsRequest(acp.ListSessionsRequest) ([]*protobyss.ACPContainer, error)
	OnNewSessionRequest(acp.NewSessionRequest) ([]*protobyss.ACPContainer, error)
	OnPromptRequest(acp.PromptRequest) ([]*protobyss.ACPContainer, error)
	OnResumeSessionRequest(acp.ResumeSessionRequest) ([]*protobyss.ACPContainer, error)
	OnSetSessionConfigOptionRequest(acp.SetSessionConfigOptionRequest) ([]*protobyss.ACPContainer, error)
	OnSetSessionModeRequest(acp.SetSessionModeRequest) ([]*protobyss.ACPContainer, error)
	OnLoadSessionRequest(acp.LoadSessionRequest) ([]*protobyss.ACPContainer, error)

	// Experimental agent requests (client -> agent).
	OnUnstableDidChangeDocumentNotification(acp.UnstableDidChangeDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidCloseDocumentNotification(acp.UnstableDidCloseDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidFocusDocumentNotification(acp.UnstableDidFocusDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidOpenDocumentNotification(acp.UnstableDidOpenDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableDidSaveDocumentNotification(acp.UnstableDidSaveDocumentNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableAcceptNesNotification(acp.UnstableAcceptNesNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableCloseNesRequest(acp.UnstableCloseNesRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableRejectNesNotification(acp.UnstableRejectNesNotification) ([]*protobyss.ACPContainer, error)
	OnUnstableStartNesRequest(acp.UnstableStartNesRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableSuggestNesRequest(acp.UnstableSuggestNesRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableDisableProviderRequest(acp.UnstableDisableProviderRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableListProvidersRequest(acp.UnstableListProvidersRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableSetProviderRequest(acp.UnstableSetProviderRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableDeleteSessionRequest(acp.UnstableDeleteSessionRequest) ([]*protobyss.ACPContainer, error)
	OnUnstableForkSessionRequest(acp.UnstableForkSessionRequest) ([]*protobyss.ACPContainer, error)
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
