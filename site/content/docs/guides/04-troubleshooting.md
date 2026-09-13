---
title: "Troubleshooting"
description: "Guides on how to find issues with abyss, diagnose configurations and more."
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

Abyss can fail for a wide variety of reasons due to its surface area.  This goal aims to help you find the root of your issue.

If you think Abyss is doing something wrong, do not hesitate to leave a bug report on [the issues page](https://github.com/SethCurry/abyss/issues).

## Troubleshooting Steps

The first place to start with troubleshooting is `abyss oneshot`, because it dumps logs directly to stderr in a way that's easy to see.
I would suggest running a test prompt via `abyss oneshot` because you can see all of the output, which makes it easier to
test changes than having to restart Zed (which doesn't love unstable ACP agents).

My go to is `abyss oneshot -f $configPath "What is the capital of France?"` because the response tends to be short, and I'm
not testing whether the agent is clever or not.

I have tried to go overboard with logging, so the issue should hopefully be clear to you, or at least to me if not that.
If you have issues, please open a report on [the issues page](https://github.com/SethCurry/abyss/issues).  I'll take a look
and either help you fix it or make a bug fix to fix issues with abyss itself.

## Log File Locations

### Host-Side Proxy Logs

This is where most of the issues are likely to occur.  The host-side proxy handles all of the connection with your ACP client,
as well as all of the Docker interactions (starting containers, copying files, etc).  If you have an issue, it's almost definitely
here.

The host-side proxy (the one that runs outside of the container) is connected directly to your ACP client via stdio.
It logs to both stderr and files, so you can either see the logs directly in your client (Ctrl+Shift+P -> "dev: open ACP logs" in Zed)
or you can find them in timestamped log files at `~/.local/var/abyss/log`.

### Container-Side Proxy Logs

This is the abyss process that runs inside the container.  It likewise logs to both a file and stderr, so you can see its output
via `docker logs $container_id` or via getting a shell in the container and checking `/root/.local/var/abyss/log`
