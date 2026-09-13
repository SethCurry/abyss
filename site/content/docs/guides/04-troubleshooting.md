---
title: "Custom Docker Images"
description: "Build your own Docker images to use for running your agents."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
weight: 5
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
