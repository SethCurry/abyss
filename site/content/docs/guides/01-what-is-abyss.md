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

## Story Time

Have you ever had your agent do the wrong thing?  It edited the wrong file, or logged in to the wrong
server, or just thought today's date is 6 months ago and proceeded to wreak havoc, small or large.

I decided I needed to approve my agent's commands.

And that sucks.

My agent used to be independent; I could plan out the next task while it was working on the current one.
Now I can't think, because every second the LLM spends thinking results in 5 tool calls that I have to approve.
Once again, I was the slowest part of my own devlopment loop.

I couldn't help but think this would be much cleaner if I could give the agent its own "rage room"
of sorts, where I don't have to watch because I don't care if it trashes its own environment.

## Enter the Abyss

Abyss is an _Agent Runtime Environment_, a term I coined to describe a
system for executing agents in an isolated and controlled environment.
It's similar to Docker's or Podman's container runtime (which we do use for
isolation), but for running an agent rather than an arbitrary application.

The short version is that Abyss is a tool that makes it easy to run your agent in  Docker
container while still giving it the resources it needs, all in a YAML file shorter than
your standup update.

## What does it do?

Abyss handles creating and managing Docker containers running your agents, alongside proxying your ACP connection
over websocket.  The upshot is that your agent can run in an isolated container (or one day even on another computer)
while neither your editor nor agent are aware.

Abyss can also:

- Stop and remove Docker containers when you disconnect from ACP
- Add bind-mounts to your agent container, so your agent can edit the same files you see in your editor
- Copy files from your PC into the container at runtime, so the agent gets a fresh copy of the files but the agent can't modify the version on your PC
- Execute startup scripts before the agent starts to install tools, dependencies, clone git repos, or whatever else.

### What will it do one day?

The current feature set is nice and useful, but Abyss wouldn't be a real project without scope creep.
This is a laundry list of features I would like to add.  Follow the blogs if you want updates!

- Allow the agent container to run on another computer, e.g. via Kubernetes or SSH tunneling
- Support for injecting RAG in a manner similar to what Docker or k8s do for volumes
  - This could allow for rapid agent development by siloing RAG into "volumes" and building an agent by combining several RAG sources
- ACP middleware features akin to what Nginx does for reverse proxies
  - Log requests centrally
  - Implement a "firewall" of sorts that can filter inbound prompts or outbound responses
  - Dynamically route requests to different agents based on code, a classifier, etc.

This is not a promise these things will happen, they're just where I currently see room for improvements.
They might get built, they might get delayed in favor of something else, or I may decide they're not a good idea
or out of scope.
