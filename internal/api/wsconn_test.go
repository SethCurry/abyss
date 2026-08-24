package api

import (
	"testing"

	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
)

func TestClassifyMessage(t *testing.T) {
	requestTypes := []wsmessage.MessageType{
		wsmessage.RequestPermissionRequestType,
		wsmessage.WriteTextFileRequestType,
		wsmessage.ReadTextFileRequestType,
		wsmessage.CreateTerminalRequestType,
		wsmessage.TerminalOutputRequestType,
		wsmessage.ReleaseTerminalRequestType,
		wsmessage.WaitForTerminalExitRequestType,
		wsmessage.KillTerminalRequestType,
		wsmessage.SetSessionModeRequestType,
		wsmessage.UnstableForkSessionRequestType,
		wsmessage.ListSessionsRequestType,
		wsmessage.ResumeSessionRequestType,
		wsmessage.SetSessionConfigOptionRequestType,
		wsmessage.LogoutRequestType,
		wsmessage.UnstableCloseNesRequestType,
		wsmessage.UnstableStartNesRequestType,
		wsmessage.UnstableSuggestNesRequestType,
		wsmessage.UnstableDisableProviderRequestType,
		wsmessage.UnstableListProvidersRequestType,
		wsmessage.UnstableSetProviderRequestType,
		wsmessage.UnstableDeleteSessionRequestType,
		wsmessage.CloseSessionRequestType,
		wsmessage.InitializeRequestType,
		wsmessage.NewSessionRequestType,
		wsmessage.AuthenticateRequestType,
		wsmessage.LoadSessionRequestType,
		wsmessage.PromptRequestType,
	}

	responseTypes := []wsmessage.MessageType{
		wsmessage.RequestPermissionResponseType,
		wsmessage.WriteTextFileResponseType,
		wsmessage.ReadTextFileResponseType,
		wsmessage.CreateTerminalResponseType,
		wsmessage.TerminalOutputResponseType,
		wsmessage.ReleaseTerminalResponseType,
		wsmessage.WaitForTerminalExitResponseType,
		wsmessage.KillTerminalResponseType,
		wsmessage.SetSessionModeResponseType,
		wsmessage.UnstableForkSessionResponseType,
		wsmessage.ListSessionsResponseType,
		wsmessage.ResumeSessionResponseType,
		wsmessage.SetSessionConfigOptionResponseType,
		wsmessage.LogoutResponseType,
		wsmessage.UnstableCloseNesResponseType,
		wsmessage.UnstableStartNesResponseType,
		wsmessage.UnstableSuggestNesResponseType,
		wsmessage.UnstableDisableProviderResponseType,
		wsmessage.UnstableListProvidersResponseType,
		wsmessage.UnstableSetProviderResponseType,
		wsmessage.UnstableDeleteSessionResponseType,
		wsmessage.CloseSessionResponseType,
		wsmessage.InitializeResponseType,
		wsmessage.NewSessionResponseType,
		wsmessage.AuthenticateResponseType,
		wsmessage.LoadSessionResponseType,
		wsmessage.PromptResponseType,
	}

	notificationTypes := []wsmessage.MessageType{
		wsmessage.SessionNotificationType,
		wsmessage.UnstableAcceptNesNotificationType,
		wsmessage.UnstableRejectNesNotificationType,
		wsmessage.UnstableDidChangeDocumentNotificationType,
		wsmessage.UnstableDidCloseDocumentNotificationType,
		wsmessage.UnstableDidFocusDocumentNotificationType,
		wsmessage.UnstableDidOpenDocumentNotificationType,
		wsmessage.UnstableDidSaveDocumentNotificationType,
		wsmessage.CancelNotificationType,
	}

	for _, mt := range requestTypes {
		if got := classifyMessage(mt); got != kindRequest {
			t.Errorf("classifyMessage(%d) = %v, want kindRequest", mt, got)
		}
	}
	for _, mt := range responseTypes {
		if got := classifyMessage(mt); got != kindResponse {
			t.Errorf("classifyMessage(%d) = %v, want kindResponse", mt, got)
		}
	}
	for _, mt := range notificationTypes {
		if got := classifyMessage(mt); got != kindNotification {
			t.Errorf("classifyMessage(%d) = %v, want kindNotification", mt, got)
		}
	}
}
