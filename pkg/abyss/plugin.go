package abyss

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
	"github.com/google/uuid"
)

var _ protobyss.ACPPlugin = (*ACPPluginRouter)(nil)

// ACPPluginRouter wraps a struct that implements one or more event handler and manages
// translating protobyss.ACPContainer messages into appropriate acp structs and then
// dispatching ACP messages to the appropriate handler if configured.
type ACPPluginRouter struct {
	// Client capability requests (agent -> client).
	OnRequestPermissionRequest    func(acp.RequestPermissionRequest) ([]*protobyss.ACPContainer, error)
	OnRequestPermissionResponse   func(acp.RequestPermissionResponse) ([]*protobyss.ACPContainer, error)
	OnWriteTextFileRequest        func(acp.WriteTextFileRequest) ([]*protobyss.ACPContainer, error)
	OnWriteTextFileResponse       func(acp.WriteTextFileResponse) ([]*protobyss.ACPContainer, error)
	OnReadTextFileRequest         func(acp.ReadTextFileRequest) ([]*protobyss.ACPContainer, error)
	OnCreateTerminalRequest       func(acp.CreateTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnTerminalOutputRequest       func(acp.TerminalOutputRequest) ([]*protobyss.ACPContainer, error)
	OnReleaseTerminalRequest      func(acp.ReleaseTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnWaitForTerminalExitRequest  func(acp.WaitForTerminalExitRequest) ([]*protobyss.ACPContainer, error)
	OnKillTerminalRequest         func(acp.KillTerminalRequest) ([]*protobyss.ACPContainer, error)
	OnSessionNotification         func(acp.SessionNotification) ([]*protobyss.ACPContainer, error)
	OnReadTextFileResponse        func(acp.ReadTextFileResponse) ([]*protobyss.ACPContainer, error)
	OnCreateTerminalResponse      func(acp.CreateTerminalResponse) ([]*protobyss.ACPContainer, error)
	OnTerminalOutputResponse      func(acp.TerminalOutputResponse) ([]*protobyss.ACPContainer, error)
	OnReleaseTerminalResponse     func(acp.ReleaseTerminalResponse) ([]*protobyss.ACPContainer, error)
	OnWaitForTerminalExitResponse func(acp.WaitForTerminalExitResponse) ([]*protobyss.ACPContainer, error)
	OnKillTerminalResponse        func(acp.KillTerminalResponse) ([]*protobyss.ACPContainer, error)

	// Agent requests (client -> agent).
	OnAuthenticateRequest            func(acp.AuthenticateRequest) ([]*protobyss.ACPContainer, error)
	OnInitializeRequest              func(acp.InitializeRequest) ([]*protobyss.ACPContainer, error)
	OnLogoutRequest                  func(acp.LogoutRequest) ([]*protobyss.ACPContainer, error)
	OnCancelNotification             func(acp.CancelNotification) ([]*protobyss.ACPContainer, error)
	OnCloseSessionRequest            func(acp.CloseSessionRequest) ([]*protobyss.ACPContainer, error)
	OnListSessionsRequest            func(acp.ListSessionsRequest) ([]*protobyss.ACPContainer, error)
	OnNewSessionRequest              func(acp.NewSessionRequest) ([]*protobyss.ACPContainer, error)
	OnPromptRequest                  func(acp.PromptRequest) ([]*protobyss.ACPContainer, error)
	OnResumeSessionRequest           func(acp.ResumeSessionRequest) ([]*protobyss.ACPContainer, error)
	OnSetSessionConfigOptionRequest  func(acp.SetSessionConfigOptionRequest) ([]*protobyss.ACPContainer, error)
	OnSetSessionModeRequest          func(acp.SetSessionModeRequest) ([]*protobyss.ACPContainer, error)
	OnLoadSessionRequest             func(acp.LoadSessionRequest) ([]*protobyss.ACPContainer, error)
	OnLoadSessionResponse            func(acp.LoadSessionResponse) ([]*protobyss.ACPContainer, error)
	OnSetSessionModeResponse         func(acp.SetSessionModeResponse) ([]*protobyss.ACPContainer, error)
	OnListSessionsResponse           func(acp.ListSessionsResponse) ([]*protobyss.ACPContainer, error)
	OnResumeSessionResponse          func(acp.ResumeSessionResponse) ([]*protobyss.ACPContainer, error)
	OnSetSessionConfigOptionResponse func(acp.SetSessionConfigOptionResponse) ([]*protobyss.ACPContainer, error)
	OnLogoutResponse                 func(acp.LogoutResponse) ([]*protobyss.ACPContainer, error)
	OnCloseSessionResponse           func(acp.CloseSessionResponse) ([]*protobyss.ACPContainer, error)
	OnInitializeResponse             func(acp.InitializeResponse) ([]*protobyss.ACPContainer, error)
	OnNewSessionResponse             func(acp.NewSessionResponse) ([]*protobyss.ACPContainer, error)
	OnAuthenticateResponse           func(acp.AuthenticateResponse) ([]*protobyss.ACPContainer, error)
	OnPromptResponse                 func(acp.PromptResponse) ([]*protobyss.ACPContainer, error)

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
	OnUnstableForkSessionResponse           func(acp.UnstableForkSessionResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableCloseNesResponse              func(acp.UnstableCloseNesResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableStartNesResponse              func(acp.UnstableStartNesResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableSuggestNesResponse            func(acp.UnstableSuggestNesResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableDisableProviderResponse       func(acp.UnstableDisableProviderResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableListProvidersResponse         func(acp.UnstableListProvidersResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableSetProviderResponse           func(acp.UnstableSetProviderResponse) ([]*protobyss.ACPContainer, error)
	OnUnstableDeleteSessionResponse         func(acp.UnstableDeleteSessionResponse) ([]*protobyss.ACPContainer, error)
}

func (r *ACPPluginRouter) HandleMessage(ctx context.Context, msg *protobyss.ACPContainer) (*protobyss.ACPContainerList, error) {
	switch MessageTypeID(msg.TypeId) {
	// Client capability requests (agent -> client).
	case RequestPermissionRequestType:
		return handleRequest(msg, r.OnRequestPermissionRequest)
	case RequestPermissionResponseType:
		return handleRequest(msg, r.OnRequestPermissionResponse)
	case WriteTextFileRequestType:
		return handleRequest(msg, r.OnWriteTextFileRequest)
	case WriteTextFileResponseType:
		return handleRequest(msg, r.OnWriteTextFileResponse)
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
	case ReadTextFileResponseType:
		return handleRequest(msg, r.OnReadTextFileResponse)
	case CreateTerminalResponseType:
		return handleRequest(msg, r.OnCreateTerminalResponse)
	case TerminalOutputResponseType:
		return handleRequest(msg, r.OnTerminalOutputResponse)
	case ReleaseTerminalResponseType:
		return handleRequest(msg, r.OnReleaseTerminalResponse)
	case WaitForTerminalExitResponseType:
		return handleRequest(msg, r.OnWaitForTerminalExitResponse)
	case KillTerminalResponseType:
		return handleRequest(msg, r.OnKillTerminalResponse)

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
	case LoadSessionResponseType:
		return handleRequest(msg, r.OnLoadSessionResponse)
	case SetSessionModeResponseType:
		return handleRequest(msg, r.OnSetSessionModeResponse)
	case ListSessionsResponseType:
		return handleRequest(msg, r.OnListSessionsResponse)
	case ResumeSessionResponseType:
		return handleRequest(msg, r.OnResumeSessionResponse)
	case SetSessionConfigOptionResponseType:
		return handleRequest(msg, r.OnSetSessionConfigOptionResponse)
	case LogoutResponseType:
		return handleRequest(msg, r.OnLogoutResponse)
	case CloseSessionResponseType:
		return handleRequest(msg, r.OnCloseSessionResponse)
	case InitializeResponseType:
		return handleRequest(msg, r.OnInitializeResponse)
	case NewSessionResponseType:
		return handleRequest(msg, r.OnNewSessionResponse)
	case AuthenticateResponseType:
		return handleRequest(msg, r.OnAuthenticateResponse)
	case PromptResponseType:
		return handleRequest(msg, r.OnPromptResponse)

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
	case UnstableForkSessionResponseType:
		return handleRequest(msg, r.OnUnstableForkSessionResponse)
	case UnstableCloseNesResponseType:
		return handleRequest(msg, r.OnUnstableCloseNesResponse)
	case UnstableStartNesResponseType:
		return handleRequest(msg, r.OnUnstableStartNesResponse)
	case UnstableSuggestNesResponseType:
		return handleRequest(msg, r.OnUnstableSuggestNesResponse)
	case UnstableDisableProviderResponseType:
		return handleRequest(msg, r.OnUnstableDisableProviderResponse)
	case UnstableListProvidersResponseType:
		return handleRequest(msg, r.OnUnstableListProvidersResponse)
	case UnstableSetProviderResponseType:
		return handleRequest(msg, r.OnUnstableSetProviderResponse)
	case UnstableDeleteSessionResponseType:
		return handleRequest(msg, r.OnUnstableDeleteSessionResponse)

	default:
		// Return an error for unhandled message types.
		// This allows testing completeness of this switch in unit tests.
		return nil, fmt.Errorf("unhandled message type: %d", msg.TypeId)
	}
}

// handleRequest unmarshals msg into T, invokes fn, and returns its response.
func handleRequest[T any](msg *protobyss.ACPContainer, fn func(T) ([]*protobyss.ACPContainer, error)) (*protobyss.ACPContainerList, error) {
	if fn == nil {
		return &protobyss.ACPContainerList{Containers: []*protobyss.ACPContainer{msg}}, nil
	}
	var params T
	if err := json.Unmarshal(msg.Content, &params); err != nil {
		return nil, err
	}

	resp, err := fn(params)
	if err != nil {
		return nil, err
	}

	for _, v := range resp {
		msgType, err := GetMessageTypeByID(v.TypeId)
		if err != nil {
			return nil, err
		}

		if msgType.IsResponse() && v.ResponseFor == "" {
			v.ResponseFor = msg.MessageId
		}

		if v.MessageId == "" {
			v.MessageId = uuid.NewString()
		}
	}

	return &protobyss.ACPContainerList{
		Containers: resp,
	}, nil
}
