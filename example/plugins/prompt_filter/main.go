//go:build wasip1

package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

type PromptFilter struct {
	bannedRegexes []*regexp.Regexp
}

func (p *PromptFilter) OnPromptRequest(req acp.PromptRequest) ([]*protobyss.ACPContainer, error) {
	allStringContents := strings.Builder{}

	for _, v := range req.Prompt {
		if v.Text != nil {
			allStringContents.WriteString(v.Text.Text)
		}
	}

	contents := []byte(allStringContents.String())

	for _, v := range p.bannedRegexes {
		if v.Match(contents) {
			msg := acp.SessionNotification{
				SessionId: req.SessionId,
				Update: acp.SessionUpdate{
					AgentMessageChunk: &acp.SessionUpdateAgentMessageChunk{
						Content: acp.TextBlock("Nuh uh, not under my roof!"),
					},
				},
			}

			marshalled, err := json.Marshal(msg)
			if err != nil {
				return nil, fmt.Errorf("failed to generate JSON for SessionNotification: %w", err)
			}

			return []*protobyss.ACPContainer{
				&protobyss.ACPContainer{
					TypeId:  int32(abyss.SessionNotificationType),
					Content: marshalled,
				},
			}, nil
		}
	}

	marshalledMsg, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate default JSON response: %w", err)
	}

	return []*protobyss.ACPContainer{
		&protobyss.ACPContainer{
			TypeId:  int32(abyss.SessionNotificationType),
			Content: marshalledMsg,
		},
	}, nil
}

func main() {}

func init() {
	bannedPromptStrings := []string{".*A_SECRET.*"}
	bannedPromptRegexes := make([]*regexp.Regexp, len(bannedPromptStrings))
	for k, v := range bannedPromptStrings {
		bannedPromptRegexes[k] = regexp.MustCompile(v)
	}

	myPlugin := &PromptFilter{bannedRegexes: bannedPromptRegexes}

	plugin := abyss.NewACPPluginRouter(myPlugin)
	protobyss.RegisterACPPlugin(plugin)
}
