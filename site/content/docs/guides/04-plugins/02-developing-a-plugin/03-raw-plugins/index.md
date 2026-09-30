---
title: "Raw Plugins"
description: "A walkthrough of how plugins work in Abyss."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2026-09-29T16:04:48+02:00
draft: false
weight: 3
toc: true
---

## Project layout

Make a new folder for your plugin somewhere convenient. Inside it, you only need a single file
called `main.go`. There's no `go.mod` file to manage because your plugin imports packages from abyss
itself, and the build command points Go at the abyss source tree.

A typical folder looks like this:

```text
my-plugin/
└── main.go
```

Every plugin file starts with the same three things: a build tag, a package declaration, and an
empty `main` function. The build tag (`//go:build wasip1`) tells Go "compile this for the WASM
sandbox, not for my normal operating system." The empty `main` looks pointless, but the plugin
machinery requires it to exist even though it does nothing — the real work happens in an `init`
function instead.

```go
//go:build wasip1

package main

func main() {}
```

## The Code

We're going to write a plugin that prints every message to the log as it flows past, then lets it
continue untouched. This is the simplest possible useful plugin, and it's a great way to see the
shape of things.

### Step 1 — Imports and registration

Open `main.go` and replace its contents with this:

```go
//go:build wasip1

package main

import (
	"context"
	"fmt"

	"github.com/SethCurry/abyss/pkg/protobyss"
)

func main() {}

func init() {
	// Hand an instance of our plugin to abyss at startup.
	protobyss.RegisterACPPlugin(&LoggerPlugin{})
}

// Compile-time check: if we ever forget to implement the interface,
// the build will fail right here instead of at runtime.
var _ protobyss.ACPPlugin = (*LoggerPlugin)(nil)

type LoggerPlugin struct{}
```

The interesting line is `protobyss.RegisterACPPlugin(...)`. That's how a plugin says "here I am,
please send messages my way." You call it once, inside the `init` function, and pass it an instance
of your plugin struct. The `var _ ...` line is a Go trick that makes the compiler verify we really
have implemented everything abyss expects — these days that's two methods, `Initialize` and
`HandleMessage`. If we've missed either one, we find out now rather than when the plugin loads.

### Step 2 — Say hello with `Initialize`

First up is the handshake we met earlier. Ours doesn't do anything fancy — our logger just
introduces itself. Add this to the bottom of your file:

```go
// Initialize is called once when abyss loads the plugin, before
// any messages are handled.
func (p *LoggerPlugin) Initialize(
	ctx context.Context,
	req *protobyss.ACPPluginInitializeRequest,
) (*protobyss.ACPPluginInitializeResponse, error) {
	return &protobyss.ACPPluginInitializeResponse{
		Name: "logger",
	}, nil
}
```

Abyss calls this once, right after loading your `.wasm` file and before handing over a single
message. Our version ignores the `OnHost` flag and the `Config` bytes and simply replies with the
name `"logger"`. If we'd forgotten this method entirely, the build would have failed back at the
`var _ protobyss.ACPPlugin = (*LoggerPlugin)(nil)` line in Step 1 — that check covers both methods
of the interface.

### Step 3 — The handler

The other half of the interface is `HandleMessage`, the method that gets to see every message. It
receives one message and returns a list of messages. Add it to the bottom of your file:

```go
// HandleMessage is called for every ACP message, in either direction.
func (p *LoggerPlugin) HandleMessage(
	ctx context.Context,
	msg *protobyss.ACPContainer,
) (*protobyss.ACPContainerList, error) {
	// msg.Content is the raw body of the message, JSON-encoded.
	fmt.Printf("saw a message: %s\n", string(msg.Content))

	// Pass the message through untouched by returning it inside a list.
	return &protobyss.ACPContainerList{
		Containers: []*protobyss.ACPContainer{msg},
	}, nil
}
```

That's the whole plugin. The `msg` you receive is an `ACPContainer`, which is the generic envelope
every message arrives in. It has four fields:

- **`TypeId`** — a number identifying which kind of ACP message this is (a prompt request, a session
  notification, etc.). You can look one up with `abyss.GetMessageTypeByID(msg.TypeId)` if you want
  to know what you're holding.
- **`MessageId`** — a unique ID for this particular message.
- **`ResponseFor`** — on replies, the `MessageId` of the request this is answering. Empty otherwise.
- **`Content`** — the message body itself, as a slice of bytes containing JSON. You can print it,
  unmarshal it into a typed struct, or hand it back untouched.

In this example we don't care which type it is — we just print the raw JSON and pass it along. To
pass a message through, you return it wrapped in an `ACPContainerList`. To drop it instead, you'd
return an empty list (`&protobyss.ACPContainerList{}`). To replace it, you'd return a different
container. To split it into several, you'd put several containers in the list.

One important difference from the typed router (which we'll meet next): when you go the raw route,
abyss does *not* fill in message IDs for you. If you invent brand-new messages, you're responsible
for setting `MessageId` and `ResponseFor` yourself. For pass-through and drop cases like this logger,
you don't have to think about it.

### Step 4 — Build it

Compile your plugin with this command, run from inside the plugin's folder:

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go
```

Let's decode that, because it looks weird:

- `GOOS=wasip1 GOARCH=wasm` tells Go "build for the WASI sandbox instead of my normal operating
  system." `wasip1` is the name of the specific WASM environment abyss uses.
- `-o plugin.wasm` is the output filename. You can name it anything; `plugin.wasm` is just a
  convention. Point your config at whatever you call it.
- `-buildmode=c-shared` is a Go flag that produces the kind of module the plugin loader expects.

When the command finishes, you'll have a `plugin.wasm` file next to your `main.go`. Point your
config's `plugins.client.path` at it, restart abyss, and you should see your `saw a message:` lines
scrolling by in the logs.

> **Heads up:** if you're used to Go, you might be tempted to drop a `go.mod` file in your plugin
> folder and tidy it. You don't need one, and adding one can cause the build to look for abyss in the
> wrong place. Build from inside the abyss repository, or point `GOPATH`/module replacements at it,
> and you'll be fine.

## Flavor 2: The Typed Router (a Prompt Filter)

The raw interface is powerful but verbose: you have to decode every message yourself, even the ones
you don't care about. The typed router flips that around. You write a plain old struct and only add
methods for the message types you're interested in. Abyss does the decoding, calls only the methods
that match, and silently passes through everything else.

We'll write a plugin that looks at each prompt the user sends and refuses to pass it along if it
contains a banned word.

### Step 1 — Imports and the struct

Start the file the same way, but this time we also pull in the `abyss` package (for the router),
the `acp-go-sdk` package (which holds the typed message structs), and `context`, which our
`Initialize` method will need:

```go
//go:build wasip1

package main

import (
	"context"
	"regexp"
	"strings"

	"github.com/SethCurry/abyss/pkg/abyss"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/coder/acp-go-sdk"
)

func main() {}

// PromptFilter holds any state our plugin needs. Ours keeps a list of
// banned regular expressions to test prompts against.
type PromptFilter struct {
	bannedRegexes []*regexp.Regexp
}
```

### Step 2 — Say hello with `Initialize`

Just like the raw interface, the typed router starts with the handshake: every plugin needs an
`Initialize` method, and `abyss.NewACPPluginRouter` won't accept a struct that's missing one — the
compiler will tell you so in no uncertain terms. Add this to the bottom of your file:

```go
// Initialize is called once when abyss loads the plugin, before
// any messages are handled.
func (p *PromptFilter) Initialize(
	ctx context.Context,
	req *protobyss.ACPPluginInitializeRequest,
) (*protobyss.ACPPluginInitializeResponse, error) {
	return &protobyss.ACPPluginInitializeResponse{
		Name: "prompt_filter",
	}, nil
}
```

The request carries the same `OnHost` flag and `Config` bytes we met earlier; our filter has no use
for either, so it simply replies with its name, `"prompt_filter"`. When abyss loads the plugin and
comes knocking, the router passes the handshake straight through to this method — no extra wiring
on your part.

### Step 3 — Register it in `init`

As before, we create an instance and hand it to abyss. The difference is that we wrap our struct in
`abyss.NewACPPluginRouter(...)` first. That wrapper is what scans our struct for `On...` methods
and wires them up to the right message types, and it carries our `Initialize` method along too, so
the handshake "just works."

```go
func init() {
	bannedStrings := []string{".*SECRET.*"}
	regexes := make([]*regexp.Regexp, len(bannedStrings))
	for i, s := range bannedStrings {
		regexes[i] = regexp.MustCompile(s)
	}

	plug := &PromptFilter{bannedRegexes: regexes}
	protobyss.RegisterACPPlugin(abyss.NewACPPluginRouter(plug))
}
```

Here we've hardcoded the banned patterns as `".*SECRET.*"`, but you could read them from the message
itself, load them another way, or swap the regexes for a plain string check — it's all just Go.

### Step 4 — Implement the methods you care about

We only care about prompts, so the only `On...` method we implement is `OnPromptRequest`. The
router figures out the rest on its own.

```go
// OnPromptRequest is called whenever the user sends a prompt to the agent.
func (p *PromptFilter) OnPromptRequest(
	req acp.PromptRequest,
) ([]*protobyss.ACPContainer, error) {
	// Gather all the text from the prompt into one string.
	allText := strings.Builder{}
	for _, v := range req.Prompt {
		if v.Text != nil {
			allText.WriteString(v.Text.Text)
		}
	}

	// If any banned pattern matches, refuse the prompt.
	for _, re := range p.bannedRegexes {
		if re.Match([]byte(allText.String())) {
			// Tell the agent we're refusing to continue.
			resp := acp.PromptResponse{StopReason: acp.StopReasonRefusal}
			respContainer, err := abyss.ACPContainer(resp)
			if err != nil {
				break
			}

			// Send the user a message explaining what happened.
			notice := acp.SessionNotification{
				SessionId: req.SessionId,
				Update: acp.SessionUpdate{
					AgentMessageChunk: &acp.SessionUpdateAgentMessageChunk{
						Content: acp.TextBlock(
							"\n\n\nNuh uh, not under my roof!\n" +
								"You have violated a security filter.",
						),
					},
				},
			}
			noticeContainer, err := abyss.ACPContainer(notice)
			if err != nil {
				break
			}

			// Returning these two replaces the original prompt
			// entirely — it never reaches the agent.
			return []*protobyss.ACPContainer{respContainer, noticeContainer}
		}
	}

	// No match: pass the original prompt through untouched.
	return abyss.ACPContainers(req)
}
```

A few things to notice:

- The argument (`req acp.PromptRequest`) arrives already decoded into a proper Go struct, courtesy
  of the router. You never touch raw bytes.
- The return type is the same `[]*protobyss.ACPContainer` as before. The helper `abyss.ACPContainers(req)`
  is the easy way to turn a typed struct back into that slice when you just want to pass it through.
  Use `abyss.ACPContainer(thing)` (singular) when you're building a single new message to return.
- Returning two containers, as we do in the refusal branch, *replaces* the original prompt with both
  of them. The original prompt is gone — only what you return continues down the pipe.
- Unlike the raw interface, the router fills in `MessageId` and `ResponseFor` for you on new
  messages you synthesize, so you don't have to set them yourself.

### Step 5 — Which methods can I implement?

`NewACPPluginRouter` looks at your struct and connects any of a long list of `On...` methods it
finds. (`Initialize` isn't one of these — it isn't a message handler, it's the required handshake
from Step 2.) A handful of the common ones:

| Method | Fires when… |
| --- | --- |
| `OnPromptRequest` | the user sends a new prompt. |
| `OnPromptResponse` | the agent finishes responding to a prompt. |
| `OnSessionNotification` | the agent streams a chunk of a reply. |
| `OnNewSessionRequest` | a new session is being created. |
| `OnWriteTextFileRequest` | the agent asks to write a file. |
| `OnReadTextFileRequest` | the agent asks to read a file. |
| `OnCreateTerminalRequest` | the agent asks to start a shell command. |
| `OnRequestPermissionRequest` | the agent asks permission to do something. |

The full list lives in [`pkg/abyss/plugin_router_builder.go`](https://github.com/SethCurry/abyss/blob/main/pkg/abyss/plugin_router_builder.go)
in the abyss source tree — one interface per ACP message type. You only ever implement the ones you
care about; the rest are silently ignored.

The handler signature is always the same shape:

```go
func (p *MyPlugin) OnSomething(req acp.Something) ([]*protobyss.ACPContainer, error)
```

### Step 6 — Build it

Same command as before, from inside the plugin's folder:

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go
```

Add it to your config under `plugins.client`, restart abyss, and try sending a prompt that contains
the word `SECRET`. Instead of reaching the agent, you'll see your refusal message come back.

## Logging From Inside a Plugin

A plugin runs in a sandbox, so the usual Go logging habits (`fmt.Println`, `log.Printf`) work, but
their output goes to the plugin's own stdout rather than into abyss's log stream. If you want your
messages to show up alongside everything else abyss logs — the best place to look when something
goes wrong — use the host logging functions abyss provides.

Grab a logger once with `protobyss.NewLogging()`, then call `Debug`, `Info`, `Warn`, or
`Error` on it. Each takes a `protobyss.LogMessage` with two fields: `Message` (the text) and
`Fields` (an optional map of string keys to string values, which show up as structured fields in
the log line). One-time setup like this also fits nicely inside your `Initialize` method, since it
runs before the first message arrives — the example below does it in `init` instead, which works
just as well.

Here's the `PromptFilter` from above, updated to stash a logger on its struct and announce itself
when it loads. The `"context"` import we added back in Step 1 covers the `context.Background()`
call:

```go
type PromptFilter struct {
	bannedRegexes []*regexp.Regexp
	log           protobyss.Logging
}

func init() {
	log := protobyss.NewLogging()

	log.Info(context.Background(), &protobyss.LogMessage{
		Message: "prompt filter loaded",
		Fields:  map[string]string{"version": "1.0"},
	})

	plug := &PromptFilter{
		bannedRegexes: []*regexp.Regexp{regexp.MustCompile(".*SECRET.*")},
		log:           log,
	}
	protobyss.RegisterACPPlugin(abyss.NewACPPluginRouter(plug))
}
```

Every log line your plugin emits is tagged with the plugin's path, so in the abyss logs you can tell
at a glance which plugin said what.

## Complete Examples in the Repository

Abyss ships two ready-to-run example plugins in the
[`example/plugins`](https://github.com/SethCurry/abyss/tree/main/example/plugins) directory of the
source tree. They're the same ideas we built above, polished and commented:

- **`global_logger`** — the raw-interface logger, for when you want to see every message regardless
  of type.
- **`prompt_filter`** — the typed-router prompt filter, for when you only care about a few message
  types.

Both already contain a built `plugin.wasm` and a `README.md` explaining the design choices they
make. Copying one of these folders is the fastest way to start your own plugin — rename the struct,
edit the handler methods, rebuild, and you're done.

## A Few Things to Keep in Mind

- **Keep the `//go:build wasip1` line at the very top of `main.go`.** Without it, the build produces
  a normal program instead of a WASM module and abyss won't be able to load it.
- **Keep the empty `main` function.** The plugin loader needs the file to be a `main` package even
  though `main` itself does nothing; all the real startup happens in `init`.
- **Every plugin must implement `Initialize`.** Abyss calls it once, as soon as it loads your
  plugin, and your plugin answers with its name. Both flavors check for the method at build time,
  so a plugin without it won't compile — let alone load.
- **Plugins run on the client side**, before messages enter the container. That means they can see
  and shape what your agent is allowed to do, but they can't see anything happening *inside* the
  container that doesn't come back out as an ACP message.
- **Plugins run in order.** When you list several, each one's output becomes the next one's input.
  Put the plugin you want to have the first or last word in the matching position.
- **Plugins are sandboxed.** They can't read your filesystem, make network calls, or spawn
  processes unless abyss grants them the ability. Logging is the one host ability exposed today.
