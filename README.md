# abyss

![Abyss Logo](./site/static/android-chrome-192x192.png)

[![Code Quality](https://github.com/SethCurry/abyss/actions/workflows/go-test.yml/badge.svg)](https://github.com/SethCurry/abyss/actions/workflows/go-test.yml)

Check out [the docs](https://abyss.scurry.io) for in-depth information.

## Stop handing your whole machine to an LLM you don't fully control.

`abyss` is a system for running your agents inside a Docker container,
without your editor or your agent being aware that they're not in the same place.

It creates a Docker container running your agent and proxies your editor's
connection into that container while creating easy, reusable facilities
for copying files, bind-mounting directories, running setup scripts and more!

Abyss also supports WASM-based plugins that are given full control over messages
flowing back and forth.  You can reject messages containing secrets, implement
tools at the ACP layer so they work with any agent and more.

![Abyss Architecture](./site/content/reference/architecture.png)

## Your agent, sandboxed in seconds.

Writing Docker Compose files to isolate your agent is a pain, and trying to
connect Zed to a containerized agent is even worse.

`abyss` handles all of that. You describe what your agent needs in a short YAML
file and `abyss` brings up a container, exposes the agent over the
[Agent Client Protocol](https://agentclientprotocol.com/) (ACP), and tears it all
down when you're done — no Docker expertise required.

`abyss` works with any ACP-compatible client, like Zed. It integrates seamlessly
and ensures that directories you mount show up at the same path in the container,
so you don't need to memorize a weird path scheme.

## Why abyss?

- **Sandboxed by default.** Each agent runs in its own container — your machine is
  never exposed to whatever the agent decides to run. Optionally copy files in
  instead of bind-mounting, so agent edits never touch your working copy.
- **No Docker configs to write.** A few lines of YAML describe the image, mounts,
  and the command that launches your agent. `abyss` handles the rest.
- **File and terminal interception.** ACP read/write file and terminal APIs are
  intercepted and executed *inside* the container, so the agent and the client
  agree on a single, consistent filesystem — and the agent never escapes it.
- **Ephemeral security by design.** Ephemeral mutual-TLS authentication; certificates
  are used for a single connection and then destroyed.
- **ACP out of the box.** `abyss` proxies agent stdio over websockets and exposes
  an ACP endpoint, so any ACP client — Zed included — can drive the agent.
- **Reproducible environments.** Run setup scripts before the agent starts to
  install dependencies or seed state, and every session begins from a known place.
- **WASM-based Plugins.** Use existing plugins or write your own, plugins have
  full access to the stream of ACP messages. See [the examples](./example/plugins).

## Sleep easier. Ship faster.

`abyss` is the firewall between your agent and your machine. Run the most capable
models you can find, give them the most aggressive prompts you want — and know
that the worst an agent can do is wreck its own sandbox, not your system.

## Supports Batch Processing

You can use the `abyss oneshot` command to execute a single prompt, allowing you to execute
bash/Python/etc cronjobs that call your agent while maintaining its isolation.

## Releases

You can grab a copy of the binary from the [releases page](https://github.com/SethCurry/abyss/releases),
and Docker images are under Packages on the right of the project home.

## Features

- Starting a Docker container with your agent
- Proxying the agent's stdio over websocket to your ACP client
- Running setup scripts before starting the agent
- Bind-mounting directories from the host into the container
- Copying files into the container (so agent edits don't impact your copy)
- Intercepting ACP read/write file and terminal APIs so they run inside the container
- Ephemeral mutual-TLS authentication; certificates are used for a single connection and then destroyed.

## AI Policy

In the interest of full transparency, AI is used in the development of this application (shocking, I know).

At present, I do not feel that AI agents are capable of maintaining a large, clean codebase on their own.

Given that, AI is used in these two roles:

### Code Generation

Code generation is targeted, and always reviewed by a human.  I intentionally ask for very targeted features; exceptionally few of my edits
result in diffs larger than 100 lines.

I am still very much aware of and in control of the codebase.  I am aware there are some rough edges (very rough edges indeed) in parts of
the codebase.  Some of that is, as you may suspect, AI slop that was "good enough for an MVP" (see `runClient` in cmd/abyss/run_client.go).
Some of it is normal human fallibility as I iterated and realized abstractions were leaky, things weren't passed around well, etc
(see the horror that is how I make plugins work with the message-ID-based websocket RPC system in pkg/abyss/message_type.go).

### Documentation

If you're here, you're probably mad at me about using AI to generate documentation.

I do use AI to generate documentation.  I do review it, and when I have time I do try to edit out the "AI tone" it tends to use.
Some parts are human written, but at this point the docs are primarily AI.

If you'll indulge me, I do it for 3 reasons:

1. It writes better documentation than me.  I am terrible at writing documentation, because it's hard for me to "pretend" I don't know how this works and write for people who don't have my context.  The tone is annoying, but the information is far more complete than if I did it.
2. I am a solo maintainer, so time spent on documentation is time not spent on bug fixes/features/etc.  That's not to devalue documentation; it is incredibly important.  Combined with the above, it doesn't make sense for me to spend time there.  I will spend a lot more time than the LLM writing much worse documentation.
3. Things are still very much in flux, so there's a lot of documentation churn.  New features, reworked features, new or updated configs, etc.  Abyss isn't at the point where the docs are largely stable and we're just nit-picking about which phrasing is clearer.

If you still disagree with me about using AI to generate the docs, feel free to open a ticket.  I am not the target audience for the docs,
so my opinion doesn't mean a ton.  If you, the user, feel that the AI-written docs diminish your ability to use abyss, let me know and
I can take the time to do a hand-editing pass over the docs.
