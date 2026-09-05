---
title: "Getting Started"
description: "Your first steps in abyss"
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
weight: 2
toc: true
params:
  math: false # enable mathematical rendering
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

Before you can use abyss, you'll need a few things:

- Docker installed, with permissions to run containers
- Authentication for your agent (API keys, or whatever your LLM requires) configured in Pi.

## Grab A Release

Releases are available on the [release page](https://github.com/SethCurry/abyss/releases).

In the future when there are fewer features that need attention, I'll get it packaged up in
Homebrew/apt/dnf/AUR/etc.

## Creating Your Configuration

The most basic, usable Pi configuration looks like this:

```yaml
docker:
  # This is a pre-built image with Abyss and Pi pre-installed
  image: "ghcr.io/sethcurry/abyss-pi:latest"

  # Abyss speaks ACP, so we invoke Pi via pi-acp
  agent_command:
    - pi-acp

host_mounts:
  # This mounts the current directory inside the container.
  # It doesn't take a target path because ACP sends a `cwd`
  # so we need the paths to line up between the host and
  # container
  - source: "./"

  # This mounts your .pi directory so that the agent has
  # API keys for your LLM, prior sessions, plugins, etc.
  # 
  # You can put this in the `copy_files` config section
  # which will prevent the LLM from modifying your config files
  # but your sessions won't end up on your host PC so the session
  # files are lost when the container gets deleted.
  - source: "~/.pi"
    destination: "/root/.pi"
```

Paste that into a file wherever you would like.  I like to store mine at the root of my code repository in a file named `abyss.yaml`,
or in an `agents` directory if I have multiple configs for the repo.

You can put it anywhere you would like, though.

## Connecting Your ACP Client

In Zed, open your settings file.  You can either click the hamburger on the top left and use the main Zed menu to click the
"Open Settings File" option.  Alternately, you can use Ctrl+Shift+P to open the command palette and use the "zed: open settings file"
command.

In the config file it opens, paste a section like this:

```json
{
  "agent_servers": {
    "abyss": {
      "type": "custom",
      "command": "abyss",
      "args": [
        "client",
        "-f",
        "/path/to/your/config/file/from/the/last/step.yaml",
      ]
    }
  }
}
```

If the `agent_servers` key already exists, add the `abyss` part to it.  If not, you can copy the entire agent_servers key.

If you would like Zed to call the agent something else, change the name of the JSON key from `"abyss"` to whatever you would like.
(don't change the value of the `"command"` key, change the `"abyss"` above that)

This is what an example of the config looks like:

![Sample Zed Config](/images/zed-config.png)

## Usage

This section is really easy: do what you were doing before.  Abyss tries to stay out of your way, so there are no additional
requirements or restrictions.  Anything that worked before will continue to work.

{{< admonition type="note" title="pi-acp Limitations" >}}
If you come from using Pi via the CLI and are new to using Pi via ACP, pi-acp does have some limitations versus
using it in a terminal.  Those are unfortunately not things that Abyss can fix, that I can think of.
{{< /admonition >}}
