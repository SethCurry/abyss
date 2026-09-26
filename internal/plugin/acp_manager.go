// Package plugin provides access control policy management for plugins.
package plugin

import (
	"context"

	"github.com/SethCurry/abyss/internal/timber"
	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/rs/zerolog"
)

func logMessage(event *zerolog.Event, msg *protobyss.LogMessage) (*emptypb.Empty, error) {
	for k, v := range msg.GetFields() {
		event.Str(k, v)
	}

	event.Msg(msg.GetMessage())

	return nil, nil
}

// NewLogging creates a new protobyss.Logging.
func NewLogging(logger zerolog.Logger, plugPath string) *Logging {
	return &Logging{
		logger: logger.With().Str("plugin_path", plugPath).Logger(),
	}
}

// Logging implements WASM-host side logging so that plugin messages are
// visible in the logs.
type Logging struct {
	logger zerolog.Logger
}

// Debug logs a message at DEBUG level.
func (l *Logging) Debug(ctx context.Context, msg *protobyss.LogMessage) (*emptypb.Empty, error) {
	return logMessage(l.logger.Debug(), msg)
}

// Info logs a message at INFO level.
func (l *Logging) Info(ctx context.Context, msg *protobyss.LogMessage) (*emptypb.Empty, error) {
	return logMessage(l.logger.Info(), msg)
}

// Warn logs a message WARN level.
func (l *Logging) Warn(ctx context.Context, msg *protobyss.LogMessage) (*emptypb.Empty, error) {
	return logMessage(l.logger.Warn(), msg)
}

// Error logs a message at ERROR level.
func (l *Logging) Error(ctx context.Context, msg *protobyss.LogMessage) (*emptypb.Empty, error) {
	return logMessage(l.logger.Error(), msg)
}

// NewACPManager creates and initializes a new ACPManager instance.
func NewACPManager(ctx context.Context) (*ACPManager, error) {
	loader, err := protobyss.NewACPPluginPlugin(ctx)
	if err != nil {
		return nil, err
	}
	return &ACPManager{
		loader:  loader,
		logger:  timber.ComponentLogger("plugin.ACPManager"),
		plugins: []protobyss.ACPPlugin{},
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
	plugin, err := a.loader.Load(ctx, path, NewLogging(a.logger, path))
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
	a.logger.Info().Msg("plugin handling message")
	allMessages := []*protobyss.ACPContainer{req}
	for i, v := range a.plugins {
		a.logger.Info().Int("plugin_index", i).Msg("executing plugin")
		var newMsgs []*protobyss.ACPContainer
		for _, msg := range allMessages {
			gotMsgs, err := v.HandleMessage(ctx, msg)
			if err != nil {
				a.logger.Error().Err(err).Msg("plugin failed")
				return allMessages, err
			}

			for _, newMsg := range gotMsgs.GetContainers() {
				msgType, err := abyss.GetMessageTypeByID(newMsg.GetTypeId())
				if err != nil {
					return allMessages, err
				}

				if msgType.IsResponse() && newMsg.GetResponseFor() == "" {
					switch {
					case msg.GetResponseFor() != "":
						newMsg.ResponseFor = msg.GetResponseFor()
					case req.GetResponseFor() != "":
						newMsg.ResponseFor = req.GetResponseFor()
					default:
						newMsg.ResponseFor = req.GetMessageId()
					}
				}
				a.logger.Info().Str("content", string(newMsg.GetContent())).Msg("got plugin results")
			}
			newMsgs = append(newMsgs, gotMsgs.GetContainers()...)
		}
		allMessages = newMsgs
	}

	allMessagesLen := len(allMessages)
	if allMessagesLen == 1 {
		allMessages[0].MessageId = req.GetMessageId()
	}
	return allMessages, nil
}
