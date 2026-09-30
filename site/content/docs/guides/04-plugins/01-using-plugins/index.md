---
title: "Using Plugins"
description: "Guides on installing, using and building plugins for abyss."
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


{{< admonition type="warning" title="Warning: Experimental" >}}
Plugins are an experimental feature and are liable to change.
I will try to keep the public APIs stable, but it is possible they will change
in the future.
{{< /admonition >}}


Plugins are little add-on programs that sit between your editor and your agent and get to look at
and even change every message that passes between the two. Want to block prompts that contain a
secret word? Rewrite a file path before the agent ever
sees it? A plugin can do all of that.

## Using an Existing Plugin

Plugins are tiny standalone files (they end in `.wasm`). To use one, you only need to tell Abyss
where the file lives. You do that in the same config file you created in
[Getting Started](02-getting-started.md), in a section called `plugins`.

Open your config file and add a `plugins` block that looks like this:

```yaml
docker:
  image: "ghcr.io/sethcurry/abyss-pi:latest"
  agent_command:
    - pi-acp
  host_mounts:
    - source: "./"
    - source: "~/.pi"
      destination: "/root/.pi"

plugins:
  client:
    - path: ./my-plugins/prompt_filter.wasm
```

The new part is everything under `plugins:`. Let's unpack it:

- **`plugins`** is the top-level section that holds all of your plugin settings.
- **`client`** is the list of plugins that run on *your* computer (the "client" side), before
  messages are sent into the container. Today every plugin is a client plugin, so this is always
  where they go.
- **`path`** is the location of the `.wasm` file on your computer. It can be a relative path (like
  the example above, relative to your config file) or an absolute one (like
  `/home/you/plugins/prompt_filter.wasm`).

### Loading More Than One

You can list as many plugins as you like. They run one after the other, in the order you wrote them,
like a bucket brigade: the first plugin hands its result to the second, the second to the third, and
so on, until the last one hands the message off to the agent (or back to your editor).

```yaml
plugins:
  client:
    - path: ./my-plugins/prompt_filter.wasm
    - path: ./my-plugins/global_logger.wasm
    - path: /home/you/plugins/audit-trail.wasm
```

If a plugin *drops* a message (we'll see how in a moment), the message never makes it to the plugins
after it. Order matters, so put the plugin you trust most first if you want it to have the final say
on what gets through.

### Checking That It Loaded

The easiest way to confirm your plugin is being picked up is `abyss oneshot`, which we met in
[Getting Started](02-getting-started.md). When abyss loads a plugin it writes a line to the logs
that looks like `loading ACP plugin` with the path next to it. If you see that line, you're in
business. If you instead see an error mentioning the plugin path, double-check the path is correct
and that the file really exists there.

```bash
abyss oneshot -f ./abyss-agent.yaml "What is the capital of France?"
```

That's all there is to using one. The rest of this guide is about *building* your own.

## Where to Go Next

- [Custom Docker Images](03-custom-docker-images.md) — pair your plugin with a custom agent image.
- [Configuration](../reference/configuration.md) — the full reference for every config option,
  including `plugins`.
- [Troubleshooting](05-troubleshooting.md) — when your plugin loads but doesn't behave, start here.
