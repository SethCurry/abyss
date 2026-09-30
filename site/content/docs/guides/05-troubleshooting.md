---
title: "Troubleshooting"
description: "Guides on how to find issues with abyss, diagnose configurations and more."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
weight: 5
toc: true
---

Abyss has a lot of moving parts — it talks to Docker, starts a container, wires up an encrypted
tunnel, launches your agent, and then translates everything your editor says into something the
agent understands. That's a big surface area, which means there are a lot of places something can
go sideways.

This guide is here to help you find *which* part is misbehaving, and to give you the tools to fix
the common ones yourself. Take a breath — most issues are small and easy to spot once you know where
to look.

If you read through this and you're still stuck — or you think abyss itself is doing something wrong
— please don't hesitate to leave a bug report on [the issues page](https://github.com/SethCurry/abyss/issues).
I'll take a look and either help you fix it or push a fix to abyss itself.

## The Number One Tool: `abyss oneshot`

When something breaks, the temptation is to keep restarting your editor and hoping it works this
time. Don't. Your editor (Zed included) is great for *using* an agent, but it's awkward for
*debugging* one — it hides abyss's logs behind a menu, and it doesn't love agents that crash on
startup.

Instead, reach for `abyss oneshot`. It runs a single prompt against your config and dumps every log
line straight to your terminal where you can see them. It's the fastest way to answer "is it my
config, or is it something else?"

My go-to test prompt is:

```bash
abyss oneshot -f /path/to/your/config.yaml "What is the capital of France?"
```

I use that exact prompt on purpose. The answer is short, so I'm not waiting around, and I'm not
trying to test whether the agent is clever — I just want to know if the plumbing works end to end.
If you get back "Paris," your whole stack is healthy and the problem is somewhere else (probably
your editor config).

If it errors out, the logs printed to your terminal will usually point you right at the culprit.
The rest of this guide walks through the most common culprits and what to do about them.

> **Tip:** `oneshot` is also how you test changes to your config without restarting your editor.
> Tweak the YAML, re-run the command above, and you'll know within seconds whether your change
> helped.

## Where Abyss Keeps Its Logs

Abyss logs to *both* your terminal (stderr) and to timestamped files on disk. The terminal output
is great for a quick look; the files are great for when something blew up an hour ago and you want
to read back through it.

### Host-Side Proxy Logs

This is where most issues live. The host-side proxy is the `abyss` process that runs on *your*
machine, outside the container. It's the one your editor talks to, and it's the one that does all
the Docker work — pulling images, starting containers, copying files in, tearing things down. If
something is broken, it's almost definitely here.

It logs to both stderr and files, so you have two ways to read it:

- **In your editor.** Your ACP client captures the proxy's stderr. In Zed, open the command palette
  (`Ctrl+Shift+P` or `Cmd+Shift+P` on macOS) and run **"dev: open ACP logs"**. You'll see abyss's
  output mixed in with ACP traffic.
- **On disk.** Timestamped log files live at `~/.local/var/abyss/log`. Each run gets its own file
  named with the time it started, like `2025-01-30T14-22-07Z.log`. Abyss automatically cleans up
  old files and keeps the most recent 10, so the folder won't grow forever.

### Container-Side Proxy Logs

The container-side proxy is the `abyss` process running *inside* the container. It's the one that
actually launches your agent and translates ACP file and terminal requests into actions inside the
sandbox. It logs to both a file and stderr too, so you have two ways in:

- **From the host**, with `docker logs <container-id>`. This is the easiest way; you don't need to
  get a shell in the container.
- **From inside the container.** If you'd rather look directly, open a shell in the container and
  check `/root/.local/var/abyss/log`.

> **Don't know the container ID?** Run `abyss docker ps` (covered below) to list every abyss
> container that's currently running, along with its ID.

### Turning the Log Volume Up (or Down)

By default abyss logs at the `debug` level, which is fairly chatty and is what you want while
troubleshooting. If you ever need *more* detail (there usually isn't much, but it's there), you can
drop to `trace`:

```bash
abyss oneshot -l trace -f /path/to/your/config.yaml "What is the capital of France?"
```

You can also set the level with the `ABYSS_LOG_LEVEL` environment variable, which is handy when
abyss is being launched by your editor and you can't easily add flags:

```bash
export ABYSS_LOG_LEVEL=trace
```

Valid levels are `trace`, `debug`, `info`, `warn`, `error`, `fatal`, and `disabled`. If you're
confident things work and just want quiet logs, `info` is a good everyday level.

## Managing Containers

Abyss starts a fresh container for each session and stops it when the session ends. Normally you
never have to think about this. But if a session crashed, or you killed abyss mid-run, you can end
up with containers that are still chugging along in the background. Abyss ships a couple of commands
for exactly this situation.

### See What's Running

```bash
abyss docker ps
```

This lists every running container that abyss started (it finds them by their `abyss` label). You'll
get back the container ID and its name, which is everything you need to inspect or stop it.

### Clean Up Everything

```bash
abyss docker gc
```

This stops *every* running abyss container. It's the "reset button" when things have gotten into a
weird state — for example, if a leftover container is holding onto a port or a bind mount and a new
session won't start because of it. Run it, then try your session again.

> **These two are safe to run any time.** They only touch containers that abyss itself created, and
> they stop containers rather than removing your files. If you're ever unsure what state things are
> in, `abyss docker ps` followed by `abyss docker gc` is a fine first move.

### Inspecting a Container by Hand

Sometimes you want to poke around inside a container that's still running — to check whether a file
got copied in, whether your agent is actually installed, or what a setup script left behind. Get a
shell with:

```bash
docker exec -it <container-id> bash
```

From there you can run `ls`, check `/root/.local/var/abyss/log`, verify your agent command exists
on the `PATH`, and so on. This is a great way to confirm "is the thing I'm mounting actually showing
up where I expect?"

## Common Problems, and How to Fix Them

The sections below cover the issues that come up most often. They're roughly in the order you'd hit
them: Docker first, then the image, then your config, then the agent, then the connection.

### "Failed to connect to Docker"

If `abyss oneshot` greets you with a message about failing to connect to Docker, the fix is almost
always one of two things:

1. **Docker isn't running.** Start the Docker daemon (or Docker Desktop, on macOS/Windows) and try
   again. A quick `docker run hello-world` will tell you whether Docker is alive and that you have
   permission to use it.
2. **You need `sudo` to use Docker.** Abyss can't prompt you for a password, so if `docker` only
   works when you prefix it with `sudo`, abyss will be locked out. Add your user to the `docker`
   group so you can run `docker` without `sudo`. Docker's
   [installation guide](https://docs.docker.com/get-started/get-docker/) walks through this.

### "Failed to pull Docker image"

This means abyss couldn't get the image named in your `docker.image` field. A few common causes:

- **The image name is wrong.** Double-check the spelling and the registry path. The ready-made
  images are listed in [Docker Images](../reference/docker-images.md).
- **You don't have access to the registry.** If you're pulling from a private registry, run
  `docker login` first so your credentials are saved.
- **You're offline, or the registry is down.** Try `docker pull <image>` by hand — if that fails,
  abyss will fail the same way, and the error from `docker` is usually clearer.
- **You set `image_pull_policy: Never` but don't have the image locally.** With `Never`, abyss
  refuses to pull and will only use an image that's already on your machine. Either build/pull the
  image yourself, or switch the policy to `IfNotPresent` (the default) or `Always`.

### The Container Starts but the Agent Doesn't

If abyss brings up the container fine but your agent never responds, the trouble is usually in how
the agent is being launched. Things to check:

- **Is your `agent_command` correct?** It should be the command that starts your agent in its ACP
  mode — for Pi that's `pi-acp`. If you've typo'd it, or pointed at a command that isn't installed
  in the image, the container-side proxy will start but the agent won't.
- **Does the agent actually exist on the container's `PATH`?** Shell in with
  `docker exec -it <container-id> bash` and run `which pi-acp` (or whatever your command is). If it
  comes back empty, the image doesn't have it installed where abyss expects.
- **Are your agent's credentials mounted in?** Most agents need their config directory to talk to
  an LLM. For Pi that's `~/.pi`, and because the container runs as `root` you usually mount it to
  `/root/.pi` (see [Getting Started](02-getting-started.md)). If the mount is missing or pointed at
  the wrong place, the agent starts but can't reach your LLM, which looks a lot like "it's just
  hanging."
- **Does the agent work on its own?** Try running your `agent_command` directly inside the
  container via `docker exec`. If it errors there too, the problem is the agent or its environment,
  not abyss.

### Bind Mounts Don't Show Up (or Show Up in the Wrong Place)

Bind mounts are the most common source of "why can't my agent see my files?" confusion. A couple of
things to keep in mind:

- **A mount with no `destination` lands at the *same* path as on your host.** That's intentional —
  it keeps your editor and your agent agreeing on where files live. But it means a relative `source`
  gets resolved to an absolute path on your machine, and that absolute path has to *exist* on your
  machine for Docker to mount it.
- **Tildes (`~`) only expand on the host, never in the `destination`.** If you write
  `destination: "~/foo"`, the `~` is expanded as *your* user on your machine, not as the container's
  user. That's almost never what you want. Use an absolute path for `destination` — for the
  container's `root` user that usually means starting with `/root/`.

If a mount seems missing, shell into the container and `ls` the path you expected it at. If the
directory is empty or doesn't exist, the mount didn't take, and the host-side logs will usually tell
you why.

### Files You Meant to Copy In Aren't There

`copy_files` runs *before* your agent starts, copying files from your host (or inline strings) into
the container. If something you expected isn't present:

- **Check the `target` path.** Abyss creates missing parent directories with mode `0755`, but the
  path still has to be one you're allowed to write to inside the container.
- **Check the `source` for `type: path` entries.** The path is read from your host, so it has to
  exist *on your machine*, not inside the image.
- **Read the host-side logs.** Every copy is logged, along with any error. If a copy failed, you'll
  see it there rather than having to guess.

### Setup Scripts Hang or Fail

`setup_scripts` run in order, once each, before the agent starts. If your session seems to hang
forever at startup, a setup script is the usual suspect — abyss waits for every script to finish
before launching the agent, so a script that's waiting on input or stuck on a network call will
pause the whole startup.

To narrow it down:

- **Run the script by hand inside the container** (`docker exec -it <container-id> bash`) and see
  where it blocks.
- **Add `set -ex` to the top of bash scripts** so each step is printed and the script stops at the
  first failure. The output shows up in the container logs.
- **Remember startup is serial.** Each script adds to your startup time. If things feel slow but
  still work, consider moving static work into your image instead. See
  [Custom Docker Images](03-custom-docker-images.md).

### TLS / Connection Errors

By default abyss generates a fresh CA and certificates for *every* session and tears them down
afterward — that's the ephemeral mutual-TLS you don't have to think about. The only time you'll see
TLS errors is if something interferes with that handshake, and the most common cause is a leftover
container from a previous crashed session holding onto the port.

The fix is the reset button from earlier:

```bash
abyss docker gc
```

Then try again. If you've deliberately set `websocket.disable_tls: true` in your config, you
shouldn't see TLS errors at all — but remember that turning TLS off is only safe for local
single-user testing, so double-check that was intentional.

### Your Editor Can't See abyss (or Says the Agent Crashed)

If `abyss oneshot` works but your editor can't connect, the problem is in how the editor is
configured, not in abyss itself. For Zed, check that your `agent_servers` block points at the real
`abyss` binary and the real path to your config file:

```json
{
  "agent_servers": {
    "abyss": {
      "type": "custom",
      "command": "abyss",
      "args": ["client", "-f", "/absolute/path/to/your/config.yaml"]
    }
  }
}
```

A couple of things that trip people up:

- **Use an absolute path to the config.** Your editor may launch abyss from a different working
  directory than you expect, so a relative path can resolve to the wrong file.
- **Make sure `abyss` is on the editor's `PATH`.** If you can run `abyss --version` in your terminal
  but the editor can't find it, the editor is probably running with a different environment. Either
  put `abyss` somewhere universal like `/usr/local/bin`, or use the full path in `"command"`.
- **Check the editor's ACP logs.** In Zed that's `Ctrl+Shift+P` → **"dev: open ACP logs"**. If abyss
  is starting and then immediately dying, the host-side logs (on disk at
  `~/.local/var/abyss/log`) will tell you why.

### File or Terminal Requests Go to the Wrong Place

By default abyss *intercepts* ACP file and terminal requests and runs them *inside the container*,
so the agent can't reach out and touch your machine. If you've flipped `acp.tools_on_host.filesystem`
or `acp.tools_on_host.terminal` to `true`, those requests instead get forwarded to your editor and
run on your host.

If reads or writes or shell commands seem to happen in the "wrong" filesystem, check those settings.
Running on the host is a deliberate escape hatch — it's powerful, but it does mean the agent can
affect your real machine, so make sure that's actually what you meant to do. The full details are in
the [Configuration reference](../reference/configuration.md).

## Still Stuck? Filing a Good Bug Report

If you've worked through the above and abyss still isn't behaving, please open an issue on
[the issues page](https://github.com/SethCurry/abyss/issues). A good report helps me help you
quickly, so try to include:

- **The abyss version** (`abyss --version`).
- **Your operating system and architecture** (e.g. macOS on Apple Silicon, Linux on x86_64).
- **The smallest config file that reproduces the problem.** Strip out anything that isn't needed to
  trigger the issue — mounts, plugins, setup scripts. The smaller it is, the faster I can reproduce
  it.
- **The exact command you ran** and the output it printed. If `abyss oneshot` reproduces it, prefer
  that over an editor session — it's much easier for me to run.
- **The relevant log file** from `~/.local/var/abyss/log`.  Only if you are comfortable or edit them;
  debug logs contain the full text of ACP messages.

I've tried to go overboard with logging, so the issue should hopefully be clear to you, or at least
to me if not. I'll take a look and either help you fix it or push a fix to abyss itself.

## Where to Go Next

- [Getting Started](02-getting-started.md) — a clean walkthrough of a working setup, useful as a
  known-good baseline to compare a broken config against.
- [Configuration](../reference/configuration.md) — the full reference for every config option, in
  case a field isn't doing what you expect.
- [Custom Docker Images](03-custom-docker-images.md) — if your troubles trace back to the image
  itself, this covers how to build or extend one.
