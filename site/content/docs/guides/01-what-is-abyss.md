---
title: "What is abyss?"
description: "Learn what abyss is, what problems it solves, and why you would want to use it."
summary: ""
date: 2026-08-01T16:04:48+02:00
lastmod: 2026-08-01T16:04:48+02:00
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

## Set your agent free

Your agent shouldn't need a babysitter. But the moment you start approving its every command, you
become one. Abyss gives your agent the freedom to experiment — and fail — safely, so you can stay
focused on the work that actually needs you.

## Run agents without the risk

Abyss is an _Agent Runtime Environment_: a dedicated, isolated space where your agent can think,
build, and break things without touching your machine. It's like Docker for your workflow — except
instead of wrapping an application, it wraps your agent.

The idea is simple. Run your agent in a Docker container, hand it the resources it needs, and walk
away. All the configuration fits in a single YAML file that's shorter than your standup update.

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

Your agent gets a safe playground. You get your focus back.

## The road ahead

The best part? We're just getting started. Here's a taste of what's coming:

- Run agent containers on another machine — via Kubernetes or SSH tunneling
- Inject RAG the way Docker handles volumes, letting you compose agents from reusable sources
- ACP middleware inspired by reverse proxies: central logging, prompt/response filtering, and
  smart routing to different agents

Follow the [blogs](/blog/) for updates as these features land. Abyss isn't just a safer way to run
agents today — it's the foundation for how they'll run tomorrow.
