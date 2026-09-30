---
title: "What is abyss?"
description: "Learn what abyss is, what problems it solves, and why you would want to use it."
summary: ""
date: 2026-08-01T16:04:48+02:00
lastmod: 2026-08-01T16:04:48+02:00
draft: false
weight: 1
toc: true
---

## The 30,000 Foot Overview

Abyss runs your agent in a Docker container, and proxies messages between your
editor and your agent.

The container provides you with peace of mind that the agent is isolated
in a box.

The proxy enables middleware so you don't even have to fully trust your
agent.

{{< asciinema url="/asciinema/basic-example-pi.cast" >}}


After Abyss:

{{< asciinema url="/asciinema/basic-example-abyss-pi.cast" >}}

It has the container's hostname, but still has access to all of the files it needs!

## What makes it great

- **Real isolation.** Your agent works in a container, not on your host. Edits and commands stay
  where they belong.
- **Zero friction.** Abyss proxies your ACP connection over websocket, so neither your editor nor
  your agent even notice the isolation is there.
- **Seamless file access.** Bind-mount folders so your agent edits the exact files you see in your
  editor — or copy files in fresh, so your local copy stays untouched.
- **Ready-to-run environments.** Execute startup scripts to install tools, pull dependencies, and
  clone repos before your agent even starts.
- **Automatic cleanup.** Containers stop and remove themselves when you disconnect.
- **WASM plugins.**  Download plugins or build your own that have full access to all ACP messages.

Your agent gets a safe playground. You get your focus back.

## The road ahead

The best part? We're just getting started. Here's a taste of what's coming:

- Run agent containers on another machine — via Kubernetes or SSH tunneling
- Inject RAG the way Docker handles volumes, letting you compose agents from reusable sources
- ACP middleware inspired by reverse proxies: central logging, prompt/response filtering, and
  smart routing to different agents

Follow the [blogs](/blog/) for updates as these features land. Abyss isn't just a safer way to run
agents today — it's the foundation for how they'll run tomorrow.
