package wsmultiplex

import (
	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
)

type messageKind int

const (
	kindRequest messageKind = iota
	kindResponse
	kindNotification
)

// classifyMessage returns the kind of a MessageType.
func classifyMessage(mt wsmessage.MessageType) messageKind {
	switch mt {
	case wsmessage.RequestPermissionRequestType,
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
		wsmessage.PromptRequestType:
		return kindRequest
	case wsmessage.RequestPermissionResponseType,
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
		wsmessage.PromptResponseType:
		return kindResponse
	default:
		return kindNotification
	}
}
