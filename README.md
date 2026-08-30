# abyss

![Abyss Logo](./site/static/android-chrome-192x192.png)

[![Code Quality](https://github.com/SethCurry/abyss/actions/workflows/go-test.yml/badge.svg)](https://github.com/SethCurry/abyss/actions/workflows/go-test.yml)

Check out [the docs](https://abyss.scurry.io) for in-depth information.

`abyss` is an _Agent Runtime Environment(s)_. The long-term goal is to provide a runtime platform
for agents, in a similar way to what Docker provides for compute.

`abyss` lets you create containers with copies or mounts of files from your desktop, talk to it
through Zed (or any ACP client), and then remove the container when you're done.

No more worrying about your agent finding the keys for the prod database on your desktop,
or `rm -rf`'ing your entire desktop.

## Status

Abyss has had its first MVP release!

You can grab a copy of the binary from the [releases page](https://github.com/SethCurry/abyss/releases),
and Docker images are under Packages on the right of the project home.

## Features

It currently supports:

- Starting a Docker container with your agent
- Proxying the stdio over websocket to your ACP client
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
