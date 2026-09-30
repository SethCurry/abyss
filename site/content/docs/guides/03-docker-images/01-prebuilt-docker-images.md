---
title: "Prebuilt Docker Images"
description: "A list of the prebuilt Docker images Abyss provides, along with an example config and a list of the pre-installed software."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2026-09-22T16:04:48+02:00
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

This is a list of all of the Docker images pre-built for use with Abyss.
All images are based on the [ubuntu](https://hub.docker.com/_/ubuntu) Docker image.

## ghcr.io/sethcurry/abyss-base

This is the base image with Abyss installed, but nothing else.

You would mostly use this when you want to build your own custom image,
as [the next section](./02-custom-docker-images) talks about.

## ghcr.io/sethcurry/abyss-pi

This is the [Pi](https://pi.dev/) image.  It comes with Abyss and Pi
pre-installed, as well as `node` and `npm`.

## ghcr.io/sethcurry/abyss-hermes

This is the [Hermes](https://github.com/nousresearch/hermes-agent) image.
It comes with Hermes pre-installed, as well as `python3`, `uv`, and the
Hermes ACP extension.

## ghcr.io/sethcurry/abyss-codex

This is the [Codex](https://openai.com/codex/) image.  It comes with [codex-acp](https://github.com/agentclientprotocol/codex-acp) pre-installed, as well as `node` and `npm`.

## ghcr.io/sethcurry/abyss-claude-agent

This is the [Claude Agent](https://claude.com/solutions/agents) image.  It comes with [claude-agent-acp](https://github.com/agentclientprotocol/claude-agent-acp) pre-installed, as well as `node` and `npm`.
