---
title: "Architecture"
description: "An overview of the underlying parts of abyss, how they interact, what they do, and visual depictions of how it all fits together."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2026-09-26T17:04:00+02:00
draft: false
weight: 4
toc: true
params:
  math: false # enable mathematical rendering
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---


abyss is two components that sit between your [ACP client](https://agentclientprotocol.com/get-started/introduction) (such as [Zed](https://zed.dev/)) and
your agent. The client proxy runs on your host; the agent proxy runs inside a Docker container. The two halves talk to each other over a websocket carrying protobuf-encoded ACP messages, secured with ephemeral mutual TLS.

![Architecture](./architecture.png)

## Client proxy

The client proxy is the binary your ACP client invokes and connects to over stdio. When it starts it:

1. Loads your [agent configuration](./configuration.md).
2. Generates a throwaway CA and certificate set for mutual TLS (unless `websocket.disable_tls` is set).
3. Pulls the configured image and builds the container &mdash; bind mounts, copied files, setup scripts, and the TLS certificates are all installed before the container starts.
4. Drops a start file into the container to signal the agent proxy to begin serving.
5. Dials the agent proxy's websocket (`wss://.../ws`) and bridges the ACP client's stdio to it.

The host-side half of the ACP connection is implemented by the `HostProxy`, which also enforces abyss-specific behavior like ensuring the requested working directory exists before a new session starts. Any [WASM plugins](#plugins) configured under `plugins.client` are loaded here and see every ACP message on the host side.

## Agent proxy

The agent proxy is `abyss server`, started inside the container by the client proxy. It waits for the start file to appear, then serves an HTTPS websocket endpoint and spawns your agent from `docker.agent_command`, bridging the agent's stdio to the websocket.

The container-side half of the ACP connection is implemented by the `ContainerProxy`. By default it intercepts ACP `fs/read_text_file`, `fs/write_text_file`, and terminal requests and executes them inside the container, so the agent can never reach your host filesystem or shell. Setting `acp.tools_on_host.files` or `acp.tools_on_host.terminal` to `true` instead forwards those requests to the host to be executed by your ACP client &mdash; convenience that trades away isolation.

## Transport

All ACP traffic between the two proxies is serialized as protobuf and sent over a single websocket. Each proxy demultiplexes the stream, routes requests to the appropriate ACP handler, and matches responses back to their requests. The connection is authenticated with mutual TLS using certificates generated for that single session and destroyed afterwards.

## Plugins

abyss ships a WASM-based plugin system. Plugins are loaded on the client proxy (and the same mechanism exists on the agent proxy) and are invoked for every ACP message flowing in either direction. A plugin can inspect, filter, transform, or augment messages &mdash; for example rewriting prompts, redacting secrets, or logging the full message stream. See the [plugin examples](https://github.com/SethCurry/abyss/tree/main/example/plugins) for working plugins.

## Lifecycle

The client proxy owns the container for the lifetime of a single ACP session. When the ACP client disconnects (or the one-shot prompt completes), the client proxy stops and removes the container, leaving nothing behind. `abyss oneshot` reuses this same pipeline to run a single prompt batch-style against an ephemeral container.
