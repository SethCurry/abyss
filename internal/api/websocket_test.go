package api

import (
	"testing"

	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/coder/acp-go-sdk"
)

func TestUnmarshalMessagePopulatesFields(t *testing.T) {
	got, err := wsmessage.UnmarshalMessage(wsmessage.CancelNotificationType, []byte(`{"sessionId":"session-123"}`))
	if err != nil {
		t.Fatalf("unmarshalMessage returned error: %v", err)
	}

	msg, ok := got.(*acp.CancelNotification)
	if !ok {
		t.Fatalf("unmarshalMessage returned %T, want *acp.CancelNotification", got)
	}
	if msg.SessionId != "session-123" {
		t.Errorf("SessionId = %q, want %q", msg.SessionId, "session-123")
	}
}

func TestUnmarshalMessageInvalidJSON(t *testing.T) {
	if _, err := wsmessage.UnmarshalMessage(wsmessage.CancelNotificationType, []byte("not json")); err == nil {
		t.Error("unmarshalMessage with invalid JSON returned nil error, want non-nil")
	}
}

func TestUnmarshalMessageNotExist(t *testing.T) {
	if _, err := wsmessage.UnmarshalMessage(wsmessage.MessageTypeNotExist, []byte("{}")); err == nil {
		t.Error("unmarshalMessage(MessageTypeNotExist) returned nil error, want non-nil")
	}
}

func TestUnmarshalMessageUnknownType(t *testing.T) {
	if _, err := wsmessage.UnmarshalMessage(wsmessage.MessageType(9999), []byte("{}")); err == nil {
		t.Error("unmarshalMessage with unknown type returned nil error, want non-nil")
	}
}

// messageTypeForType dispatches to the generic messageTypeFor with the
// concrete type of msg.
func messageTypeForType(msg any) wsmessage.MessageType {
	switch v := msg.(type) {
	case acp.RequestPermissionRequest:
		return messageTypeFor(v)
	case acp.RequestPermissionResponse:
		return messageTypeFor(v)
	case acp.WriteTextFileRequest:
		return messageTypeFor(v)
	case acp.WriteTextFileResponse:
		return messageTypeFor(v)
	case acp.ReadTextFileRequest:
		return messageTypeFor(v)
	case acp.ReadTextFileResponse:
		return messageTypeFor(v)
	case acp.CreateTerminalRequest:
		return messageTypeFor(v)
	case acp.CreateTerminalResponse:
		return messageTypeFor(v)
	case acp.TerminalOutputRequest:
		return messageTypeFor(v)
	case acp.TerminalOutputResponse:
		return messageTypeFor(v)
	case acp.ReleaseTerminalRequest:
		return messageTypeFor(v)
	case acp.ReleaseTerminalResponse:
		return messageTypeFor(v)
	case acp.WaitForTerminalExitRequest:
		return messageTypeFor(v)
	case acp.WaitForTerminalExitResponse:
		return messageTypeFor(v)
	case acp.KillTerminalRequest:
		return messageTypeFor(v)
	case acp.KillTerminalResponse:
		return messageTypeFor(v)
	case acp.SessionNotification:
		return messageTypeFor(v)
	case acp.SetSessionModeRequest:
		return messageTypeFor(v)
	case acp.SetSessionModeResponse:
		return messageTypeFor(v)
	case acp.UnstableForkSessionRequest:
		return messageTypeFor(v)
	case acp.UnstableForkSessionResponse:
		return messageTypeFor(v)
	case acp.ListSessionsRequest:
		return messageTypeFor(v)
	case acp.ListSessionsResponse:
		return messageTypeFor(v)
	case acp.ResumeSessionRequest:
		return messageTypeFor(v)
	case acp.ResumeSessionResponse:
		return messageTypeFor(v)
	case acp.SetSessionConfigOptionRequest:
		return messageTypeFor(v)
	case acp.SetSessionConfigOptionResponse:
		return messageTypeFor(v)
	case acp.LogoutRequest:
		return messageTypeFor(v)
	case acp.LogoutResponse:
		return messageTypeFor(v)
	case acp.UnstableCloseNesRequest:
		return messageTypeFor(v)
	case acp.UnstableCloseNesResponse:
		return messageTypeFor(v)
	case acp.UnstableStartNesRequest:
		return messageTypeFor(v)
	case acp.UnstableStartNesResponse:
		return messageTypeFor(v)
	case acp.UnstableSuggestNesRequest:
		return messageTypeFor(v)
	case acp.UnstableSuggestNesResponse:
		return messageTypeFor(v)
	case acp.UnstableAcceptNesNotification:
		return messageTypeFor(v)
	case acp.UnstableRejectNesNotification:
		return messageTypeFor(v)
	case acp.UnstableDidChangeDocumentNotification:
		return messageTypeFor(v)
	case acp.UnstableDidCloseDocumentNotification:
		return messageTypeFor(v)
	case acp.UnstableDidFocusDocumentNotification:
		return messageTypeFor(v)
	case acp.UnstableDidOpenDocumentNotification:
		return messageTypeFor(v)
	case acp.UnstableDidSaveDocumentNotification:
		return messageTypeFor(v)
	case acp.UnstableDisableProviderRequest:
		return messageTypeFor(v)
	case acp.UnstableDisableProviderResponse:
		return messageTypeFor(v)
	case acp.UnstableListProvidersRequest:
		return messageTypeFor(v)
	case acp.UnstableListProvidersResponse:
		return messageTypeFor(v)
	case acp.UnstableSetProviderRequest:
		return messageTypeFor(v)
	case acp.UnstableSetProviderResponse:
		return messageTypeFor(v)
	case acp.UnstableDeleteSessionRequest:
		return messageTypeFor(v)
	case acp.UnstableDeleteSessionResponse:
		return messageTypeFor(v)
	case acp.CloseSessionRequest:
		return messageTypeFor(v)
	case acp.CloseSessionResponse:
		return messageTypeFor(v)
	case acp.InitializeRequest:
		return messageTypeFor(v)
	case acp.InitializeResponse:
		return messageTypeFor(v)
	case acp.NewSessionRequest:
		return messageTypeFor(v)
	case acp.NewSessionResponse:
		return messageTypeFor(v)
	case acp.AuthenticateRequest:
		return messageTypeFor(v)
	case acp.AuthenticateResponse:
		return messageTypeFor(v)
	case acp.LoadSessionRequest:
		return messageTypeFor(v)
	case acp.LoadSessionResponse:
		return messageTypeFor(v)
	case acp.PromptRequest:
		return messageTypeFor(v)
	case acp.PromptResponse:
		return messageTypeFor(v)
	case acp.CancelNotification:
		return messageTypeFor(v)
	}
	return wsmessage.MessageTypeNotExist
}
