---
title: "Abyss v0.0.2"
description: "Abyss v0.0.2 release announcement"
summary: "Abyss 0.0.2 is now available for download! See what's new in this release."
date: 2026-09-10T00:00:00+00:00
lastmod: 2026-09-10T00:00:00+00:00
draft: false
weight: 50
categories: []
tags: []
contributors: []
pinned: false
homepage: false
params:
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

This release was primarily focused on:

A: Cleaning up some bugs in the new user onboarding

B: Cleaning up some of the AI spaghetti code into more maintainable patterns

I'm hoping to get back to new features after this.  My intent is to get base
images for a few more agents, and add some dynamicism to the containers via
ACP (i.e. seeing ReadFile calls for paths that only exist on the host, and asking
the user if they want to read from the host, copy the file into the container,
or just let it fail).

### New Features

- Added an `image_pull_policy` to the Docker config that mimics Docker Compose's `image_pull_policy` behavior
- Added mutual TLS authentication by default, using single-use certificates
  - There is a config flag to disable this if you want plaintext

### Fixes

- Docker images now get pulled by default if not present
- Setup scripts now correctly run before an ACP connection is established

### Other Changes

- Websocket messages are now wrapped in protobuf, allowing non-ACP messages to be transmitted
  - In the future this can be used for things like getting a shell inside the container
- Creating containers now uses a "builder" pattern with options, cutting down on some spaghetti code
