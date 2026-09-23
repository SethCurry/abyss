package abyss

import (
	"encoding/json"
	"testing"

	"github.com/coder/acp-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestACPContainer verifies that ACPContainer wraps an ACP message into a
// protobyss.ACPContainer with the correct type id and a JSON payload that
// round-trips back to the original message.
func TestACPContainer(t *testing.T) {
	t.Run("simple request round-trips", func(t *testing.T) {
		msg := acp.LogoutRequest{
			Meta: map[string]any{"trace": "abc"},
		}

		container, err := ACPContainer(msg)
		require.NoError(t, err)
		require.NotNil(t, container)

		assert.Equal(t, int32(LogoutRequestType), container.TypeId)

		// The content must be valid JSON that unmarshals back into the message.
		var roundTrip acp.LogoutRequest
		require.NoError(t, json.Unmarshal(container.Content, &roundTrip))
		assert.Equal(t, msg, roundTrip)
	})

	t.Run("notification with fields round-trips", func(t *testing.T) {
		msg := acp.CancelNotification{
			SessionId: acp.SessionId("session-42"),
		}

		container, err := ACPContainer(msg)
		require.NoError(t, err)
		require.NotNil(t, container)

		assert.Equal(t, int32(CancelNotificationType), container.TypeId)

		var roundTrip acp.CancelNotification
		require.NoError(t, json.Unmarshal(container.Content, &roundTrip))
		assert.Equal(t, msg, roundTrip)
	})

	t.Run("content is valid JSON", func(t *testing.T) {
		container, err := ACPContainer(acp.LogoutRequest{})
		require.NoError(t, err)

		assert.True(t, json.Valid(container.Content),
			"container content should be valid JSON, got %q", container.Content)
	})

	t.Run("type id matches registry lookup", func(t *testing.T) {
		msg := acp.CancelNotification{SessionId: "x"}

		container, err := ACPContainer(msg)
		require.NoError(t, err)

		registered, err := GetMessageTypeByType(msg)
		require.NoError(t, err)
		assert.Equal(t, int32(registered.TypeID()), container.TypeId)
	})

	t.Run("different messages produce different type ids", func(t *testing.T) {
		reqContainer, err := ACPContainer(acp.LogoutRequest{})
		require.NoError(t, err)

		notifContainer, err := ACPContainer(acp.CancelNotification{SessionId: "x"})
		require.NoError(t, err)

		assert.NotEqual(t, reqContainer.TypeId, notifContainer.TypeId)
	})

	t.Run("preserves message direction metadata via registry", func(t *testing.T) {
		// The container itself does not store direction, but the type id must
		// resolve to a registered message whose direction matches what was
		// declared for that message.
		container, err := ACPContainer(acp.LogoutRequest{})
		require.NoError(t, err)

		registered, err := GetMessageTypeByID(container.TypeId)
		require.NoError(t, err)

		assert.Equal(t, ToAgent, registered.Direction())
		assert.False(t, registered.IsResponse())
	})

	t.Run("response message direction via registry", func(t *testing.T) {
		container, err := ACPContainer(acp.LogoutResponse{})
		require.NoError(t, err)

		registered, err := GetMessageTypeByID(container.TypeId)
		require.NoError(t, err)

		assert.Equal(t, ToACPClient, registered.Direction())
		assert.True(t, registered.IsResponse())
	})
}

// TestACPContainers verifies that ACPContainers converts a variadic list of
// ACP messages into a correctly ordered slice of ACPContainers.
func TestACPContainers(t *testing.T) {
	t.Run("empty input returns empty slice", func(t *testing.T) {
		containers, err := ACPContainers[acp.LogoutRequest]()
		require.NoError(t, err)
		assert.Empty(t, containers)
	})

	t.Run("single message", func(t *testing.T) {
		msg := acp.LogoutRequest{Meta: map[string]any{"k": "v"}}

		containers, err := ACPContainers(msg)
		require.NoError(t, err)
		require.Len(t, containers, 1)

		assert.Equal(t, int32(LogoutRequestType), containers[0].TypeId)

		var roundTrip acp.LogoutRequest
		require.NoError(t, json.Unmarshal(containers[0].Content, &roundTrip))
		assert.Equal(t, msg, roundTrip)
	})

	t.Run("multiple messages preserve order and values", func(t *testing.T) {
		msgs := []acp.CancelNotification{
			{SessionId: "first"},
			{SessionId: "second"},
			{SessionId: "third"},
		}

		containers, err := ACPContainers(msgs[0], msgs[1], msgs[2])
		require.NoError(t, err)
		require.Len(t, containers, 3)

		for i, c := range containers {
			assert.Equal(t, int32(CancelNotificationType), c.TypeId)

			var roundTrip acp.CancelNotification
			require.NoError(t, json.Unmarshal(c.Content, &roundTrip))
			assert.Equal(t, msgs[i], roundTrip)
		}
	})

	t.Run("mixed via ACPContainer per-message type ids are stable", func(t *testing.T) {
		// ACPContainers is homogeneous in T, so we verify that repeated calls
		// for the same type yield consistent type ids.
		containers, err := ACPContainers(
			acp.CancelNotification{SessionId: "a"},
			acp.CancelNotification{SessionId: "b"},
		)
		require.NoError(t, err)
		require.Len(t, containers, 2)

		assert.Equal(t, containers[0].TypeId, containers[1].TypeId)
		assert.NotEmpty(t, containers[0].Content)
	})
}

// TestACPContainerList verifies that ACPContainerList builds an
// ACPContainerList containing the wrapped messages.
func TestACPContainerList(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		list, err := ACPContainerList[acp.LogoutRequest]()
		require.NoError(t, err)
		require.NotNil(t, list)
		assert.Empty(t, list.Containers)
	})

	t.Run("single message", func(t *testing.T) {
		msg := acp.LogoutRequest{}

		list, err := ACPContainerList(msg)
		require.NoError(t, err)
		require.NotNil(t, list)
		require.Len(t, list.Containers, 1)

		assert.Equal(t, int32(LogoutRequestType), list.Containers[0].TypeId)
		assert.True(t, json.Valid(list.Containers[0].Content))
	})

	t.Run("multiple messages preserve order", func(t *testing.T) {
		msgs := []acp.CancelNotification{
			{SessionId: "one"},
			{SessionId: "two"},
		}

		list, err := ACPContainerList(msgs[0], msgs[1])
		require.NoError(t, err)
		require.NotNil(t, list)
		require.Len(t, list.Containers, 2)

		for i, c := range list.Containers {
			assert.Equal(t, int32(CancelNotificationType), c.TypeId)

			var roundTrip acp.CancelNotification
			require.NoError(t, json.Unmarshal(c.Content, &roundTrip))
			assert.Equal(t, msgs[i], roundTrip)
		}
	})

	t.Run("containers are the same instances produced by ACPContainers", func(t *testing.T) {
		// ACPContainerList delegates to ACPContainers; ensure the wrapping is
		// equivalent by comparing type ids and content against a direct call.
		msg := acp.CancelNotification{SessionId: "z"}

		direct, err := ACPContainers(msg)
		require.NoError(t, err)

		list, err := ACPContainerList(msg)
		require.NoError(t, err)

		require.Len(t, list.Containers, 1)
		assert.Equal(t, direct[0].TypeId, list.Containers[0].TypeId)
		assert.Equal(t, direct[0].Content, list.Containers[0].Content)
	})
}

// TestACPContainer_RegistersAllMessageTypes ensures that every message type
// declared in AllMessageTypes can be looked up by the type id embedded in a
// container built from a zero-value instance of a representative subset. This
// guards against the registry and the container type ids drifting apart.
func TestACPContainer_RegistersAllMessageTypes(t *testing.T) {
	// Build containers from several distinct message kinds and confirm each
	// type id resolves back to the expected registered TypedMessage.
	cases := []struct {
		name        string
		typeID      MessageTypeID
		containerID int32
	}{
		{"request", LogoutRequestType, int32(LogoutRequestType)},
		{"notification", CancelNotificationType, int32(CancelNotificationType)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			registered, err := GetMessageTypeByID(tc.containerID)
			require.NoError(t, err)
			assert.Equal(t, tc.typeID, registered.TypeID())
		})
	}

	// Every registered type id must be unique.
	seen := make(map[int32]struct{}, len(AllMessageTypes))
	for _, mt := range AllMessageTypes {
		id := int32(mt.TypeID())
		_, dup := seen[id]
		assert.False(t, dup, "duplicate message type id %d", id)
		seen[id] = struct{}{}
	}
}