package abyss

import (
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

// OnRequestPermissionRequestHandler is implemented by objects that provide the OnRequestPermissionRequest callback.
type OnRequestPermissionRequestHandler interface {
	OnRequestPermissionRequest(acp.RequestPermissionRequest) ([]*protobyss.ACPContainer, error)
}

// OnRequestPermissionResponseHandler is implemented by objects that provide the OnRequestPermissionResponse callback.
type OnRequestPermissionResponseHandler interface {
	OnRequestPermissionResponse(acp.RequestPermissionResponse) ([]*protobyss.ACPContainer, error)
}

// OnWriteTextFileRequestHandler is implemented by objects that provide the OnWriteTextFileRequest callback.
type OnWriteTextFileRequestHandler interface {
	OnWriteTextFileRequest(acp.WriteTextFileRequest) ([]*protobyss.ACPContainer, error)
}

// OnWriteTextFileResponseHandler is implemented by objects that provide the OnWriteTextFileResponse callback.
type OnWriteTextFileResponseHandler interface {
	OnWriteTextFileResponse(acp.WriteTextFileResponse) ([]*protobyss.ACPContainer, error)
}

// OnReadTextFileRequestHandler is implemented by objects that provide the OnReadTextFileRequest callback.
type OnReadTextFileRequestHandler interface {
	OnReadTextFileRequest(acp.ReadTextFileRequest) ([]*protobyss.ACPContainer, error)
}

// OnCreateTerminalRequestHandler is implemented by objects that provide the OnCreateTerminalRequest callback.
type OnCreateTerminalRequestHandler interface {
	OnCreateTerminalRequest(acp.CreateTerminalRequest) ([]*protobyss.ACPContainer, error)
}

// OnTerminalOutputRequestHandler is implemented by objects that provide the OnTerminalOutputRequest callback.
type OnTerminalOutputRequestHandler interface {
	OnTerminalOutputRequest(acp.TerminalOutputRequest) ([]*protobyss.ACPContainer, error)
}

// OnReleaseTerminalRequestHandler is implemented by objects that provide the OnReleaseTerminalRequest callback.
type OnReleaseTerminalRequestHandler interface {
	OnReleaseTerminalRequest(acp.ReleaseTerminalRequest) ([]*protobyss.ACPContainer, error)
}

// OnWaitForTerminalExitRequestHandler is implemented by objects that provide the OnWaitForTerminalExitRequest callback.
type OnWaitForTerminalExitRequestHandler interface {
	OnWaitForTerminalExitRequest(acp.WaitForTerminalExitRequest) ([]*protobyss.ACPContainer, error)
}

// OnKillTerminalRequestHandler is implemented by objects that provide the OnKillTerminalRequest callback.
type OnKillTerminalRequestHandler interface {
	OnKillTerminalRequest(acp.KillTerminalRequest) ([]*protobyss.ACPContainer, error)
}

// OnSessionNotificationHandler is implemented by objects that provide the OnSessionNotification callback.
type OnSessionNotificationHandler interface {
	OnSessionNotification(acp.SessionNotification) ([]*protobyss.ACPContainer, error)
}

// OnReadTextFileResponseHandler is implemented by objects that provide the OnReadTextFileResponse callback.
type OnReadTextFileResponseHandler interface {
	OnReadTextFileResponse(acp.ReadTextFileResponse) ([]*protobyss.ACPContainer, error)
}

// OnCreateTerminalResponseHandler is implemented by objects that provide the OnCreateTerminalResponse callback.
type OnCreateTerminalResponseHandler interface {
	OnCreateTerminalResponse(acp.CreateTerminalResponse) ([]*protobyss.ACPContainer, error)
}

// OnTerminalOutputResponseHandler is implemented by objects that provide the OnTerminalOutputResponse callback.
type OnTerminalOutputResponseHandler interface {
	OnTerminalOutputResponse(acp.TerminalOutputResponse) ([]*protobyss.ACPContainer, error)
}

// OnReleaseTerminalResponseHandler is implemented by objects that provide the OnReleaseTerminalResponse callback.
type OnReleaseTerminalResponseHandler interface {
	OnReleaseTerminalResponse(acp.ReleaseTerminalResponse) ([]*protobyss.ACPContainer, error)
}

// OnWaitForTerminalExitResponseHandler is implemented by objects that provide the OnWaitForTerminalExitResponse callback.
type OnWaitForTerminalExitResponseHandler interface {
	OnWaitForTerminalExitResponse(acp.WaitForTerminalExitResponse) ([]*protobyss.ACPContainer, error)
}

// OnKillTerminalResponseHandler is implemented by objects that provide the OnKillTerminalResponse callback.
type OnKillTerminalResponseHandler interface {
	OnKillTerminalResponse(acp.KillTerminalResponse) ([]*protobyss.ACPContainer, error)
}

// OnAuthenticateRequestHandler is implemented by objects that provide the OnAuthenticateRequest callback.
type OnAuthenticateRequestHandler interface {
	OnAuthenticateRequest(acp.AuthenticateRequest) ([]*protobyss.ACPContainer, error)
}

// OnInitializeRequestHandler is implemented by objects that provide the OnInitializeRequest callback.
type OnInitializeRequestHandler interface {
	OnInitializeRequest(acp.InitializeRequest) ([]*protobyss.ACPContainer, error)
}

// OnLogoutRequestHandler is implemented by objects that provide the OnLogoutRequest callback.
type OnLogoutRequestHandler interface {
	OnLogoutRequest(acp.LogoutRequest) ([]*protobyss.ACPContainer, error)
}

// OnCancelNotificationHandler is implemented by objects that provide the OnCancelNotification callback.
type OnCancelNotificationHandler interface {
	OnCancelNotification(acp.CancelNotification) ([]*protobyss.ACPContainer, error)
}

// OnCloseSessionRequestHandler is implemented by objects that provide the OnCloseSessionRequest callback.
type OnCloseSessionRequestHandler interface {
	OnCloseSessionRequest(acp.CloseSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnListSessionsRequestHandler is implemented by objects that provide the OnListSessionsRequest callback.
type OnListSessionsRequestHandler interface {
	OnListSessionsRequest(acp.ListSessionsRequest) ([]*protobyss.ACPContainer, error)
}

// OnNewSessionRequestHandler is implemented by objects that provide the OnNewSessionRequest callback.
type OnNewSessionRequestHandler interface {
	OnNewSessionRequest(acp.NewSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnPromptRequestHandler is implemented by objects that provide the OnPromptRequest callback.
type OnPromptRequestHandler interface {
	OnPromptRequest(acp.PromptRequest) ([]*protobyss.ACPContainer, error)
}

// OnResumeSessionRequestHandler is implemented by objects that provide the OnResumeSessionRequest callback.
type OnResumeSessionRequestHandler interface {
	OnResumeSessionRequest(acp.ResumeSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionConfigOptionRequestHandler is implemented by objects that provide the OnSetSessionConfigOptionRequest callback.
type OnSetSessionConfigOptionRequestHandler interface {
	OnSetSessionConfigOptionRequest(acp.SetSessionConfigOptionRequest) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionModeRequestHandler is implemented by objects that provide the OnSetSessionModeRequest callback.
type OnSetSessionModeRequestHandler interface {
	OnSetSessionModeRequest(acp.SetSessionModeRequest) ([]*protobyss.ACPContainer, error)
}

// OnLoadSessionRequestHandler is implemented by objects that provide the OnLoadSessionRequest callback.
type OnLoadSessionRequestHandler interface {
	OnLoadSessionRequest(acp.LoadSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnLoadSessionResponseHandler is implemented by objects that provide the OnLoadSessionResponse callback.
type OnLoadSessionResponseHandler interface {
	OnLoadSessionResponse(acp.LoadSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionModeResponseHandler is implemented by objects that provide the OnSetSessionModeResponse callback.
type OnSetSessionModeResponseHandler interface {
	OnSetSessionModeResponse(acp.SetSessionModeResponse) ([]*protobyss.ACPContainer, error)
}

// OnListSessionsResponseHandler is implemented by objects that provide the OnListSessionsResponse callback.
type OnListSessionsResponseHandler interface {
	OnListSessionsResponse(acp.ListSessionsResponse) ([]*protobyss.ACPContainer, error)
}

// OnResumeSessionResponseHandler is implemented by objects that provide the OnResumeSessionResponse callback.
type OnResumeSessionResponseHandler interface {
	OnResumeSessionResponse(acp.ResumeSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnSetSessionConfigOptionResponseHandler is implemented by objects that provide the OnSetSessionConfigOptionResponse callback.
type OnSetSessionConfigOptionResponseHandler interface {
	OnSetSessionConfigOptionResponse(acp.SetSessionConfigOptionResponse) ([]*protobyss.ACPContainer, error)
}

// OnLogoutResponseHandler is implemented by objects that provide the OnLogoutResponse callback.
type OnLogoutResponseHandler interface {
	OnLogoutResponse(acp.LogoutResponse) ([]*protobyss.ACPContainer, error)
}

// OnCloseSessionResponseHandler is implemented by objects that provide the OnCloseSessionResponse callback.
type OnCloseSessionResponseHandler interface {
	OnCloseSessionResponse(acp.CloseSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnInitializeResponseHandler is implemented by objects that provide the OnInitializeResponse callback.
type OnInitializeResponseHandler interface {
	OnInitializeResponse(acp.InitializeResponse) ([]*protobyss.ACPContainer, error)
}

// OnNewSessionResponseHandler is implemented by objects that provide the OnNewSessionResponse callback.
type OnNewSessionResponseHandler interface {
	OnNewSessionResponse(acp.NewSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnAuthenticateResponseHandler is implemented by objects that provide the OnAuthenticateResponse callback.
type OnAuthenticateResponseHandler interface {
	OnAuthenticateResponse(acp.AuthenticateResponse) ([]*protobyss.ACPContainer, error)
}

// OnPromptResponseHandler is implemented by objects that provide the OnPromptResponse callback.
type OnPromptResponseHandler interface {
	OnPromptResponse(acp.PromptResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidChangeDocumentNotificationHandler is implemented by objects that provide the OnUnstableDidChangeDocumentNotification callback.
type OnUnstableDidChangeDocumentNotificationHandler interface {
	OnUnstableDidChangeDocumentNotification(acp.UnstableDidChangeDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidCloseDocumentNotificationHandler is implemented by objects that provide the OnUnstableDidCloseDocumentNotification callback.
type OnUnstableDidCloseDocumentNotificationHandler interface {
	OnUnstableDidCloseDocumentNotification(acp.UnstableDidCloseDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidFocusDocumentNotificationHandler is implemented by objects that provide the OnUnstableDidFocusDocumentNotification callback.
type OnUnstableDidFocusDocumentNotificationHandler interface {
	OnUnstableDidFocusDocumentNotification(acp.UnstableDidFocusDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidOpenDocumentNotificationHandler is implemented by objects that provide the OnUnstableDidOpenDocumentNotification callback.
type OnUnstableDidOpenDocumentNotificationHandler interface {
	OnUnstableDidOpenDocumentNotification(acp.UnstableDidOpenDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDidSaveDocumentNotificationHandler is implemented by objects that provide the OnUnstableDidSaveDocumentNotification callback.
type OnUnstableDidSaveDocumentNotificationHandler interface {
	OnUnstableDidSaveDocumentNotification(acp.UnstableDidSaveDocumentNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableAcceptNesNotificationHandler is implemented by objects that provide the OnUnstableAcceptNesNotification callback.
type OnUnstableAcceptNesNotificationHandler interface {
	OnUnstableAcceptNesNotification(acp.UnstableAcceptNesNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableCloseNesRequestHandler is implemented by objects that provide the OnUnstableCloseNesRequest callback.
type OnUnstableCloseNesRequestHandler interface {
	OnUnstableCloseNesRequest(acp.UnstableCloseNesRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableRejectNesNotificationHandler is implemented by objects that provide the OnUnstableRejectNesNotification callback.
type OnUnstableRejectNesNotificationHandler interface {
	OnUnstableRejectNesNotification(acp.UnstableRejectNesNotification) ([]*protobyss.ACPContainer, error)
}

// OnUnstableStartNesRequestHandler is implemented by objects that provide the OnUnstableStartNesRequest callback.
type OnUnstableStartNesRequestHandler interface {
	OnUnstableStartNesRequest(acp.UnstableStartNesRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSuggestNesRequestHandler is implemented by objects that provide the OnUnstableSuggestNesRequest callback.
type OnUnstableSuggestNesRequestHandler interface {
	OnUnstableSuggestNesRequest(acp.UnstableSuggestNesRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDisableProviderRequestHandler is implemented by objects that provide the OnUnstableDisableProviderRequest callback.
type OnUnstableDisableProviderRequestHandler interface {
	OnUnstableDisableProviderRequest(acp.UnstableDisableProviderRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableListProvidersRequestHandler is implemented by objects that provide the OnUnstableListProvidersRequest callback.
type OnUnstableListProvidersRequestHandler interface {
	OnUnstableListProvidersRequest(acp.UnstableListProvidersRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSetProviderRequestHandler is implemented by objects that provide the OnUnstableSetProviderRequest callback.
type OnUnstableSetProviderRequestHandler interface {
	OnUnstableSetProviderRequest(acp.UnstableSetProviderRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDeleteSessionRequestHandler is implemented by objects that provide the OnUnstableDeleteSessionRequest callback.
type OnUnstableDeleteSessionRequestHandler interface {
	OnUnstableDeleteSessionRequest(acp.UnstableDeleteSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableForkSessionRequestHandler is implemented by objects that provide the OnUnstableForkSessionRequest callback.
type OnUnstableForkSessionRequestHandler interface {
	OnUnstableForkSessionRequest(acp.UnstableForkSessionRequest) ([]*protobyss.ACPContainer, error)
}

// OnUnstableForkSessionResponseHandler is implemented by objects that provide the OnUnstableForkSessionResponse callback.
type OnUnstableForkSessionResponseHandler interface {
	OnUnstableForkSessionResponse(acp.UnstableForkSessionResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableCloseNesResponseHandler is implemented by objects that provide the OnUnstableCloseNesResponse callback.
type OnUnstableCloseNesResponseHandler interface {
	OnUnstableCloseNesResponse(acp.UnstableCloseNesResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableStartNesResponseHandler is implemented by objects that provide the OnUnstableStartNesResponse callback.
type OnUnstableStartNesResponseHandler interface {
	OnUnstableStartNesResponse(acp.UnstableStartNesResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSuggestNesResponseHandler is implemented by objects that provide the OnUnstableSuggestNesResponse callback.
type OnUnstableSuggestNesResponseHandler interface {
	OnUnstableSuggestNesResponse(acp.UnstableSuggestNesResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDisableProviderResponseHandler is implemented by objects that provide the OnUnstableDisableProviderResponse callback.
type OnUnstableDisableProviderResponseHandler interface {
	OnUnstableDisableProviderResponse(acp.UnstableDisableProviderResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableListProvidersResponseHandler is implemented by objects that provide the OnUnstableListProvidersResponse callback.
type OnUnstableListProvidersResponseHandler interface {
	OnUnstableListProvidersResponse(acp.UnstableListProvidersResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableSetProviderResponseHandler is implemented by objects that provide the OnUnstableSetProviderResponse callback.
type OnUnstableSetProviderResponseHandler interface {
	OnUnstableSetProviderResponse(acp.UnstableSetProviderResponse) ([]*protobyss.ACPContainer, error)
}

// OnUnstableDeleteSessionResponseHandler is implemented by objects that provide the OnUnstableDeleteSessionResponse callback.
type OnUnstableDeleteSessionResponseHandler interface {
	OnUnstableDeleteSessionResponse(acp.UnstableDeleteSessionResponse) ([]*protobyss.ACPContainer, error)
}

// NewACPPluginRouter builds an ACPPluginRouter from an object that implements
// any subset of the *Handler interfaces. Unmatched callbacks are left nil.
func NewACPPluginRouter(v any) ACPPluginRouter {
	var r ACPPluginRouter

	if h, ok := v.(OnRequestPermissionRequestHandler); ok {
		r.OnRequestPermissionRequest = h.OnRequestPermissionRequest
	}
	if h, ok := v.(OnRequestPermissionResponseHandler); ok {
		r.OnRequestPermissionResponse = h.OnRequestPermissionResponse
	}
	if h, ok := v.(OnWriteTextFileRequestHandler); ok {
		r.OnWriteTextFileRequest = h.OnWriteTextFileRequest
	}
	if h, ok := v.(OnWriteTextFileResponseHandler); ok {
		r.OnWriteTextFileResponse = h.OnWriteTextFileResponse
	}
	if h, ok := v.(OnReadTextFileRequestHandler); ok {
		r.OnReadTextFileRequest = h.OnReadTextFileRequest
	}
	if h, ok := v.(OnCreateTerminalRequestHandler); ok {
		r.OnCreateTerminalRequest = h.OnCreateTerminalRequest
	}
	if h, ok := v.(OnTerminalOutputRequestHandler); ok {
		r.OnTerminalOutputRequest = h.OnTerminalOutputRequest
	}
	if h, ok := v.(OnReleaseTerminalRequestHandler); ok {
		r.OnReleaseTerminalRequest = h.OnReleaseTerminalRequest
	}
	if h, ok := v.(OnWaitForTerminalExitRequestHandler); ok {
		r.OnWaitForTerminalExitRequest = h.OnWaitForTerminalExitRequest
	}
	if h, ok := v.(OnKillTerminalRequestHandler); ok {
		r.OnKillTerminalRequest = h.OnKillTerminalRequest
	}
	if h, ok := v.(OnSessionNotificationHandler); ok {
		r.OnSessionNotification = h.OnSessionNotification
	}
	if h, ok := v.(OnReadTextFileResponseHandler); ok {
		r.OnReadTextFileResponse = h.OnReadTextFileResponse
	}
	if h, ok := v.(OnCreateTerminalResponseHandler); ok {
		r.OnCreateTerminalResponse = h.OnCreateTerminalResponse
	}
	if h, ok := v.(OnTerminalOutputResponseHandler); ok {
		r.OnTerminalOutputResponse = h.OnTerminalOutputResponse
	}
	if h, ok := v.(OnReleaseTerminalResponseHandler); ok {
		r.OnReleaseTerminalResponse = h.OnReleaseTerminalResponse
	}
	if h, ok := v.(OnWaitForTerminalExitResponseHandler); ok {
		r.OnWaitForTerminalExitResponse = h.OnWaitForTerminalExitResponse
	}
	if h, ok := v.(OnKillTerminalResponseHandler); ok {
		r.OnKillTerminalResponse = h.OnKillTerminalResponse
	}
	if h, ok := v.(OnAuthenticateRequestHandler); ok {
		r.OnAuthenticateRequest = h.OnAuthenticateRequest
	}
	if h, ok := v.(OnInitializeRequestHandler); ok {
		r.OnInitializeRequest = h.OnInitializeRequest
	}
	if h, ok := v.(OnLogoutRequestHandler); ok {
		r.OnLogoutRequest = h.OnLogoutRequest
	}
	if h, ok := v.(OnCancelNotificationHandler); ok {
		r.OnCancelNotification = h.OnCancelNotification
	}
	if h, ok := v.(OnCloseSessionRequestHandler); ok {
		r.OnCloseSessionRequest = h.OnCloseSessionRequest
	}
	if h, ok := v.(OnListSessionsRequestHandler); ok {
		r.OnListSessionsRequest = h.OnListSessionsRequest
	}
	if h, ok := v.(OnNewSessionRequestHandler); ok {
		r.OnNewSessionRequest = h.OnNewSessionRequest
	}
	if h, ok := v.(OnPromptRequestHandler); ok {
		r.OnPromptRequest = h.OnPromptRequest
	}
	if h, ok := v.(OnResumeSessionRequestHandler); ok {
		r.OnResumeSessionRequest = h.OnResumeSessionRequest
	}
	if h, ok := v.(OnSetSessionConfigOptionRequestHandler); ok {
		r.OnSetSessionConfigOptionRequest = h.OnSetSessionConfigOptionRequest
	}
	if h, ok := v.(OnSetSessionModeRequestHandler); ok {
		r.OnSetSessionModeRequest = h.OnSetSessionModeRequest
	}
	if h, ok := v.(OnLoadSessionRequestHandler); ok {
		r.OnLoadSessionRequest = h.OnLoadSessionRequest
	}
	if h, ok := v.(OnLoadSessionResponseHandler); ok {
		r.OnLoadSessionResponse = h.OnLoadSessionResponse
	}
	if h, ok := v.(OnSetSessionModeResponseHandler); ok {
		r.OnSetSessionModeResponse = h.OnSetSessionModeResponse
	}
	if h, ok := v.(OnListSessionsResponseHandler); ok {
		r.OnListSessionsResponse = h.OnListSessionsResponse
	}
	if h, ok := v.(OnResumeSessionResponseHandler); ok {
		r.OnResumeSessionResponse = h.OnResumeSessionResponse
	}
	if h, ok := v.(OnSetSessionConfigOptionResponseHandler); ok {
		r.OnSetSessionConfigOptionResponse = h.OnSetSessionConfigOptionResponse
	}
	if h, ok := v.(OnLogoutResponseHandler); ok {
		r.OnLogoutResponse = h.OnLogoutResponse
	}
	if h, ok := v.(OnCloseSessionResponseHandler); ok {
		r.OnCloseSessionResponse = h.OnCloseSessionResponse
	}
	if h, ok := v.(OnInitializeResponseHandler); ok {
		r.OnInitializeResponse = h.OnInitializeResponse
	}
	if h, ok := v.(OnNewSessionResponseHandler); ok {
		r.OnNewSessionResponse = h.OnNewSessionResponse
	}
	if h, ok := v.(OnAuthenticateResponseHandler); ok {
		r.OnAuthenticateResponse = h.OnAuthenticateResponse
	}
	if h, ok := v.(OnPromptResponseHandler); ok {
		r.OnPromptResponse = h.OnPromptResponse
	}
	if h, ok := v.(OnUnstableDidChangeDocumentNotificationHandler); ok {
		r.OnUnstableDidChangeDocumentNotification = h.OnUnstableDidChangeDocumentNotification
	}
	if h, ok := v.(OnUnstableDidCloseDocumentNotificationHandler); ok {
		r.OnUnstableDidCloseDocumentNotification = h.OnUnstableDidCloseDocumentNotification
	}
	if h, ok := v.(OnUnstableDidFocusDocumentNotificationHandler); ok {
		r.OnUnstableDidFocusDocumentNotification = h.OnUnstableDidFocusDocumentNotification
	}
	if h, ok := v.(OnUnstableDidOpenDocumentNotificationHandler); ok {
		r.OnUnstableDidOpenDocumentNotification = h.OnUnstableDidOpenDocumentNotification
	}
	if h, ok := v.(OnUnstableDidSaveDocumentNotificationHandler); ok {
		r.OnUnstableDidSaveDocumentNotification = h.OnUnstableDidSaveDocumentNotification
	}
	if h, ok := v.(OnUnstableAcceptNesNotificationHandler); ok {
		r.OnUnstableAcceptNesNotification = h.OnUnstableAcceptNesNotification
	}
	if h, ok := v.(OnUnstableCloseNesRequestHandler); ok {
		r.OnUnstableCloseNesRequest = h.OnUnstableCloseNesRequest
	}
	if h, ok := v.(OnUnstableRejectNesNotificationHandler); ok {
		r.OnUnstableRejectNesNotification = h.OnUnstableRejectNesNotification
	}
	if h, ok := v.(OnUnstableStartNesRequestHandler); ok {
		r.OnUnstableStartNesRequest = h.OnUnstableStartNesRequest
	}
	if h, ok := v.(OnUnstableSuggestNesRequestHandler); ok {
		r.OnUnstableSuggestNesRequest = h.OnUnstableSuggestNesRequest
	}
	if h, ok := v.(OnUnstableDisableProviderRequestHandler); ok {
		r.OnUnstableDisableProviderRequest = h.OnUnstableDisableProviderRequest
	}
	if h, ok := v.(OnUnstableListProvidersRequestHandler); ok {
		r.OnUnstableListProvidersRequest = h.OnUnstableListProvidersRequest
	}
	if h, ok := v.(OnUnstableSetProviderRequestHandler); ok {
		r.OnUnstableSetProviderRequest = h.OnUnstableSetProviderRequest
	}
	if h, ok := v.(OnUnstableDeleteSessionRequestHandler); ok {
		r.OnUnstableDeleteSessionRequest = h.OnUnstableDeleteSessionRequest
	}
	if h, ok := v.(OnUnstableForkSessionRequestHandler); ok {
		r.OnUnstableForkSessionRequest = h.OnUnstableForkSessionRequest
	}
	if h, ok := v.(OnUnstableForkSessionResponseHandler); ok {
		r.OnUnstableForkSessionResponse = h.OnUnstableForkSessionResponse
	}
	if h, ok := v.(OnUnstableCloseNesResponseHandler); ok {
		r.OnUnstableCloseNesResponse = h.OnUnstableCloseNesResponse
	}
	if h, ok := v.(OnUnstableStartNesResponseHandler); ok {
		r.OnUnstableStartNesResponse = h.OnUnstableStartNesResponse
	}
	if h, ok := v.(OnUnstableSuggestNesResponseHandler); ok {
		r.OnUnstableSuggestNesResponse = h.OnUnstableSuggestNesResponse
	}
	if h, ok := v.(OnUnstableDisableProviderResponseHandler); ok {
		r.OnUnstableDisableProviderResponse = h.OnUnstableDisableProviderResponse
	}
	if h, ok := v.(OnUnstableListProvidersResponseHandler); ok {
		r.OnUnstableListProvidersResponse = h.OnUnstableListProvidersResponse
	}
	if h, ok := v.(OnUnstableSetProviderResponseHandler); ok {
		r.OnUnstableSetProviderResponse = h.OnUnstableSetProviderResponse
	}
	if h, ok := v.(OnUnstableDeleteSessionResponseHandler); ok {
		r.OnUnstableDeleteSessionResponse = h.OnUnstableDeleteSessionResponse
	}

	return r
}
