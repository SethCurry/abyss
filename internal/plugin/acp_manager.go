package plugin

import (
	"context"

	"github.com/SethCurry/abyss/internal/timber"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/rs/zerolog"
)

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

type ACPManager struct {
	loader  *protobyss.ACPPluginPlugin
	plugins []protobyss.ACPPlugin
	logger  zerolog.Logger
}

func (a *ACPManager) Load(ctx context.Context, path string) error {
	a.logger.Info().Str("path", path).Msg("loading ACP plugin")
	plugin, err := a.loader.Load(ctx, path)
	if err != nil {
		return err
	}
	a.plugins = append(a.plugins, plugin)
	return nil
}

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
