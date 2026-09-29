//go:build wasip1

package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

// SizedStringBuffer retains a bounded history of strings, dropping the
// oldest entries once the fixed capacity is reached.
type SizedStringBuffer struct {
	maxSize int
	buf     []string
}

// NewSizedStringBuffer creates a buffer that holds at most maxSize strings.
func NewSizedStringBuffer(maxSize int) *SizedStringBuffer {
	if maxSize < 0 {
		maxSize = 0
	}

	return &SizedStringBuffer{
		maxSize: maxSize,
		buf:     make([]string, 0, maxSize),
	}
}

func (b *SizedStringBuffer) Reset() {
	b.buf = b.buf[:0]
}

// Append adds s to the buffer, truncating the oldest content if the
// buffer is already at capacity.
func (b *SizedStringBuffer) Append(s string) {
	if b.maxSize == 0 {
		return
	}

	b.buf = append(b.buf, s)
	if len(b.buf) > b.maxSize {
		b.buf = b.buf[len(b.buf)-b.maxSize:]
	}
}

// String returns the concatenation of the buffer's current contents.
func (b *SizedStringBuffer) String() string {
	return strings.Join(b.buf, "")
}

type SecretsFilter struct {
	detectors        []SecretDetector
	promptRequestBuf *SizedStringBuffer
	agentMessageBuf  *SizedStringBuffer
	agentThoughtBuf  *SizedStringBuffer
	logging          protobyss.Logging
	initDone         bool
	currentSession   acp.SessionId
	isCancelled      bool
}

func (p *SecretsFilter) Initialize(ctx context.Context, request *protobyss.ACPPluginInitializeRequest) (*protobyss.ACPPluginInitializeResponse, error) {
	p.logging = protobyss.NewLogging()
	return &protobyss.ACPPluginInitializeResponse{
		Name: "secrets_filter",
	}, nil
}

func (p *SecretsFilter) getCancelMessages(matchedDetector SecretDetector) ([]*protobyss.ACPContainer, error) {
	// cancel session to agent
	// notification to user
	if p.isCancelled {
		return []*protobyss.ACPContainer{}, nil
	}
	p.isCancelled = true

	cancelMsg := acp.CancelNotification{
		SessionId: p.currentSession,
	}
	cancelProto, err := abyss.ACPContainer(cancelMsg)
	if err != nil {
		return nil, err
	}

	userNotification := acp.SessionNotification{
		SessionId: p.currentSession,
		Update: acp.SessionUpdate{
			AgentMessageChunk: &acp.SessionUpdateAgentMessageChunk{
				Content: acp.TextBlock("This session has been cancelled due to a security filter match.\n" +
					"Matched detector: " + matchedDetector.Name() + "\n"),
			},
		},
	}
	userNotificationProto, err := abyss.ACPContainer(userNotification)
	if err != nil {
		return nil, err
	}

	return []*protobyss.ACPContainer{userNotificationProto, cancelProto}, nil
}

/*
{
   "sessionId": "01a0ee12-1023-7322-a5e2-15a4736f6228",
   "update": {
     "_meta": {
       "terminal_output": {
         "data": "\treturn b.maxLength\n}\n\nfunc init() {\n\tdetectors := []SecretDetector{\n\t\t&BannedRegex{\n\t\t\t// test string sk-or-v1-abcd-efg-hiklkmno-p-qrs-tuv-wxyz\n\t\t\tname:      \"openrouter\",\n\t\t\tRegex:     regexp.MustCompile(\".*sk-or-v1-[A-Za-z0-9_-]{32,128}.*\"),\n\t\t\tmaxLength: 140,\n\t\t},\n\t}\n\n\tmaxLength := 0\n\tfor _, d := range detectors {\n\t\tif l := d.MaxLength(); l > maxLength {\n\t\t\tmaxLength = l\n\t\t}\n\t}\n\n\tmyPlugin := &SecretsFilter{\n\t\tdetectors:        detectors,\n\t\tpromptRequestBuf: NewSizedStringBuffer(maxLength),\n\t\tagentMessageBuf:  NewSizedStringBuffer(maxLength),\n\t\tagentThoughtBuf:  NewSizedStringBuffer(maxLength),\n\t}\n\n\tplugin := abyss.NewACPPluginRouter(myPlugin)\n\tprotobyss.RegisterACPPlugin(plugin)\n}\n",
         "terminal_id": "call_4p229rk6"
       }
     },
     "sessionUpdate": "tool_call_update",
     "status": "in_progress",
     "toolCallId": "call_4p229rk6"
   }
 }
*/

func (p *SecretsFilter) handleSessionID(sessionId acp.SessionId) {
	if sessionId != p.currentSession {
		p.currentSession = sessionId
		p.promptRequestBuf.Reset()
		p.agentMessageBuf.Reset()
		p.agentThoughtBuf.Reset()
		p.isCancelled = false
	}
}

func (p *SecretsFilter) OnSessionNotification(notif acp.SessionNotification) ([]*protobyss.ACPContainer, error) {
	if p.isCancelled {
		return []*protobyss.ACPContainer{}, nil
	}
	p.handleSessionID(notif.SessionId)

	agentMessage := ""
	piTerminalOutput := ""

	if notif.Update.ToolCallUpdate != nil && notif.Update.ToolCallUpdate.Meta != nil {
		if terminalOutput, ok := notif.Update.ToolCallUpdate.Meta["terminal_output"]; ok {
			if terminalOutputStr, ok := terminalOutput.(map[string]any); ok {
				if message, ok := terminalOutputStr["data"]; ok {
					if messageStr, ok := message.(string); ok {
						piTerminalOutput = messageStr
					} else {
						p.logging.Info(context.Background(), &protobyss.LogMessage{
							Message: fmt.Sprintf("data %T", message),
							Fields:  map[string]string{},
						})
					}
				} else {
					p.logging.Info(context.Background(), &protobyss.LogMessage{
						Message: "data not in terminal_output",
						Fields:  map[string]string{},
					})
				}
			} else {
				p.logging.Info(context.Background(), &protobyss.LogMessage{
					Message: fmt.Sprintf("terminal_output is not a map: %T", terminalOutput),
					Fields:  map[string]string{},
				})
			}
		} else {
			p.logging.Info(context.Background(), &protobyss.LogMessage{
				Message: "no terminal_output in message",
				Fields:  map[string]string{},
			})
		}
	}

	if notif.Update.AgentMessageChunk != nil && notif.Update.AgentMessageChunk.Content.Text != nil {
		agentMessage = notif.Update.AgentMessageChunk.Content.Text.Text

		p.agentMessageBuf.Append(agentMessage)

		if p.agentMessageBuf.maxSize > len(agentMessage) {
			agentMessage = p.agentMessageBuf.String()
		}
	}

	agentThought := ""

	if notif.Update.AgentThoughtChunk != nil && notif.Update.AgentThoughtChunk.Content.Text != nil {
		agentThought = notif.Update.AgentThoughtChunk.Content.Text.Text
		p.agentThoughtBuf.Append(agentThought)

		if p.agentThoughtBuf.maxSize > len(agentThought) {
			agentThought = p.agentThoughtBuf.String()
		}
	}

	for _, v := range p.detectors {
		if agentMessage != "" && v.Match(agentMessage) {
			return p.getCancelMessages(v)
		}
		if agentThought != "" && v.Match(agentThought) {
			return p.getCancelMessages(v)
		}
		if piTerminalOutput != "" && v.Match(piTerminalOutput) {
			newStatus := acp.ToolCallStatusFailed
			newUpdate := acp.SessionNotification{
				SessionId: notif.SessionId,
				Update: acp.SessionUpdate{
					ToolCallUpdate: &acp.SessionToolCallUpdate{
						Status:     &newStatus,
						ToolCallId: notif.Update.ToolCallUpdate.ToolCallId,
					},
				},
			}

			statusMsg, err := abyss.ACPContainer(newUpdate)
			if err != nil {
				return nil, err
			}
			msgs, err := p.getCancelMessages(v)
			if err != nil {
				return nil, err
			}
			return append([]*protobyss.ACPContainer{statusMsg}, msgs...), nil
		}
	}

	return abyss.ACPContainers(notif)
}

func (p *SecretsFilter) OnTerminalOutputResponse(resp acp.TerminalOutputResponse) ([]*protobyss.ACPContainer, error) {
	for _, v := range p.detectors {
		if v.Match(resp.Output) {
			return p.getCancelMessages(v)
		}
	}
	return abyss.ACPContainers(resp)
}

func (p *SecretsFilter) OnReadTextFileResponse(resp acp.ReadTextFileResponse) ([]*protobyss.ACPContainer, error) {
	for _, v := range p.detectors {
		if v.Match(resp.Content) {
			return p.getCancelMessages(v)
		}
	}
	return abyss.ACPContainers(resp)
}

func (p *SecretsFilter) OnPromptRequest(req acp.PromptRequest) ([]*protobyss.ACPContainer, error) {
	p.handleSessionID(req.SessionId)
	allStringContents := strings.Builder{}

	for _, v := range req.Prompt {
		if v.Text != nil {
			allStringContents.WriteString(v.Text.Text)
		}
	}

	contents := allStringContents.String()
	p.promptRequestBuf.Append(contents)

	if len(contents) < p.promptRequestBuf.maxSize {
		contents = p.promptRequestBuf.String()
	}

	for _, v := range p.detectors {
		if v.Match(contents) {
			resp := acp.PromptResponse{
				StopReason: acp.StopReasonEndTurn,
			}

			respCont, err := abyss.ACPContainer(resp)
			if err != nil {
				break
			}

			msg := acp.SessionNotification{
				SessionId: req.SessionId,
				Update: acp.SessionUpdate{
					AgentMessageChunk: &acp.SessionUpdateAgentMessageChunk{
						Content: acp.TextBlock("This prompt matches security filter " + v.Name() +
							"The agent will ignore your message."),
					},
				},
			}
			msgCont, err := abyss.ACPContainer(msg)
			if err != nil {
				break
			}

			return []*protobyss.ACPContainer{respCont, msgCont}, nil
		}
	}

	return abyss.ACPContainers(
		req,
	)
}

func main() {
}

type SecretDetector interface {
	Name() string
	Match(string) bool
	MaxLength() int
}

type BannedRegex struct {
	name      string
	Regex     *regexp.Regexp
	maxLength int
}

func (b *BannedRegex) Name() string {
	return b.name
}

func (b *BannedRegex) Match(s string) bool {
	return b.Regex.MatchString(s)
}

func (b *BannedRegex) MaxLength() int {
	return b.maxLength
}

func init() {
	detectors := []SecretDetector{
		&BannedRegex{
			// test string sk-or-v1-abcd-efg-hiklkmno-p-qrs-tuv-wxyz
			name:      "openrouter",
			Regex:     regexp.MustCompile(".*sk-or-v1-[A-Za-z0-9_-]{32,128}.*"),
			maxLength: 140,
		},
	}

	maxLength := 0
	for _, d := range detectors {
		if l := d.MaxLength(); l > maxLength {
			maxLength = l
		}
	}

	myPlugin := &SecretsFilter{
		detectors:        detectors,
		promptRequestBuf: NewSizedStringBuffer(maxLength),
		agentMessageBuf:  NewSizedStringBuffer(maxLength),
		agentThoughtBuf:  NewSizedStringBuffer(maxLength),
	}

	plugin := abyss.NewACPPluginRouter(myPlugin)
	protobyss.RegisterACPPlugin(plugin)
}
