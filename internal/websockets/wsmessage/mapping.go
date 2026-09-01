package wsmessage

import (
	"encoding/json"
	"fmt"

	"github.com/coder/acp-go-sdk"
)

var MessageTypeToMessage = map[MessageType]func() any{
	RequestPermissionRequestType: func() any {
		return &acp.RequestPermissionRequest{}
	},
	RequestPermissionResponseType: func() any {
		return &acp.RequestPermissionResponse{}
	},
	WriteTextFileRequestType: func() any {
		return &acp.WriteTextFileRequest{}
	},
	WriteTextFileResponseType: func() any {
		return &acp.WriteTextFileResponse{}
	},
	ReadTextFileRequestType: func() any {
		return &acp.ReadTextFileRequest{}
	},
	ReadTextFileResponseType: func() any {
		return &acp.ReadTextFileResponse{}
	},
	CreateTerminalRequestType: func() any {
		return &acp.CreateTerminalRequest{}
	},
	CreateTerminalResponseType: func() any {
		return &acp.CreateTerminalResponse{}
	},
	TerminalOutputRequestType: func() any {
		return &acp.TerminalOutputRequest{}
	},
	TerminalOutputResponseType: func() any {
		return &acp.TerminalOutputResponse{}
	},
	ReleaseTerminalRequestType: func() any {
		return &acp.ReleaseTerminalRequest{}
	},
	ReleaseTerminalResponseType: func() any {
		return &acp.ReleaseTerminalResponse{}
	},
	WaitForTerminalExitRequestType: func() any {
		return &acp.WaitForTerminalExitRequest{}
	},
	WaitForTerminalExitResponseType: func() any {
		return &acp.WaitForTerminalExitResponse{}
	},
	KillTerminalRequestType: func() any {
		return &acp.KillTerminalRequest{}
	},
	KillTerminalResponseType: func() any {
		return &acp.KillTerminalResponse{}
	},
	SessionNotificationType: func() any {
		return &acp.SessionNotification{}
	},
	SetSessionModeRequestType: func() any {
		return &acp.SetSessionModeRequest{}
	},
	SetSessionModeResponseType: func() any {
		return &acp.SetSessionModeResponse{}
	},
	UnstableForkSessionRequestType: func() any {
		return &acp.UnstableForkSessionRequest{}
	},
	UnstableForkSessionResponseType: func() any {
		return &acp.UnstableForkSessionResponse{}
	},
	ListSessionsRequestType: func() any {
		return &acp.ListSessionsRequest{}
	},
	ListSessionsResponseType: func() any {
		return &acp.ListSessionsResponse{}
	},
	ResumeSessionRequestType: func() any {
		return &acp.ResumeSessionRequest{}
	},
	ResumeSessionResponseType: func() any {
		return &acp.ResumeSessionResponse{}
	},
	SetSessionConfigOptionRequestType: func() any {
		return &acp.SetSessionConfigOptionRequest{}
	},
	SetSessionConfigOptionResponseType: func() any {
		return &acp.SetSessionConfigOptionResponse{}
	},
	LogoutRequestType: func() any {
		return &acp.LogoutRequest{}
	},
	LogoutResponseType: func() any {
		return &acp.LogoutResponse{}
	},
	UnstableCloseNesRequestType: func() any {
		return &acp.UnstableCloseNesRequest{}
	},
	UnstableCloseNesResponseType: func() any {
		return &acp.UnstableCloseNesResponse{}
	},
	UnstableStartNesRequestType: func() any {
		return &acp.UnstableStartNesRequest{}
	},
	UnstableStartNesResponseType: func() any {
		return &acp.UnstableStartNesResponse{}
	},
	UnstableSuggestNesRequestType: func() any {
		return &acp.UnstableSuggestNesRequest{}
	},
	UnstableSuggestNesResponseType: func() any {
		return &acp.UnstableSuggestNesResponse{}
	},
	UnstableAcceptNesNotificationType: func() any {
		return &acp.UnstableAcceptNesNotification{}
	},
	UnstableRejectNesNotificationType: func() any {
		return &acp.UnstableRejectNesNotification{}
	},
	UnstableDidChangeDocumentNotificationType: func() any {
		return &acp.UnstableDidChangeDocumentNotification{}
	},
	UnstableDidCloseDocumentNotificationType: func() any {
		return &acp.UnstableDidCloseDocumentNotification{}
	},
	UnstableDidFocusDocumentNotificationType: func() any {
		return &acp.UnstableDidFocusDocumentNotification{}
	},
	UnstableDidOpenDocumentNotificationType: func() any {
		return &acp.UnstableDidOpenDocumentNotification{}
	},
	UnstableDidSaveDocumentNotificationType: func() any {
		return &acp.UnstableDidSaveDocumentNotification{}
	},
	UnstableDisableProviderRequestType: func() any {
		return &acp.UnstableDisableProviderRequest{}
	},
	UnstableDisableProviderResponseType: func() any {
		return &acp.UnstableDisableProviderResponse{}
	},
	UnstableListProvidersRequestType: func() any {
		return &acp.UnstableListProvidersRequest{}
	},
	UnstableListProvidersResponseType: func() any {
		return &acp.UnstableListProvidersResponse{}
	},
	UnstableSetProviderRequestType: func() any {
		return &acp.UnstableSetProviderRequest{}
	},
	UnstableSetProviderResponseType: func() any {
		return &acp.UnstableSetProviderResponse{}
	},
	UnstableDeleteSessionRequestType: func() any {
		return &acp.UnstableDeleteSessionRequest{}
	},
	UnstableDeleteSessionResponseType: func() any {
		return &acp.UnstableDeleteSessionResponse{}
	},
	CloseSessionRequestType: func() any {
		return &acp.CloseSessionRequest{}
	},
	CloseSessionResponseType: func() any {
		return &acp.CloseSessionResponse{}
	},
	InitializeRequestType: func() any {
		return &acp.InitializeRequest{}
	},
	InitializeResponseType: func() any {
		return &acp.InitializeResponse{}
	},
	NewSessionRequestType: func() any {
		return &acp.NewSessionRequest{}
	},
	NewSessionResponseType: func() any {
		return &acp.NewSessionResponse{}
	},
	AuthenticateRequestType: func() any {
		return &acp.AuthenticateRequest{}
	},
	AuthenticateResponseType: func() any {
		return &acp.AuthenticateResponse{}
	},
	LoadSessionRequestType: func() any {
		return &acp.LoadSessionRequest{}
	},
	LoadSessionResponseType: func() any {
		return &acp.LoadSessionResponse{}
	},
	PromptRequestType: func() any {
		return &acp.PromptRequest{}
	},
	PromptResponseType: func() any {
		return &acp.PromptResponse{}
	},
	CancelNotificationType: func() any {
		return &acp.CancelNotification{}
	},
}

// unmarshalMessage decodes a raw websocket message into the acp struct that
// corresponds to the given MessageType.
func UnmarshalMessage(mt MessageType, data []byte) (any, error) {
	constructor, ok := MessageTypeToMessage[mt]
	if !ok {
		return nil, fmt.Errorf("unknown websocket message type: %d", mt)
	}

	target := constructor()
	if err := json.Unmarshal(data, target); err != nil {
		return nil, err
	}

	return target, nil
}
