---
title: "Abyss v0.0.4 is now available"
description: "Abyss v0.0.4 is now available and adds support for persistent containers as well as some bug fixes."
date: 2026-10-06T00:00:00+00:00
lastmod: 2026-10-06T00:00:00+00:00
draft: false
weight: 50
categories: []
tags: []
contributors: []
pinned: false
homepage: false
---

Abyss v0.0.4 is now available!

My primary focus over the past week or two has been improving the
documentation.

Several sections were hand re-written by me to remove some of
the AI-ness as well as to add additional context I thought would be helpful.

Notably, tons of Markdown links were broken. I apologize for that,
it was a mess. I
fixed those, as well as using a macro for links now to raise
errors on broken links when they build.

## New Features

- Persistent Containers: Configs with `docker.persistent_name` value will re-use the same container repeatedly.
  - `abyss docker gc` will now ignore containers marked as persistent by metadata
- Secrets Filter Plugin: Created a sample plugin that cancels sessions when a secret is accessed via tools or prompts.

## Bug Fixes

- `abyss oneshot` now loads plugins, as it should
