//go:build wasip1

package main

import (
	"context"
	"fmt"

	"github.com/SethCurry/abyss/pkg/protobyss"
)

// This has to exist to compile, per go-plugin.  It can be empty, however.
func main() {}

func init() {
	// Register our plugin in the init
	protobyss.RegisterACPPlugin(NewACPLoggerPlugin())
}

func NewACPLoggerPlugin() *ACPLoggerPlugin {
	return &ACPLoggerPlugin{}
}

// Compile-time type check
var _ protobyss.ACPPlugin = (*ACPLoggerPlugin)(nil)

type ACPLoggerPlugin struct{}

// We get all messages.  This is easier than ACPPluginRouter here since we need to log all messages and don't care what type they are.
func (p *ACPLoggerPlugin) HandleMessage(ctx context.Context, message *protobyss.ACPContainer) (*protobyss.ACPContainerList, error) {
	fmt.Printf("%s\n", string(message.Content))
	return &protobyss.ACPContainerList{
		Containers: []*protobyss.ACPContainer{message},
	}, nil
}
