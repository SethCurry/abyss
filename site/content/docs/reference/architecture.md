---
title: "Architecture"
description: "An overview of the underlying parts of abyss, how they interact, what they do, and visual depictions of how it all fits together."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
weight: 4
toc: true
params:
  math: false # enable mathematical rendering
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---


abyss itself is 2 components that sit between your [ACP client](https://agentclientprotocol.com/get-started/introduction) like [Zed](https://zed.dev/) and
your agent.

The client proxy runs on your host and is what Zed invokes and connects to.
The client proxy then starts a Docker container based on your configuration.

The agent proxy runs inside the Docker container, and is what runs your agent and proxies stdio to the client proxy, and then your ACP client.

![Architecture](./architecture.png)
