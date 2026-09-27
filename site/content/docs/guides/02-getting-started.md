---
title: "Getting Started"
description: "Install abyss, write your first agent configuration, and connect it to your editor."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2026-09-27T19:51:49+00:00
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

This guide walks you through everything you need to go from "just heard about abyss" to "running
your agent inside a sandbox." By the end, you'll have abyss installed, a working configuration file,
and your agent connected to your editor.

Don't worry if the word *container* or *YAML* sounds intimidating — we'll explain each step as we go.

## Before You Begin

You'll need three things before we start. None of them are complicated, and you probably already
have at least two of them.

### 1. Docker

Abyss runs your agent inside a [Docker](https://www.docker.com/) container, so you need Docker
installed and running on your computer.

If you're not sure whether you have it, open a terminal and type:

```bash
docker run hello-world
```

If Docker is installed and working, you'll see a friendly "Hello from Docker!" message. If instead
you see an error about the command not being found, head over to Docker's
[installation guide](https://docs.docker.com/get-started/get-docker/) and come back when it's ready.

> **One thing to check:** make sure you can run `docker` commands *without* typing `sudo` first.
> If `sudo` is required, you may need to add your user to the `docker` group. The Docker
> installation guide above explains how.

### 2. An ACP Client

Abyss speaks the [Agent Client Protocol](https://agentclientprotocol.com/) (ACP). ACP is just the
common language that your editor and your agent use to talk to each other. Any editor that supports
ACP will work with abyss.

The easiest editor to get started with is [Zed](https://zed.dev/), and that's what we'll use in the
examples below. If you'd rather use a different editor, the concepts are identical — only the
settings screen looks different.

### 3. An Agent and Its Credentials

Abyss doesn't include an agent of its own; it provides a safe place for *your* agent to run. So you
need an agent to put inside it, plus whatever API keys that agent needs to talk to your LLM.

The good news: abyss ships ready-made images for four popular agents, so you most likely won't
have to build anything yourself:

- **[Pi](https://pi.dev/)** — the agent this guide uses for its detailed walkthrough.
- **[Hermes](https://hermes-agent.nousresearch.com/)** — Nous Research's agent.
- **[Codex](https://github.com/openai/codex)** — OpenAI's coding agent.
- **[Claude Agent](https://claude.com/product/claude-code)** — Anthropic's coding agent.

This guide explains the Pi configuration in detail, because every agent's file is built from the
same three pieces — once you've seen one, the others look very familiar. If you use one of the
other three, you'll find a section with your exact configuration in
[Create Your Configuration](#create-your-configuration) below. And if you use a different agent
entirely, the setup is the same idea — you'll just point at your own image and command. (We cover
that in [Custom Docker Images](03-custom-docker-images.md).)

Whatever agent you pick, make sure it's already set up and can talk to your LLM on its own before
you put it in a sandbox. In practice that means the folder where it keeps your API keys and
settings — for Pi, that's your `~/.pi` directory — exists and works.

## Install Abyss

Abyss is a single program called `abyss`. Right now it's distributed as a downloadable binary from
the [release page](https://github.com/SethCurry/abyss/releases).

1. Go to the [release page](https://github.com/SethCurry/abyss/releases) and download the archive
   for your operating system and processor. For example, if you're on an Apple Silicon Mac, grab
   the `Darwin_arm64` archive; on an Intel/AMD Linux machine, grab `Linux_x86_64`.
2. Unpack the archive. Inside you'll find a single file named `abyss`.
3. Move that file somewhere your terminal can always find it. On Linux and macOS, a good home is
   `/usr/local/bin/abyss`. On Windows, put it in a folder that's listed in your `PATH`.

On Linux or macOS, steps 2 and 3 usually look like this (adjust the archive name to the one you
downloaded):

```bash
tar -xzf abyss_Linux_x86_64.tar.gz
sudo mv abyss /usr/local/bin/abyss
```

> **In the future**, once the feature set settles down, abyss will be packaged for
> Homebrew, apt, dnf, AUR, and friends so you can skip this manual step. For now, it's a simple
> download-and-move.

### Check That It Worked

Confirm the installation by asking abyss for its version:

```bash
abyss --version
```

You should see a version number and a date. If you get "command not found" instead, your terminal
can't find the `abyss` file — double-check the folder you moved it into is on your `PATH`.

## Create Your Configuration

Abyss is driven by one small YAML file. That file describes three things:

1. **Which container image** to run your agent in.
2. **What command** starts your agent.
3. **Which folders** on your computer the agent should be able to see.

Every agent's file is built from those same three pieces — only the values change. Below there's
a section for each of the ready-made images: Pi, Hermes, Codex, and Claude Agent. Start with the
Pi section even if you use a different agent, since it explains how everything fits together. Then
copy your own agent's configuration and meet us at [Where to Put the File](#where-to-put-the-file).

### Using Pi

Here's a complete, working configuration for Pi:

```yaml
docker:
  # A pre-built image with abyss and Pi already installed.
  image: "ghcr.io/sethcurry/abyss-pi:latest"

  # Abyss speaks ACP, so we start Pi in its ACP mode.
  agent_command:
    - pi-acp

  # Folders from your computer that get shared into the container.
  host_mounts:
    # Share the current directory. There's no destination path on
    # purpose: ACP sends your "current working directory", so the path
    # inside the container must match the path on your computer for
    # everything to line up.
    - source: "./"

    # Share your Pi settings so the agent has your API keys, prior
    # sessions, and plugins. Because the container runs as the `root`
    # user, we point it at `/root/.pi` instead of your own home folder.
    - source: "~/.pi"
      destination: "/root/.pi"
```

Let's unpack the important bits.

- **`image`** is the Docker image that provides the environment your agent runs in. We're using
  `ghcr.io/sethcurry/abyss-pi:latest`, which comes with both abyss and Pi already installed, so
  there's nothing extra to set up. You can find the full list of ready-made images in
  [Docker Images](../reference/docker-images.md).

- **`agent_command`** is the command that starts your agent once the container is up. For Pi, that's
  simply `pi-acp`.

- **`host_mounts`** is how you share folders between your computer and the container. This is the
  piece that lets your agent read and edit the files you actually care about.
  - The first mount shares your current directory (`./`). Because abyss mounts it at the *same*
    path inside the container, your editor and your agent always agree on where files live — no
    confusing path translation to memorize.
  - The second mount shares your `~/.pi` folder, which holds your API keys. We give it an explicit
    `destination` because the container's `root` user keeps its settings in `/root/.pi`, not in
    your own home directory.

### Using Hermes

[Hermes](https://hermes-agent.nousresearch.com/) follows the same recipe as Pi — its own image,
own command — with one important difference, so give this section a careful read.

The `abyss-hermes` image comes with a complete copy of Hermes already installed, and that copy
lives in the container's `/root/.hermes` folder. That's also where Hermes looks for your settings
when it runs. Your computer has its own `~/.hermes` folder too: that's where Hermes keeps your API
keys (in `~/.hermes/.env`) and your preferences (in `~/.hermes/config.yaml`).

Here's the catch. With Pi, we shared the entire `~/.pi` folder with the container. If we did the
same with `~/.hermes`, your folder would cover up the pre-installed copy of Hermes that the image
needs in order to run — and if your computer isn't the same kind of machine as the container (a
Mac, for instance), the agent won't start at all. So instead of sharing the whole folder, we copy
in just the two files Hermes needs:

```yaml
docker:
  # A pre-built image with abyss and Hermes already installed.
  image: "ghcr.io/sethcurry/abyss-hermes:latest"

  # Hermes ships with its own ACP mode.
  agent_command:
    - hermes-acp

  host_mounts:
    # Share the current directory, just like with Pi.
    - source: "./"

# Copy your Hermes credentials and settings into the container.
# Both files live in ~/.hermes on your computer.
copy_files:
  - type: "path"
    source: "~/.hermes/.env"
    target: "/root/.hermes/.env"
  - type: "path"
    source: "~/.hermes/config.yaml"
    target: "/root/.hermes/config.yaml"
```

Because we're *copying* files in rather than sharing a live folder, anything Hermes saves while
it's inside the sandbox — new sessions and memories, for example — stays inside the sandbox and
disappears when it shuts down, instead of showing up in `~/.hermes` on your computer.

### Using Codex

If OpenAI's [Codex](https://github.com/openai/codex) is your agent, good news: its configuration
is nearly identical to Pi's. The image comes with Codex and its
[ACP adapter](https://github.com/agentclientprotocol/codex-acp) already installed, so the only
thing you need to bring along is your login:

```yaml
docker:
  # A pre-built image with abyss, Codex, and codex-acp already installed.
  image: "ghcr.io/sethcurry/abyss-codex:latest"

  # codex-acp is a small adapter that lets Codex speak ACP.
  agent_command:
    - codex-acp

  host_mounts:
    # Share the current directory, just like with Pi.
    - source: "./"

    # Share your Codex login and settings. Codex keeps these in
    # ~/.codex on your computer, and the container's root user looks
    # for them in /root/.codex.
    - source: "~/.codex"
      destination: "/root/.codex"
```

For this to work, you'll need to have signed in to Codex on your computer at least once, so that
your `~/.codex` folder exists and contains your login.

### Using Claude Agent

Anthropic's [Claude Agent](https://claude.com/product/claude-code) follows the exact same recipe as
Codex, with its own image, adapter, and settings folder:

```yaml
docker:
  # A pre-built image with abyss, Claude Agent, and claude-agent-acp
  # already installed.
  image: "ghcr.io/sethcurry/abyss-claude-agent:latest"

  # claude-agent-acp is a small adapter that lets Claude Agent speak ACP.
  agent_command:
    - claude-agent-acp

  host_mounts:
    # Share the current directory, just like with Pi.
    - source: "./"

    # Share your Claude login and settings. Claude keeps these in
    # ~/.claude on your computer, and the container's root user looks
    # for them in /root/.claude.
    - source: "~/.claude"
      destination: "/root/.claude"
```

Just like with Codex, this assumes you've already signed in to Claude on your computer, so your
`~/.claude` folder exists and holds your credentials.

### Where to Put the File

Save the file anywhere you like. A few common spots:

- `abyss.yaml` at the root of your project, if the agent is only ever used there.
- An `agents/` directory inside the project, if you keep several agent configs for one repo.

Any location works — you'll tell abyss where the file lives in the next step, so pick whatever feels
tidy to you.

## Connect Your Editor

Now let's point your editor at abyss. We'll use Zed here, but the idea is the same in any ACP client:
you register abyss as an "agent server" and give it the path to your config file.

In Zed, open your settings file. The quickest way is to press `Ctrl+Shift+P` (or `Cmd+Shift+P` on
macOS) to open the command palette, then type **"zed: open settings file"** and press Enter.

In that file, paste a section like this:

```json
{
  "agent_servers": {
    "abyss": {
      "type": "custom",
      "command": "abyss",
      "args": [
        "client",
        "-f",
        "/path/to/your/config/file.yaml"
      ]
    }
  }
}
```

A few notes on this block:

- Replace `/path/to/your/config/file.yaml` with the real path to the file you created in the
  previous step.
- The `"abyss"` key is just the name Zed will show for this agent. You can rename it to anything you
  like (for example `"my-agent"`) without changing anything else. Just leave the `"command"` value
  as `"abyss"`.
- If your settings file already has an `"agent_servers"` section, merge the `"abyss"` entry into it
  rather than adding a second one.

For reference, here's what a finished Zed config looks like:

![Sample Zed Config](/images/zed-config.png)

Once you've saved your settings, Zed should now list abyss as one of its available agents. Select it
and start chatting — abyss will spin up a container in the background, start your agent inside it,
and hand the conversation back to Zed.

## Test It Without Your Editor

If you'd like to confirm everything works before involving your editor, abyss has a built-in
one-shot mode that runs a single prompt and prints the answer. It's a great way to sanity-check a
new config:

```bash
abyss oneshot -f /path/to/your/config/file.yaml "What is the capital of France?"
```

If you get a sensible answer (it's Paris, in case you're wondering), your configuration is solid and
you're ready to go. `oneshot` also prints its logs directly to the terminal, which makes it the
easiest way to see what's going wrong if something isn't working. We lean on it heavily in the
[Troubleshooting](04-troubleshooting.md) guide.

## Using It Day to Day

Here's the good news: once abyss is connected, using it is exactly like using your agent before.
Abyss is designed to stay out of your way. There are no extra commands to memorize and no special
workflow — you just ask your agent to do things, and it works the same as always, except everything
now happens inside a sandbox.

A couple of tips as you settle in:

- **Startup takes a moment.** The first time you use an agent, abyss has to pull the Docker image and
  start a container. Subsequent launches are much faster.
- **Want to run setup steps or install dependencies?** You can add a `setup_scripts` section to your
  config so things run automatically before your agent starts.
- **Prefer the agent *not* to touch your files directly?** Swap `host_mounts` for `copy_files` to
  hand the agent a fresh copy instead of a live mount.

Both of these are covered in the [Configuration reference](../reference/configuration.md), which
lists every option abyss supports.

## Where to Go Next

- [What is abyss?](01-what-is-abyss.md) — the big picture of why abyss exists and how it works.
- [Custom Docker Images](03-custom-docker-images.md) — run a different agent or build your own image.
- [Configuration](../reference/configuration.md) — the full reference for every config option.
- [Troubleshooting](04-troubleshooting.md) — what to do when something doesn't work.

{{< admonition type="note" title="A Note About pi-acp" >}}
If you're coming from using Pi in a terminal, `pi-acp` (the mode abyss uses) has a few limitations
compared to the full terminal experience. Those are limitations of `pi-acp` itself rather than
something abyss can work around — worth keeping in mind if a feature you relied on seems missing.
{{< /admonition >}}
