---
title: "How Plugins Work"
description: "A walkthrough of how plugins work in Abyss."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2026-09-29T16:04:48+02:00
draft: false
weight: 1
toc: true
params:
  math: false # enable mathematical rendering
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

This section is broken down into two parts:

1. **How plugins work**, so the rest of the guide makes sense.
2. **Building your own plugin**, with two complete walk-throughs.

Don't worry if you've never written a line of Go or heard the word "WASM" before. We'll go step by
step and explain everything as we meet it.

## How Plugins Work

Before we write any code, a quick tour of the moving parts. Skim this once now; come back to it
whenever a later step feels mysterious.

### Messages, messages everywhere

Your editor and your agent are constantly sending little JSON messages back and forth using a
shared language called the [Agent Client Protocol](https://agentclientprotocol.com/) (ACP). When you
type a prompt, your editor sends a *prompt request* message. When the agent replies, it sends a
*session notification* message. When the agent wants to run a shell command, it sends a
*create-terminal request* message. There are a few dozen message types in total, and every single
one of them funnels through the same pipe.

This is what one request/response pair looks like:

![Example ACP Connection](./acp-sequence-basic.png)

A plugin is a tap on that pipe. Every message that flows through abyss is handed to your plugin, and
your plugin gets to decide what happens next:

- **Pass it through** unchanged, as if the plugin weren't there.
- **Modify it** — tweak the contents and let the edited version continue.
- **Drop it** — swallow the message so nothing continues.
- **Replace it** — throw away the original and send something else instead.
- **Split it** — turn one message into several, each of which continues in order.

Here's an example with a plugin that cancels sessions
that mention a secret like an API key:

![Example ACP Sequence with Plugin](./acp-sequence-plugin-example.png)

### WASM, in one paragraph

A plugin is a small program compiled to a format called **WASM** (short for "WebAssembly", but it's
useful far beyond the web). WASM is a neat choice here for three reasons. First, it's sandboxed: a
plugin can't reach out and read your files or call the internet unless abyss explicitly lets it, so
running someone else's plugin is much safer than running a random script. Second, plugins are
written in [Go](https://go.dev/), which is the same language abyss itself is written in, so the
whole experience stays in one comfortable place.
Third, plugins aren't tied to using the exact version of everything Abyss does.
That's an issue with Go's native plugins.

You won't need to think about WASM day to day. The only place it shows up is the command you run to
turn your Go code into a plugin, which we'll get to shortly.

### The two flavors of plugin

Abyss gives you two ways to write a plugin, and which one you pick is mostly a matter of taste:

1. **The raw interface.** You implement two methods: `Initialize`, which every plugin must have,
   and `HandleMessage`, which receives *every* message, no matter its type, as a generic blob.
   You're responsible for figuring out what kind of message it is and decoding it yourself. This is
   the right choice when you want to see all traffic — a logger, a tracing tool, a generic
   middleware.

2. **The typed router.** You write a struct with a method for each message type you care about
   (`OnPromptRequest`, `OnSessionNotification`, and so on), plus the required `Initialize` method.
   Abyss decodes each message for you and only calls the methods that match. Message types you
   didn't implement simply pass through untouched. This is the right choice when you only care
   about a handful of message types — a prompt filter, a path rewriter, a permission policy.

We'll build one of each in this guide so you can see the difference firsthand.

### Saying hello: the `Initialize` handshake

Whichever flavor you pick, there's one method every plugin must have: `Initialize`. The
requirement comes from the shared contract that defines plugins (`schema/proto/abyss.proto` in the
abyss source tree), and the compiler enforces it — a plugin without an `Initialize` method won't
even build.

The idea is simple. Before abyss lets your plugin see a single message, it knocks on the door and
asks "who are you?" That knock is `Initialize`. Abyss calls it exactly once, right after loading
your `.wasm` file and before any messages flow through the pipe.

The request your `Initialize` method receives carries two pieces of information:

- **`OnHost`** — a true/false flag telling your plugin whether it's running directly on the host
  machine or somewhere else.
- **`Config`** — a set of bytes containing a JSON-marshalled copy of the plugin's
  options from the configuration file.

In return, your plugin hands back a response with a single field:

- **`Name`** — the name your plugin wants to be known by, like `"prompt_filter"` or
  `"global_logger"`.

Because `Initialize` runs before the first message arrives, it's also the natural home for any
one-time setup your plugin needs — compiling regular expressions, creating a logger, preparing
internal state, and so on.

## Up Next

The next page depends on whether you're writing a raw or typed plugin.  If you're writing a raw plugin, go see [Raw Plugins](./03-raw-plugins).
If you're writing a typed plugin, go see [Typed Plugins](./02-typed-plugins)
