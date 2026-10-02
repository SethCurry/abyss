---
title: "What is abyss?"
description: "Learn what abyss is, what it does, what problems it solves, and why you would want to use it."
summary: ""
date: 2026-08-01T16:04:48+02:00
lastmod: 2026-08-01T16:04:48+02:00
draft: false
weight: 1
toc: true
---

## The 30,000 Foot Overview

{{< mermaid >}}

architecture-beta
    group host(pixel:computer-old-electronics)[Host]

    group container(pixel:logo-social-media-dropbox)[Container] in host
    
    service acpclient(pixel:content-files-notepad)[ACP Client like Zed] in host
    service hostproxy(pixel:internet-network-computer-upload)[Abyss Host Proxy] in host
    service containerproxy(pixel:internet-network-computer-download)[Abyss Container Proxy] in container
    service agent(pixel:business-products-network-user)[Agent] in container

    acpclient:R -- L:hostproxy
    hostproxy:R -- L:containerproxy
    containerproxy:R -- L:agent

{{< /mermaid >}}

Abyss is 2 things at heart:

- A system for creating Docker containers running agents, like DevContainers for agents
- A stdio to WebSocket proxy with a plugin system, like nginx for agent traffic

The goal of those and Abyss generally is to provide:

- A secure environment to run agents in
- An easy way to manage those environments
- Agent and client agnostic features via ACP middleware

Those goals come out of my own experiences and pains.

### I Am Not A Babysitter

I'm unwilling to babysit my agents enough to manage
dozens of permission prompts, but I can't accept
them having unfettered access to my desktop.
I want to give them a walled garden that they
can play in and destroy.

I want something secure enough to hit "Accept All"
on my permission prompts.

Here's Abyss ending the agent's session and informing
the user after the agent started regurgitating
an API key.  This works with _any_ agent:

![Abyss ending a session after the agent starts reproducing an API Key](./abyss-secrets-refusal-example.jpg)

### I Am Impatient

I don't want a tool that is always in my way.
I don't want to update a Dockerfile 4 times a day,
or realize that copies of my agent image are taking
up 100GB of space on my machine.

I want a tool that I forget I'm using.


See how long it takes you to notice which one
is Abyss (if you don't look at the agent name).

{{< asciinema url="/asciinema/basic-example-pi.cast" >}}

{{< asciinema url="/asciinema/basic-example-abyss-pi.cast" >}}


### I Want To Pick My Own Tools

A lot of what Abyss does could be implemented as
plugins in the agents themselves.  I thought about
that, but I don't want to be eternally tied to
a particular editor or agent.  I want something
agnostic to both.


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
