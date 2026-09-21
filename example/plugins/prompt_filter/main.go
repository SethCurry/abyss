//go:build wasip1

package main

import (
	"regexp"
	"strings"

	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

type PromptFilter struct {
	bannedRegexes []*regexp.Regexp
}

// We only need to implement OnPromptRequest since we don't care about the other types.
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

			return abyss.ACPContainers(msg)
		}
	}

	return abyss.ACPContainers(
		req,
	)
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
