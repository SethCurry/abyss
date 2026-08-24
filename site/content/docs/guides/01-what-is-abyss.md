---
title: "What is abyss?"
description: "Learn what abyss is, what problems it solves, and why you would want to use it."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
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

Abyss is an _Agent Runtime Environment_, a term I coined to describe a
system for executing agents in an isolated and controlled environment.
It's similar to Docker's or Podman's container runtime (which we do use for
isolation), but for running an agent rather than an arbitrary application.

There are a lot of parallels and overlaps with containerization, but I
believe there is value in a system tailored to the specific needs of
agents.

## What does it do?

Abyss automates a lot of the tedium that makes manually managing agents
difficult.  It provides integrations to help set up the inside of the
container, a stdio -> websocket proxy for ACP to work with, and lifecycle management of your agent.

A quick overview of some helpful features of abyss:

- Handles bind-mounting directories from your PC
- Copies files from your PC into the container on startup
- Executes startup scripts inside the container

## Architecture Overview

abyss itself is 2 components that sit between your [ACP client](https://agentclientprotocol.com/get-started/introduction) like [Zed](https://zed.dev/) and
your agent.

The client proxy runs on your host and is what Zed invokes and connects to.
The client proxy then starts a Docker container based on your configuration.

The agent proxy runs inside the Docker container, and is what runs your agent and proxies stdio to the client proxy, and then your ACP client.

![Architecture](./architecture.png)
