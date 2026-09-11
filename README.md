# abyss

![Abyss Logo](./site/static/android-chrome-192x192.png)

[![Code Quality](https://github.com/SethCurry/abyss/actions/workflows/go-test.yml/badge.svg)](https://github.com/SethCurry/abyss/actions/workflows/go-test.yml)

Check out [the docs](https://abyss.scurry.io) for in-depth information.

`abyss` is an _Agent Runtime Environment_ — a platform that runs your LLM agents
inside isolated containers, the same way Docker runs your compute.

Writing Docker Compose files to isolate your agent is a pain, and trying to
connect Zed to a containerized agent is even worse.

`abyss` handles all of
that. You describe what your agent needs in a short YAML file and `abyss` brings
up a container, exposes the agent over the [Agent Client Protocol](https://agentclientprotocol.com/)
(ACP), and tears it all down when you're done — no Docker expertise required.

`abyss` works with any ACP-compatible client, like Zed.  It integrates
seamlessly and ensures that directories you mount show up at the same path
in the container so you don't need to worry about memorizing a weird path scheme.

No more worrying about the agent finding the keys for the prod database on your
desktop, or `rm -rf`'ing your entire home directory.

## Why abyss?

- **No Docker configs to write.** A few lines of YAML describe the image, mounts,
  and the command that launches your agent. `abyss` handles the rest.
- **Sandboxed by default.** Each agent runs in its own container. Optionally copy
  files in instead of bind-mounting, so agent edits never touch your working copy.
- **ACP out of the box.** `abyss` proxies agent stdio over websockets and exposes
  an ACP endpoint, so any ACP client — Zed included — can drive the agent.
- **File and terminal interception.** ACP read/write file and terminal APIs are
  intercepted and executed inside the container, so the agent and the client agree
  on a single, consistent filesystem.
- **Reproducible environments.** Run setup scripts before the agent starts to
  install dependencies or seed state, and every session begins from a known place.

## Status

Abyss has had its first MVP release!

You can grab a copy of the binary from the [releases page](https://github.com/SethCurry/abyss/releases),
and Docker images are under Packages on the right of the project home.

## Features

- Starting a Docker container with your agent
- Proxying the agent's stdio over websocket to your ACP client
- Running setup scripts before starting the agent
- Bind-mounting directories from the host into the container
- Copying files into the container (so agent edits don't impact your copy)
- Intercepting ACP read/write file and terminal APIs so they run inside the container

## Vague and Unorganized TODO

These are not done, but are a laundry list of things I would like to accomplish:

- Allow agents to communicate with each other (requires them to be able to ACP to each other)
- Authentication on client-server comms
- Encryption on client-server comms
- Per-agent storage and shared storage
  - Unsure what this looks like.  Is it RAG?  Is it literal directories?  Both?
