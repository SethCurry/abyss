---
title: "Configuration"
description: "Learn how to configure Ares agents using the agent configuration file."
summary: ""
date: 2023-09-07T16:13:18+02:00
lastmod: 2023-09-07T16:13:18+02:00
draft: false
weight: 1
toc: true
params:
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

Abyss is configured using an agent configuration file.  It controls
all of the aspects of your agent as well as the environment in which it runs.

## Full Example Configuration

```yaml
docker:
  # The name of the image to use
  image: "abyss:latest"

  # A list of host mounts to bind mount into the container
  # You can specify relative paths, but they will be mounted
  # inside the container at the absolute path on your host machine
  host_mounts:
    - "./src"

  # The command to run inside the container
  agent_command:
    - npx
    - "@earendil-works/pi-coding-agent"

# Setup scripts are shell scripts that are run before the agent starts
# They allow executing custom steps at run-time, for things that aren't
# suitable for building into the Docker image itself.
setup_scripts:
  - type: "inline"
    source: "cp /something /something-else"
  - type: "file"
    source: "./scripts/setup.sh"
copy_files:
  - type: "inline"
    source: "nameserver 8.8.8.8"
    target: "/etc/resolv.conf"
  - type: "path"
    source: "./config/some-config.yml"
    target: "/etc/some-config.yml"
```
