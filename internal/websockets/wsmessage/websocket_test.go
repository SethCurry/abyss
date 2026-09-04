package wsmessage_test

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
