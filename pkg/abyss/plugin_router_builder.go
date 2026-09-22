package abyss

import (
	"github.com/SethCurry/abyss/internal/timber"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

// OnRequestPermissionRequestHandler handles ACP request-permission-request
// callbacks.
type OnRequestPermissionRequestHandler interface {
	OnRequestPermissionRequest(acp.RequestPermissionRequest) ([]*protobyss.ACPContainer, error)
}

// OnRequestPermissionResponseHandler handles ACP request-permission-response
// callbacks.
type OnRequestPermissionResponseHandler interface {
	OnRequestPermissionResponse(acp.RequestPermissionResponse) ([]*protobyss.ACPContainer, error)
}

// OnWriteTextFileRequestHandler handles ACP write-text-file-request
// callbacks.
type OnWriteTextFileRequestHandler interface {
	OnWriteTextFileRequest(acp.WriteTextFileRequest) ([]*protobyss.ACPContainer, error)
}

// OnWriteTextFileResponseHandler handles ACP write-text-file-response
// callbacks.
type OnWriteTextFileResponseHandler interface {
	OnWriteTextFileResponse(acp.WriteTextFileResponse) ([]*protobyss.ACPContainer, error)
}

// OnReadTextFileRequestHandler handles ACP read-text-file-request callbacks.
type OnReadTextFileRequestHandler interface {
	OnReadTextFileRequest(acp.ReadTextFileRequest) ([]*protobyss.ACPContainer, error)
}

// OnCreateTerminalRequestHandler handles ACP create-terminal-request
// callbacks.
type OnCreateTerminalRequestHandler interface {
	OnCreateTerminalRequest(acp.CreateTerminalRequest) ([]*protobyss.ACPContainer, error)
}

// OnTerminalOutputRequestHandler handles ACP terminal-output-request
// callbacks.
type OnTerminalOutputRequestHandler interface {
	OnTerminalOutputRequest(acp.TerminalOutputRequest) ([]*protobyss.ACPContainer, error)
}

// OnReleaseTerminalRequestHandler handles ACP release-terminal-request
// callbacks.
type OnReleaseTerminalRequestHandler interface {
	OnReleaseTerminalRequest(acp.ReleaseTerminalRequest) ([]*protobyss.ACPContainer, error)
}

// OnWaitForTerminalExitRequestHandler handles ACP
// wait-for-terminal-exit-request callbacks.
type OnWaitForTerminalExitRequestHandler interface {
	OnWaitForTerminalExitRequest(acp.WaitForTerminalExitRequest) ([]*protobyss.ACPContainer, error)
}

// OnKillTerminalRequestHandler handles ACP kill-terminal-request callbacks.
type OnKillTerminalRequestHandler interface {
	OnKillTerminalRequest(acp.KillTerminalRequest) ([]*protobyss.ACPContainer, error)
}

// OnSessionNotificationHandler handles ACP session-notification callbacks.
type OnSessionNotificationHandler interface {
	OnSessionNotification(acp.SessionNotification) ([]*protobyss.ACPContainer, error)
}

// OnReadTextFileResponseHandler handles ACP read-text-file-response
// callbacks.
type OnReadTextFileResponseHandler interface {
	OnReadTextFileResponse(acp.ReadTextFileResponse) ([]*protobyss.ACPContainer, error)
}

// OnCreateTerminalResponseHandler handles ACP create-terminal-response
// callbacks.
type OnCreateTerminalResponseHandler interface {
	OnCreateTerminalResponse(acp.CreateTerminalResponse) ([]*protobyss.ACPContainer, error)
}

// OnTerminalOutputResponseHandler handles ACP terminal-output-response
// callbacks.
type OnTerminalOutputResponseHandler interface {
	OnTerminalOutputResponse(acp.TerminalOutputResponse) ([]*protobyss.ACPContainer, error)
}

// OnReleaseTerminalResponseHandler handles ACP release-terminal-response
// callbacks.
type OnReleaseTerminalResponseHandler interface {
	OnReleaseTerminalResponse(acp.ReleaseTerminalResponse) ([]*protobyss.ACPContainer, error)
}

// OnWaitForTerminalExitResponseHandler handles ACP
// wait-for-terminal-exit-response callbacks.
type OnWaitForTerminalExitResponseHandler interface {
	OnWaitForTerminalExitResponse(acp.WaitForTerminalExitResponse) ([]*protobyss.ACPContainer, error)
}

// OnKillTerminalResponseHandler handles ACP kill-terminal-response
// callbacks.
type OnKillTerminalResponseHandler interface {
	OnKillTerminalResponse(acp.KillTerminalResponse) ([]*protobyss.ACPContainer, error)
}

// OnAuthenticateRequestHandler handles ACP authenticate-request callbacks.
type OnAuthenticateRequestHandler interface {
	OnAuthenticateRequest(acp.AuthenticateRequest) ([]*protobyss.ACPContainer, error)
}

// OnInitializeRequestHandler handles ACP initialize-request callbacks.
type OnInitializeRequestHandler interface {
	OnInitializeRequest(acp.InitializeRequest) ([]*protobyss.ACPContainer, error)
}

// OnLogoutRequestHandler handles ACP logout-request callbacks.
type OnLogoutRequestHandler interface {
	OnLogoutRequest(acp.LogoutRequest) ([]*protobyss.ACPContainer, error)
}

// OnCancelNotificationHandler handles ACP cancel-notification callbacks.
type OnCancelNotificationHandler interface {
	OnCancelNotification(acp.CancelNotification) ([]*protobyss.ACPContainer, error)
}

// OnCloseSessionRequestHandler handles ACP close-session-request callbacks.
type OnCloseSessionRequestHandler interface {
	OnCloseSessionRequest(acp.CloseSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnListSessionsRequestHandler handles ACP list-sessions-request callbacks.
type OnListSessionsRequestHandler interface {
	OnListSessionsRequest(acp.ListSessionsRequest) ([]*protobyss.ACPContainer, error)
}

// OnNewSessionRequestHandler handles ACP new-session-request callbacks.
type OnNewSessionRequestHandler interface {
	OnNewSessionRequest(acp.NewSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnPromptRequestHandler handles ACP prompt-request callbacks.
type OnPromptRequestHandler interface {
	OnPromptRequest(acp.PromptRequest) ([]*protobyss.ACPContainer, error)
}

// OnResumeSessionRequestHandler handles ACP resume-session-request
// callbacks.
type OnResumeSessionRequestHandler interface {
	OnResumeSessionRequest(acp.ResumeSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionConfigOptionRequestHandler handles ACP
// set-session-config-option-request callbacks.
type OnSetSessionConfigOptionRequestHandler interface {
	OnSetSessionConfigOptionRequest(acp.SetSessionConfigOptionRequest) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionModeRequestHandler handles ACP set-session-mode-request
// callbacks.
type OnSetSessionModeRequestHandler interface {
	OnSetSessionModeRequest(acp.SetSessionModeRequest) ([]*protobyss.ACPContainer, error)
}

// OnLoadSessionRequestHandler handles ACP load-session-request callbacks.
type OnLoadSessionRequestHandler interface {
	OnLoadSessionRequest(acp.LoadSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnLoadSessionResponseHandler handles ACP load-session-response callbacks.
type OnLoadSessionResponseHandler interface {
	OnLoadSessionResponse(acp.LoadSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionModeResponseHandler handles ACP set-session-mode-response
// callbacks.
type OnSetSessionModeResponseHandler interface {
	OnSetSessionModeResponse(acp.SetSessionModeResponse) ([]*protobyss.ACPContainer, error)
}

// OnListSessionsResponseHandler handles ACP list-sessions-response
// callbacks.
type OnListSessionsResponseHandler interface {
	OnListSessionsResponse(acp.ListSessionsResponse) ([]*protobyss.ACPContainer, error)
}

// OnResumeSessionResponseHandler handles ACP resume-session-response
// callbacks.
type OnResumeSessionResponseHandler interface {
	OnResumeSessionResponse(acp.ResumeSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionConfigOptionResponseHandler handles ACP
// set-session-config-option-response callbacks.
type OnSetSessionConfigOptionResponseHandler interface {
	OnSetSessionConfigOptionResponse(acp.SetSessionConfigOptionResponse) ([]*protobyss.ACPContainer, error)
}

// OnLogoutResponseHandler handles ACP logout-response callbacks.
type OnLogoutResponseHandler interface {
	OnLogoutResponse(acp.LogoutResponse) ([]*protobyss.ACPContainer, error)
}

// OnCloseSessionResponseHandler handles ACP close-session-response
// callbacks.
type OnCloseSessionResponseHandler interface {
	OnCloseSessionResponse(acp.CloseSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnInitializeResponseHandler handles ACP initialize-response callbacks.
type OnInitializeResponseHandler interface {
	OnInitializeResponse(acp.InitializeResponse) ([]*protobyss.ACPContainer, error)
}

// OnNewSessionResponseHandler handles ACP new-session-response callbacks.
type OnNewSessionResponseHandler interface {
	OnNewSessionResponse(acp.NewSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnAuthenticateResponseHandler handles ACP authenticate-response callbacks.
type OnAuthenticateResponseHandler interface {
	OnAuthenticateResponse(acp.AuthenticateResponse) ([]*protobyss.ACPContainer, error)
}

// OnPromptResponseHandler handles ACP prompt-response callbacks.
type OnPromptResponseHandler interface {
	OnPromptResponse(acp.PromptResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidChangeDocumentNotificationHandler handles ACP
// unstable-did-change-document-notification callbacks.
type OnUnstableDidChangeDocumentNotificationHandler interface {
	OnUnstableDidChangeDocumentNotification(acp.UnstableDidChangeDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidCloseDocumentNotificationHandler handles ACP
// unstable-did-close-document-notification callbacks.
type OnUnstableDidCloseDocumentNotificationHandler interface {
	OnUnstableDidCloseDocumentNotification(acp.UnstableDidCloseDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidFocusDocumentNotificationHandler handles ACP
// unstable-did-focus-document-notification callbacks.
type OnUnstableDidFocusDocumentNotificationHandler interface {
	OnUnstableDidFocusDocumentNotification(acp.UnstableDidFocusDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidOpenDocumentNotificationHandler handles ACP
// unstable-did-open-document-notification callbacks.
type OnUnstableDidOpenDocumentNotificationHandler interface {
	OnUnstableDidOpenDocumentNotification(acp.UnstableDidOpenDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidSaveDocumentNotificationHandler handles ACP
// unstable-did-save-document-notification callbacks.
type OnUnstableDidSaveDocumentNotificationHandler interface {
	OnUnstableDidSaveDocumentNotification(acp.UnstableDidSaveDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableAcceptNesNotificationHandler handles ACP
// unstable-accept-nes-notification callbacks.
type OnUnstableAcceptNesNotificationHandler interface {
	OnUnstableAcceptNesNotification(acp.UnstableAcceptNesNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableCloseNesRequestHandler handles ACP unstable-close-nes-request
// callbacks.
type OnUnstableCloseNesRequestHandler interface {
	OnUnstableCloseNesRequest(acp.UnstableCloseNesRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableRejectNesNotificationHandler handles ACP
// unstable-reject-nes-notification callbacks.
type OnUnstableRejectNesNotificationHandler interface {
	OnUnstableRejectNesNotification(acp.UnstableRejectNesNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableStartNesRequestHandler handles ACP unstable-start-nes-request
// callbacks.
type OnUnstableStartNesRequestHandler interface {
	OnUnstableStartNesRequest(acp.UnstableStartNesRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSuggestNesRequestHandler handles ACP
// unstable-suggest-nes-request callbacks.
type OnUnstableSuggestNesRequestHandler interface {
	OnUnstableSuggestNesRequest(acp.UnstableSuggestNesRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDisableProviderRequestHandler handles ACP
// unstable-disable-provider-request callbacks.
type OnUnstableDisableProviderRequestHandler interface {
	OnUnstableDisableProviderRequest(acp.UnstableDisableProviderRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableListProvidersRequestHandler handles ACP
// unstable-list-providers-request callbacks.
type OnUnstableListProvidersRequestHandler interface {
	OnUnstableListProvidersRequest(acp.UnstableListProvidersRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSetProviderRequestHandler handles ACP
// unstable-set-provider-request callbacks.
type OnUnstableSetProviderRequestHandler interface {
	OnUnstableSetProviderRequest(acp.UnstableSetProviderRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDeleteSessionRequestHandler handles ACP
// unstable-delete-session-request callbacks.
type OnUnstableDeleteSessionRequestHandler interface {
	OnUnstableDeleteSessionRequest(acp.UnstableDeleteSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableForkSessionRequestHandler handles ACP
// unstable-fork-session-request callbacks.
type OnUnstableForkSessionRequestHandler interface {
	OnUnstableForkSessionRequest(acp.UnstableForkSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableForkSessionResponseHandler handles ACP
// unstable-fork-session-response callbacks.
type OnUnstableForkSessionResponseHandler interface {
	OnUnstableForkSessionResponse(acp.UnstableForkSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableCloseNesResponseHandler handles ACP unstable-close-nes-response
// callbacks.
type OnUnstableCloseNesResponseHandler interface {
	OnUnstableCloseNesResponse(acp.UnstableCloseNesResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableStartNesResponseHandler handles ACP unstable-start-nes-response
// callbacks.
type OnUnstableStartNesResponseHandler interface {
	OnUnstableStartNesResponse(acp.UnstableStartNesResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSuggestNesResponseHandler handles ACP
// unstable-suggest-nes-response callbacks.
type OnUnstableSuggestNesResponseHandler interface {
	OnUnstableSuggestNesResponse(acp.UnstableSuggestNesResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDisableProviderResponseHandler handles ACP
// unstable-disable-provider-response callbacks.
type OnUnstableDisableProviderResponseHandler interface {
	OnUnstableDisableProviderResponse(acp.UnstableDisableProviderResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableListProvidersResponseHandler handles ACP
// unstable-list-providers-response callbacks.
type OnUnstableListProvidersResponseHandler interface {
	OnUnstableListProvidersResponse(acp.UnstableListProvidersResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSetProviderResponseHandler handles ACP
// unstable-set-provider-response callbacks.
type OnUnstableSetProviderResponseHandler interface {
	OnUnstableSetProviderResponse(acp.UnstableSetProviderResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDeleteSessionResponseHandler handles ACP
// unstable-delete-session-response callbacks.
type OnUnstableDeleteSessionResponseHandler interface {
	OnUnstableDeleteSessionResponse(acp.UnstableDeleteSessionResponse) ([]*protobyss.ACPContainer, error)
}

// NewACPPluginRouter builds an ACPPluginRouter from an object that
// implements any subset of the *Handler interfaces. Unmatched callbacks
// are left nil.
func NewACPPluginRouter(v any) *ACPPluginRouter {
	r := &ACPPluginRouter{}
	var methods []string

	if h, ok := v.(OnRequestPermissionRequestHandler); ok {
		r.OnRequestPermissionRequest = h.OnRequestPermissionRequest
		methods = append(methods, "OnRequestPermissionRequest")
	}
	if h, ok := v.(OnRequestPermissionResponseHandler); ok {
		r.OnRequestPermissionResponse = h.OnRequestPermissionResponse
		methods = append(methods, "OnRequestPermissionResponse")
	}
	if h, ok := v.(OnWriteTextFileRequestHandler); ok {
		r.OnWriteTextFileRequest = h.OnWriteTextFileRequest
		methods = append(methods, "OnWriteTextFileRequest")
	}
	if h, ok := v.(OnWriteTextFileResponseHandler); ok {
		r.OnWriteTextFileResponse = h.OnWriteTextFileResponse
		methods = append(methods, "OnWriteTextFileResponse")
	}
	if h, ok := v.(OnReadTextFileRequestHandler); ok {
		r.OnReadTextFileRequest = h.OnReadTextFileRequest
		methods = append(methods, "OnReadTextFileRequest")
	}
	if h, ok := v.(OnCreateTerminalRequestHandler); ok {
		r.OnCreateTerminalRequest = h.OnCreateTerminalRequest
		methods = append(methods, "OnCreateTerminalRequest")
	}
	if h, ok := v.(OnTerminalOutputRequestHandler); ok {
		r.OnTerminalOutputRequest = h.OnTerminalOutputRequest
		methods = append(methods, "OnTerminalOutputRequest")
	}
	if h, ok := v.(OnReleaseTerminalRequestHandler); ok {
		r.OnReleaseTerminalRequest = h.OnReleaseTerminalRequest
		methods = append(methods, "OnReleaseTerminalRequest")
	}
	if h, ok := v.(OnWaitForTerminalExitRequestHandler); ok {
		r.OnWaitForTerminalExitRequest = h.OnWaitForTerminalExitRequest
		methods = append(methods, "OnWaitForTerminalExitRequest")
	}
	if h, ok := v.(OnKillTerminalRequestHandler); ok {
		r.OnKillTerminalRequest = h.OnKillTerminalRequest
		methods = append(methods, "OnKillTerminalRequest")
	}
	if h, ok := v.(OnSessionNotificationHandler); ok {
		r.OnSessionNotification = h.OnSessionNotification
		methods = append(methods, "OnSessionNotification")
	}
	if h, ok := v.(OnReadTextFileResponseHandler); ok {
		r.OnReadTextFileResponse = h.OnReadTextFileResponse
		methods = append(methods, "OnReadTextFileResponse")
	}
	if h, ok := v.(OnCreateTerminalResponseHandler); ok {
		r.OnCreateTerminalResponse = h.OnCreateTerminalResponse
		methods = append(methods, "OnCreateTerminalResponse")
	}
	if h, ok := v.(OnTerminalOutputResponseHandler); ok {
		r.OnTerminalOutputResponse = h.OnTerminalOutputResponse
		methods = append(methods, "OnTerminalOutputResponse")
	}
	if h, ok := v.(OnReleaseTerminalResponseHandler); ok {
		r.OnReleaseTerminalResponse = h.OnReleaseTerminalResponse
		methods = append(methods, "OnReleaseTerminalResponse")
	}
	if h, ok := v.(OnWaitForTerminalExitResponseHandler); ok {
		r.OnWaitForTerminalExitResponse = h.OnWaitForTerminalExitResponse
		methods = append(methods, "OnWaitForTerminalExitResponse")
	}
	if h, ok := v.(OnKillTerminalResponseHandler); ok {
		r.OnKillTerminalResponse = h.OnKillTerminalResponse
		methods = append(methods, "OnKillTerminalResponse")
	}
	if h, ok := v.(OnAuthenticateRequestHandler); ok {
		r.OnAuthenticateRequest = h.OnAuthenticateRequest
		methods = append(methods, "OnAuthenticateRequest")
	}
	if h, ok := v.(OnInitializeRequestHandler); ok {
		r.OnInitializeRequest = h.OnInitializeRequest
		methods = append(methods, "OnInitializeRequest")
	}
	if h, ok := v.(OnLogoutRequestHandler); ok {
		r.OnLogoutRequest = h.OnLogoutRequest
		methods = append(methods, "OnLogoutRequest")
	}
	if h, ok := v.(OnCancelNotificationHandler); ok {
		r.OnCancelNotification = h.OnCancelNotification
		methods = append(methods, "OnCancelNotification")
	}
	if h, ok := v.(OnCloseSessionRequestHandler); ok {
		r.OnCloseSessionRequest = h.OnCloseSessionRequest
		methods = append(methods, "OnCloseSessionRequest")
	}
	if h, ok := v.(OnListSessionsRequestHandler); ok {
		r.OnListSessionsRequest = h.OnListSessionsRequest
		methods = append(methods, "OnListSessionsRequest")
	}
	if h, ok := v.(OnNewSessionRequestHandler); ok {
		r.OnNewSessionRequest = h.OnNewSessionRequest
		methods = append(methods, "OnNewSessionRequest")
	}
	if h, ok := v.(OnPromptRequestHandler); ok {
		r.OnPromptRequest = h.OnPromptRequest
		methods = append(methods, "OnPromptRequest")
	}
	if h, ok := v.(OnResumeSessionRequestHandler); ok {
		r.OnResumeSessionRequest = h.OnResumeSessionRequest
		methods = append(methods, "OnResumeSessionRequest")
	}
	if h, ok := v.(OnSetSessionConfigOptionRequestHandler); ok {
		r.OnSetSessionConfigOptionRequest = h.OnSetSessionConfigOptionRequest
		methods = append(methods, "OnSetSessionConfigOptionRequest")
	}
	if h, ok := v.(OnSetSessionModeRequestHandler); ok {
		r.OnSetSessionModeRequest = h.OnSetSessionModeRequest
		methods = append(methods, "OnSetSessionModeRequest")
	}
	if h, ok := v.(OnLoadSessionRequestHandler); ok {
		r.OnLoadSessionRequest = h.OnLoadSessionRequest
		methods = append(methods, "OnLoadSessionRequest")
	}
	if h, ok := v.(OnLoadSessionResponseHandler); ok {
		r.OnLoadSessionResponse = h.OnLoadSessionResponse
		methods = append(methods, "OnLoadSessionResponse")
	}
	if h, ok := v.(OnSetSessionModeResponseHandler); ok {
		r.OnSetSessionModeResponse = h.OnSetSessionModeResponse
		methods = append(methods, "OnSetSessionModeResponse")
	}
	if h, ok := v.(OnListSessionsResponseHandler); ok {
		r.OnListSessionsResponse = h.OnListSessionsResponse
		methods = append(methods, "OnListSessionsResponse")
	}
	if h, ok := v.(OnResumeSessionResponseHandler); ok {
		r.OnResumeSessionResponse = h.OnResumeSessionResponse
		methods = append(methods, "OnResumeSessionResponse")
	}
	if h, ok := v.(OnSetSessionConfigOptionResponseHandler); ok {
		r.OnSetSessionConfigOptionResponse = h.OnSetSessionConfigOptionResponse
		methods = append(methods, "OnSetSessionConfigOptionResponse")
	}
	if h, ok := v.(OnLogoutResponseHandler); ok {
		r.OnLogoutResponse = h.OnLogoutResponse
		methods = append(methods, "OnLogoutResponse")
	}
	if h, ok := v.(OnCloseSessionResponseHandler); ok {
		r.OnCloseSessionResponse = h.OnCloseSessionResponse
		methods = append(methods, "OnCloseSessionResponse")
	}
	if h, ok := v.(OnInitializeResponseHandler); ok {
		r.OnInitializeResponse = h.OnInitializeResponse
		methods = append(methods, "OnInitializeResponse")
	}
	if h, ok := v.(OnNewSessionResponseHandler); ok {
		r.OnNewSessionResponse = h.OnNewSessionResponse
		methods = append(methods, "OnNewSessionResponse")
	}
	if h, ok := v.(OnAuthenticateResponseHandler); ok {
		r.OnAuthenticateResponse = h.OnAuthenticateResponse
		methods = append(methods, "OnAuthenticateResponse")
	}
	if h, ok := v.(OnPromptResponseHandler); ok {
		r.OnPromptResponse = h.OnPromptResponse
		methods = append(methods, "OnPromptResponse")
	}
	if h, ok := v.(OnUnstableDidChangeDocumentNotificationHandler); ok {
		r.OnUnstableDidChangeDocumentNotification = h.OnUnstableDidChangeDocumentNotification
		methods = append(methods, "OnUnstableDidChangeDocumentNotification")
	}
	if h, ok := v.(OnUnstableDidCloseDocumentNotificationHandler); ok {
		r.OnUnstableDidCloseDocumentNotification = h.OnUnstableDidCloseDocumentNotification
		methods = append(methods, "OnUnstableDidCloseDocumentNotification")
	}
	if h, ok := v.(OnUnstableDidFocusDocumentNotificationHandler); ok {
		r.OnUnstableDidFocusDocumentNotification = h.OnUnstableDidFocusDocumentNotification
		methods = append(methods, "OnUnstableDidFocusDocumentNotification")
	}
	if h, ok := v.(OnUnstableDidOpenDocumentNotificationHandler); ok {
		r.OnUnstableDidOpenDocumentNotification = h.OnUnstableDidOpenDocumentNotification
		methods = append(methods, "OnUnstableDidOpenDocumentNotification")
	}
	if h, ok := v.(OnUnstableDidSaveDocumentNotificationHandler); ok {
		r.OnUnstableDidSaveDocumentNotification = h.OnUnstableDidSaveDocumentNotification
		methods = append(methods, "OnUnstableDidSaveDocumentNotification")
	}
	if h, ok := v.(OnUnstableAcceptNesNotificationHandler); ok {
		r.OnUnstableAcceptNesNotification = h.OnUnstableAcceptNesNotification
		methods = append(methods, "OnUnstableAcceptNesNotification")
	}
	if h, ok := v.(OnUnstableCloseNesRequestHandler); ok {
		r.OnUnstableCloseNesRequest = h.OnUnstableCloseNesRequest
		methods = append(methods, "OnUnstableCloseNesRequest")
	}
	if h, ok := v.(OnUnstableRejectNesNotificationHandler); ok {
		r.OnUnstableRejectNesNotification = h.OnUnstableRejectNesNotification
		methods = append(methods, "OnUnstableRejectNesNotification")
	}
	if h, ok := v.(OnUnstableStartNesRequestHandler); ok {
		r.OnUnstableStartNesRequest = h.OnUnstableStartNesRequest
		methods = append(methods, "OnUnstableStartNesRequest")
	}
	if h, ok := v.(OnUnstableSuggestNesRequestHandler); ok {
		r.OnUnstableSuggestNesRequest = h.OnUnstableSuggestNesRequest
		methods = append(methods, "OnUnstableSuggestNesRequest")
	}
	if h, ok := v.(OnUnstableDisableProviderRequestHandler); ok {
		r.OnUnstableDisableProviderRequest = h.OnUnstableDisableProviderRequest
		methods = append(methods, "OnUnstableDisableProviderRequest")
	}
	if h, ok := v.(OnUnstableListProvidersRequestHandler); ok {
		r.OnUnstableListProvidersRequest = h.OnUnstableListProvidersRequest
		methods = append(methods, "OnUnstableListProvidersRequest")
	}
	if h, ok := v.(OnUnstableSetProviderRequestHandler); ok {
		r.OnUnstableSetProviderRequest = h.OnUnstableSetProviderRequest
		methods = append(methods, "OnUnstableSetProviderRequest")
	}
	if h, ok := v.(OnUnstableDeleteSessionRequestHandler); ok {
		r.OnUnstableDeleteSessionRequest = h.OnUnstableDeleteSessionRequest
		methods = append(methods, "OnUnstableDeleteSessionRequest")
	}
	if h, ok := v.(OnUnstableForkSessionRequestHandler); ok {
		r.OnUnstableForkSessionRequest = h.OnUnstableForkSessionRequest
		methods = append(methods, "OnUnstableForkSessionRequest")
	}
	if h, ok := v.(OnUnstableForkSessionResponseHandler); ok {
		r.OnUnstableForkSessionResponse = h.OnUnstableForkSessionResponse
		methods = append(methods, "OnUnstableForkSessionResponse")
	}
	if h, ok := v.(OnUnstableCloseNesResponseHandler); ok {
		r.OnUnstableCloseNesResponse = h.OnUnstableCloseNesResponse
		methods = append(methods, "OnUnstableCloseNesResponse")
	}
	if h, ok := v.(OnUnstableStartNesResponseHandler); ok {
		r.OnUnstableStartNesResponse = h.OnUnstableStartNesResponse
		methods = append(methods, "OnUnstableStartNesResponse")
	}
	if h, ok := v.(OnUnstableSuggestNesResponseHandler); ok {
		r.OnUnstableSuggestNesResponse = h.OnUnstableSuggestNesResponse
		methods = append(methods, "OnUnstableSuggestNesResponse")
	}
	if h, ok := v.(OnUnstableDisableProviderResponseHandler); ok {
		r.OnUnstableDisableProviderResponse = h.OnUnstableDisableProviderResponse
		methods = append(methods, "OnUnstableDisableProviderResponse")
	}
	if h, ok := v.(OnUnstableListProvidersResponseHandler); ok {
		r.OnUnstableListProvidersResponse = h.OnUnstableListProvidersResponse
		methods = append(methods, "OnUnstableListProvidersResponse")
	}
	if h, ok := v.(OnUnstableSetProviderResponseHandler); ok {
		r.OnUnstableSetProviderResponse = h.OnUnstableSetProviderResponse
		methods = append(methods, "OnUnstableSetProviderResponse")
	}
	if h, ok := v.(OnUnstableDeleteSessionResponseHandler); ok {
		r.OnUnstableDeleteSessionResponse = h.OnUnstableDeleteSessionResponse
		methods = append(methods, "OnUnstableDeleteSessionResponse")
	}

	logger := timber.ComponentLogger("acp.plugin_router")
	logger.Info().Strs("methods", methods).Msg("enabled methods for plugin")

	return r
}
