// Package plugin provides access control policy management for plugins.
package plugin

import (
	"context"

	"github.com/SethCurry/abyss/internal/timber"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/rs/zerolog"
)

// NewACPManager creates and initializes a new ACPManager instance.
func NewACPManager(ctx context.Context) (*ACPManager, error) {
	loader, err := protobyss.NewACPPluginPlugin(ctx)
	if err != nil {
		return nil, err
	}
	return &ACPManager{
		loader: loader,
		logger: timber.ComponentLogger("plugin.ACPManager"),
	}, nil
}

// ACPManager manages access control policy plugins and their lifecycle.
type ACPManager struct {
	loader  *protobyss.ACPPluginPlugin
	plugins []protobyss.ACPPlugin
	logger  zerolog.Logger
}

// Load loads an ACP plugin from the given path.
func (a *ACPManager) Load(ctx context.Context, path string) error {
	a.logger.Info().Str("path", path).Msg("loading ACP plugin")
	plugin, err := a.loader.Load(ctx, path)
	if err != nil {
		return err
	}
	a.plugins = append(a.plugins, plugin)
	return nil
}

// HandleMessage passes a message through all loaded ACP plugins and returns
// the resulting messages.
func (a *ACPManager) HandleMessage(
	ctx context.Context, req *protobyss.ACPContainer,
) ([]*protobyss.ACPContainer, error) {
	allMessages := make([]*protobyss.ACPContainer, 1)
	allMessages[0] = req
	for _, v := range a.plugins {
		var newMsgs []*protobyss.ACPContainer
		for _, msg := range allMessages {
			gotMsgs, err := v.HandleMessage(ctx, msg)
			if err != nil {
				return nil, err
			}
			newMsgs = append(newMsgs, gotMsgs.GetContainers()...)
		}
		allMessages = newMsgs
	}
	return allMessages, nil
}
