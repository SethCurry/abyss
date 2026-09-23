package abyss

import (
	"context"
	"testing"

	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// promptRequestPlugin only implements OnPromptRequest, allowing the test to
// verify that NewACPPluginRouter wires the handler up and HandleMessage
// dispatches to it.
type promptRequestPlugin struct {
	called   bool
	received acp.PromptRequest
}

func (p *promptRequestPlugin) OnPromptRequest(
	req acp.PromptRequest,
) ([]*protobyss.ACPContainer, error) {
	p.called = true
	p.received = req
	return nil, nil
}

func Test_ACPPluginRouter_HandleMessage_DispatchesToHandler(t *testing.T) {
	plugin := &promptRequestPlugin{}
	router := NewACPPluginRouter(plugin)

	req := acp.PromptRequest{
		SessionId: "test-session",
	}

	msg, err := ACPContainer(req)
	require.NoError(t, err)

	_, err = router.HandleMessage(context.Background(), msg)
	require.NoError(t, err)

	assert.True(t, plugin.called, "OnPromptRequest handler was not invoked")
	assert.Equal(t, acp.SessionId("test-session"), plugin.received.SessionId)
}

func Test_ACPPluginRouter_Completeness(t *testing.T) {
	plugin := &ACPPluginRouter{}
	for _, v := range AllMessageTypes {
		_, err := plugin.HandleMessage(context.Background(), &protobyss.ACPContainer{
			TypeId:  int32(v.TypeID()),
			Content: []byte("{}"),
		})
		assert.Nil(t, err, "expected no error, got %v", err)
	}
}
