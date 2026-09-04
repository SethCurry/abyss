package wsmessage_test

import (
	"testing"

	"github.com/SethCurry/abyss/internal/websockets/wsmessage"
	"github.com/coder/acp-go-sdk"
)

func TestGetMessageTypeByID(t *testing.T) {
	got, err := wsmessage.GetMessageTypeByID(int32(wsmessage.CancelNotificationType))
	if err != nil {
		t.Fatalf("GetMessageTypeByID returned error: %v", err)
	}
	if got.TypeID != wsmessage.CancelNotificationType {
		t.Errorf("TypeID = %v, want %v", got.TypeID, wsmessage.CancelNotificationType)
	}
}

func TestGetMessageTypeByIDUnknown(t *testing.T) {
	if _, err := wsmessage.GetMessageTypeByID(9999); err == nil {
		t.Error("GetMessageTypeByID with unknown ID returned nil error, want non-nil")
	}
}

func TestGetMessageTypeByType(t *testing.T) {
	got, err := wsmessage.GetMessageTypeByType(acp.CancelNotification{})
	if err != nil {
		t.Fatalf("GetMessageTypeByType returned error: %v", err)
	}
	if got.TypeID != wsmessage.CancelNotificationType {
		t.Errorf("TypeID = %v, want %v", got.TypeID, wsmessage.CancelNotificationType)
	}
}

func TestGetMessageTypeByTypeUnknown(t *testing.T) {
	if _, err := wsmessage.GetMessageTypeByType(struct{}{}); err == nil {
		t.Error("GetMessageTypeByType with unknown type returned nil error, want non-nil")
	}
}
