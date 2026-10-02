---
title: "Setting Up To Build Images"
description: "The pre-work for building Abyss Docker images"
summary: ""
date: 2023-09-07T16:05:48+02:00
lastmod: 2026-09-22T16:04:48+02:00
draft: false
weight: 1
toc: true
---

The ready-made images on the [Docker Images]({{% ref "/docs/guides/03-docker-images/01-prebuilt-docker-images.md" %}}) list are great for
getting started, but eventually you might want something they don't ship. Maybe your agent needs a
language runtime they don't include, a specific version of a tool, or a private config file baked
in. The good news is that abyss is completely happy running inside an image *you* build yourself.

There are two ways to go about it, and we'll walk through both step by step:

- **Extend an existing image** — start from one of the images I already publish and add what you
  need on top. This is the easier path, and it's the one most people should take.
- **Build an entirely custom image** — start from a plain base (like Ubuntu) and install abyss
  yourself. You'd pick this if you want full control over what's inside, or you're targeting an
  unusual platform like ARM.

Don't worry if you've never written a Dockerfile before — we'll explain each line as we go, and by
the end you'll have a working image you can point abyss at.

## Before You Begin

You'll need a couple of things ready before we start.

### 1. Docker

You're going to be building Docker images, so you need Docker installed and running. The quickest check is to run:

```bash
docker run hello-world
```

If you see a friendly "Hello from Docker!" message, you're set. Otherwise, follow Docker's
[installation guide](https://docs.docker.com/get-started/get-docker/) and come back when it works.

> Make sure you can run `docker` commands *without* typing `sudo` first. If `sudo` is required, add
> your user to the `docker` group — the Docker installation guide explains how.

### 2. A Place to Work

Create an empty folder somewhere on your computer. This is where we'll put the small text file
(called a *Dockerfile*) that describes your image. It can live anywhere; a common spot is a new
folder inside your project:

```bash
mkdir my-agent-image
cd my-agent-image
```

### 3. The Two Rules

No matter which path you take, abyss only asks two things of the image it runs in:

1. **The `abyss` program must be at `/usr/local/bin/abyss`.** Abyss launches itself inside the
   container, so it needs to know exactly where to find its own binary.
2. **`bash` must be installed.** Abyss uses bash to run your setup scripts and the agent's startup
   command.

That's it. Everything else — the operating system, the architecture, the extra tools — is entirely
up to you. Keep those two rules in mind and you can't go far wrong.

## Which Option Should I Pick?

Still not sure which path to take? Here's a quick rule of thumb:

- **Extend an existing image** if one of my published images is *almost* what you want and you just
  need to add a few tools or files. It's faster, simpler, and you automatically pick up fixes when
  the base image is updated.
- **Build a custom image** if you need a different base operating system, a very small image, a
  platform I don't publish for, or you simply want full control over every byte inside.

Both paths produce an image abyss is happy to run. When in doubt, start with Option 1 — you can
always switch to Option 2 later.

## Where to Go Next

- [Configuration]({{% ref "/docs/reference/configuration/" %}}) — every config option, including `image_pull_policy`
  and `setup_scripts` for run-time customization.
- [Troubleshooting]({{% ref "/docs/guides/05-troubleshooting.md" %}}) — what to do when your image won't build or abyss can't
  start it.
