package plugin

import (
	"context"

	"github.com/SethCurry/abyss/pkg/protobyss"
)

func NewACPManager(ctx context.Context) (*ACPManager, error) {
	loader, err := protobyss.NewACPPluginPlugin(ctx)
	if err != nil {
		return nil, err
	}
	return &ACPManager{
		loader: loader,
	}, nil
}

type ACPManager struct {
	loader  *protobyss.ACPPluginPlugin
	plugins []protobyss.ACPPlugin
}

func (a *ACPManager) Load(ctx context.Context, path string) error {
	plugin, err := a.loader.Load(ctx, path)
	if err != nil {
		return err
	}
	a.plugins = append(a.plugins, plugin)
	return nil
}

func (a *ACPManager) Handle(ctx context.Context, req *protobyss.ACPContainer) ([]*protobyss.ACPContainer, error) {
	allMessages := make([]*protobyss.ACPContainer, 1)
	allMessages[0] = req
	for _, v := range a.plugins {
		var newMsgs []*protobyss.ACPContainer
		for _, msg := range allMessages {
			gotMsgs, err := v.HandleMessage(ctx, msg)
			if err != nil {
				return nil, err
			}
			newMsgs = append(newMsgs, gotMsgs.Containers...)
		}
		allMessages = newMsgs
	}
	return allMessages, nil
}
