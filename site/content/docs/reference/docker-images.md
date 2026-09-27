---
title: "Docker Images"
description: "A list of the Docker images I provide to make using Abyss easier."
summary: ""
date: 2026-08-25T16:13:18+02:00
lastmod: 2026-09-27T19:51:49+00:00
draft: false
weight: 2
toc: true
params:
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

This is the master list of all of the Docker images I publish for Abyss to help you get started without
needing to build your own Docker container first.

The [Getting Started guide](../guides/02-getting-started.md) has a step-by-step onboarding section for
each agent image, including the `agent_command` to start it with.

| Image | Agent | `agent_command` | Description |
|-------|-------|-----------------|-------------|
| ghcr.io/sethcurry/abyss-base:latest | None | — | The base image with Abyss installed that other images are built on top of. |
| ghcr.io/sethcurry/abyss-pi:latest | [Pi](https://pi.dev/) | `pi-acp` | An image with both Abyss and Pi installed. See [onboarding](../guides/02-getting-started.md#using-pi). |
| ghcr.io/sethcurry/abyss-hermes:latest | [Hermes](https://hermes-agent.nousresearch.com/) | `hermes-acp` | An image with Abyss and Hermes installed. See [onboarding](../guides/02-getting-started.md#using-hermes). |
| ghcr.io/sethcurry/abyss-codex:latest | [Codex](https://github.com/openai/codex) | `codex-acp` | An image with Abyss, Codex, and the [codex-acp](https://github.com/agentclientprotocol/codex-acp) adapter installed. See [onboarding](../guides/02-getting-started.md#using-codex). |
| ghcr.io/sethcurry/abyss-claude-agent:latest | [Claude Agent](https://claude.com/product/claude-code) | `claude-agent-acp` | An image with Abyss, Claude Agent, and the [claude-agent-acp](https://github.com/agentclientprotocol/claude-agent-acp) adapter installed. See [onboarding](../guides/02-getting-started.md#using-claude-agent). |
