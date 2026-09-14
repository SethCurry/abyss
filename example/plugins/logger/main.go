//go:build wasip1

package main

import (
	"context"
	"fmt"

	"github.com/SethCurry/abyss/pkg/protobyss"
)

func main() {}

func init() {
	protobyss.RegisterACPPlugin(NewACPLoggerPlugin())
}

func NewACPLoggerPlugin() *ACPLoggerPlugin {
	return &ACPLoggerPlugin{}
}

var _ protobyss.ACPPlugin = (*ACPLoggerPlugin)(nil)

type ACPLoggerPlugin struct{}

func (p *ACPLoggerPlugin) HandleMessage(ctx context.Context, message *protobyss.ACPContainer) (*protobyss.ACPContainerList, error) {
	fmt.Printf("%s\n", string(message.Content))
	return &protobyss.ACPContainerList{
		Containers: []*protobyss.ACPContainer{message},
	}, nil
}
