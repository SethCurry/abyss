//go:build wasip1

package main

import (
	"regexp"

	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

type PromptFilter struct {
	bannedRegexes []*regexp.Regexp
	logging       protobyss.Logging
	initDone      bool
}

func (p *PromptFilter) OnSessionNotification(notif acp.SessionNotification) ([]*protobyss.ACPContainer, error) {
	//if notif.Update.AgentMessageChunk != nil && notif.Update.AgentMessageChunk.Content.Text != nil {
	//	for _, v := range p.bannedRegexes {
	//		if v.MatchString(notif.Update.AgentMessageChunk.Content.Text.Text) {
	msg := acp.SessionNotification{
		SessionId: notif.SessionId,
		Update: acp.SessionUpdate{
			AgentMessageChunk: &acp.SessionUpdateAgentMessageChunk{
				Content: acp.TextBlock("content blocked by plugin"),
			},
		},
	}

	return abyss.ACPContainers(msg)
	//		}
	//}
	//}

	//return abyss.ACPContainers(notif)
}

func main() {
}

func init() {
	bannedPromptStrings := []string{".*SECRET.*"}
	bannedPromptRegexes := make([]*regexp.Regexp, len(bannedPromptStrings))
	for k, v := range bannedPromptStrings {
		bannedPromptRegexes[k] = regexp.MustCompile(v)
	}

	myPlugin := &PromptFilter{
		bannedRegexes: bannedPromptRegexes,
	}

	plugin := abyss.NewACPPluginRouter(myPlugin)
	protobyss.RegisterACPPlugin(plugin)
}
