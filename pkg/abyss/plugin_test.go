package abyss

import (
	"context"
	"testing"

	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/stretchr/testify/assert"
)

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
