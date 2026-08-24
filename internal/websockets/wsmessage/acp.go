package wsmessage

import "github.com/coder/acp-go-sdk"

// acpMessageType is a type constraint that limits the types accepted by
// messageTypeFor to those handled by unmarshalMessage.
type ACPMessageType interface {
	acp.RequestPermissionRequest |
		acp.RequestPermissionResponse |
		acp.WriteTextFileRequest |
		acp.WriteTextFileResponse |
		acp.ReadTextFileRequest |
		acp.ReadTextFileResponse |
		acp.CreateTerminalRequest |
		acp.CreateTerminalResponse |
		acp.TerminalOutputRequest |
		acp.TerminalOutputResponse |
		acp.ReleaseTerminalRequest |
		acp.ReleaseTerminalResponse |
		acp.WaitForTerminalExitRequest |
		acp.WaitForTerminalExitResponse |
		acp.KillTerminalRequest |
		acp.KillTerminalResponse |
		acp.SessionNotification |
		acp.SetSessionModeRequest |
		acp.SetSessionModeResponse |
		acp.UnstableForkSessionRequest |
		acp.UnstableForkSessionResponse |
		acp.ListSessionsRequest |
		acp.ListSessionsResponse |
		acp.ResumeSessionRequest |
		acp.ResumeSessionResponse |
		acp.SetSessionConfigOptionRequest |
		acp.SetSessionConfigOptionResponse |
		acp.LogoutRequest |
		acp.LogoutResponse |
		acp.UnstableCloseNesRequest |
		acp.UnstableCloseNesResponse |
		acp.UnstableStartNesRequest |
		acp.UnstableStartNesResponse |
		acp.UnstableSuggestNesRequest |
		acp.UnstableSuggestNesResponse |
		acp.UnstableAcceptNesNotification |
		acp.UnstableRejectNesNotification |
		acp.UnstableDidChangeDocumentNotification |
		acp.UnstableDidCloseDocumentNotification |
		acp.UnstableDidFocusDocumentNotification |
		acp.UnstableDidOpenDocumentNotification |
		acp.UnstableDidSaveDocumentNotification |
		acp.UnstableDisableProviderRequest |
		acp.UnstableDisableProviderResponse |
		acp.UnstableListProvidersRequest |
		acp.UnstableListProvidersResponse |
		acp.UnstableSetProviderRequest |
		acp.UnstableSetProviderResponse |
		acp.UnstableDeleteSessionRequest |
		acp.UnstableDeleteSessionResponse |
		acp.CloseSessionRequest |
		acp.CloseSessionResponse |
		acp.InitializeRequest |
		acp.InitializeResponse |
		acp.NewSessionRequest |
		acp.NewSessionResponse |
		acp.AuthenticateRequest |
		acp.AuthenticateResponse |
		acp.LoadSessionRequest |
		acp.LoadSessionResponse |
		acp.PromptRequest |
		acp.PromptResponse |
		acp.CancelNotification
}
