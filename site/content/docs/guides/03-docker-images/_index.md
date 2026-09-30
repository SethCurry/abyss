---
title: "Docker Images"
description: "Build your own Docker images to use for running your agents."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2026-09-22T16:04:48+02:00
draft: false
weight: 3
toc: true
params:
  math: false # enable mathematical rendering
  seo:
    title: "" # custom title (optional)
    description: "" # custom description (recommended)
    canonical: "" # custom canonical URL (optional)
    robots: "" # custom robot tags (optional)
---

The ready-made images on the [Docker Images](../reference/docker-images.md) list are great for
getting started, but eventually you might want something they don't ship. Maybe your agent needs a
language runtime they don't include, a specific version of a tool, or a private config file baked
in. The good news is that abyss is completely happy running inside an image *you* build yourself.

There are two ways to go about it, and we'll walk through both step by step:

- **Extend an existing image** — start from one of the images I already publish and add what you
  need on top. This is the easier path, and it's the one most people should take.
- **Build an entirely custom image** — start from a plain base (like Ubuntu) and install abyss
  yourself. You'd pick this if you want full control over what's inside, or you're targeting an
  unusual platform like ARM.

Don't worry if you've never written a Dockerfile before — we'll explain each line as we go, and by
the end you'll have a working image you can point abyss at.

## Before You Begin

You'll need a couple of things ready before we start.

### 1. Docker

You're going to be building Docker images, so you need Docker installed and running. If you already
followed the [Getting Started](02-getting-started.md) guide, you have this. If not, the quickest
check is to run:

```bash
docker run hello-world
```

If you see a friendly "Hello from Docker!" message, you're set. Otherwise, follow Docker's
[installation guide](https://docs.docker.com/get-started/get-docker/) and come back when it works.

> Make sure you can run `docker` commands *without* typing `sudo` first. If `sudo` is required, add
> your user to the `docker` group — the Docker installation guide explains how.

### 2. A Place to Work

Create an empty folder somewhere on your computer. This is where we'll put the small text file
(called a *Dockerfile*) that describes your image. It can live anywhere; a common spot is a new
folder inside your project:

```bash
mkdir my-agent-image
cd my-agent-image
```

### 3. The Two Rules

No matter which path you take, abyss only asks two things of the image it runs in:

1. **The `abyss` program must be at `/usr/local/bin/abyss`.** Abyss launches itself inside the
   container, so it needs to know exactly where to find its own binary.
2. **`bash` must be installed.** Abyss uses bash to run your setup scripts and the agent's startup
   command.

That's it. Everything else — the operating system, the architecture, the extra tools — is entirely
up to you. Keep those two rules in mind and you can't go far wrong.

## Option 1: Extending an Existing Image

This is the easy path. You start from an image that already has `abyss` and `bash` in the right
place, and you add whatever *extra* things your agent needs on top. Because you're building on a
known-good foundation, there's very little that can go wrong.

You can find the list of images to build on at [Docker Images](../reference/docker-images.md). The
two most common starting points are:

- `ghcr.io/sethcurry/abyss-base:latest` — abyss and bash, nothing else. Great when you want to
  bring your own agent.
- `ghcr.io/sethcurry/abyss-pi:latest` — abyss, bash, and Pi already installed. Great when you use
  Pi but need a few extra tools.

### Step 1: Write the Dockerfile

A *Dockerfile* is just a short text file of instructions. Each line is a step Docker follows to
build your image, one on top of the next. Create a file named `Dockerfile` (no extension) in the
folder you made, and open it in your text editor.

Let's say you use Pi, but your project needs Python 3 to run some tests your agent should be able to
execute. Here's a complete Dockerfile that adds Python on top of the Pi image:

```dockerfile
# Start from the official Pi image. This already has abyss, bash,
# and Pi installed, so we only need to add Python.
FROM ghcr.io/sethcurry/abyss-pi:latest

# Update the package list and install Python 3. We combine the
# commands into one RUN line so the image stays small and tidy.
RUN apt-get update && apt-get install -y python3 python3-pip

# Clean up after the install so we don't ship leftover files.
RUN apt-get clean autoclean \
    && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/

# (Optional) Install any Python packages your agent should be able
# to use.
RUN pip3 install --no-cache-dir pytest ruff
```

Let's walk through what each line does.

- **`FROM ghcr.io/sethcurry/abyss-pi:latest`** tells Docker, "start building from this image."
  Everything we add goes on top of it. The base image already satisfies both of abyss' rules
  (`abyss` is at `/usr/local/bin/abyss` and `bash` is installed), so we don't have to think about
  them.

- **`RUN apt-get update && apt-get install -y python3 python3-pip`** runs commands *inside* the
  image while it's being built. Here we use Ubuntu's package manager (`apt-get`) to install Python
  3 and `pip`. The `-y` flag answers "yes" automatically so the build doesn't pause to ask you
  questions.

- **`RUN apt-get clean ...`** tidies up the package caches. This isn't required, but it keeps your
  image smaller, which means faster downloads and less disk usage.

- **`RUN pip3 install ...`** installs a couple of Python tools your agent might want to use. You can
  list as many as you like, or leave this line out entirely.

### Step 2: Build the Image

Now ask Docker to build the image from your Dockerfile. Run this command from inside the folder
that contains your `Dockerfile`:

```bash
docker build -t my-abyss-pi .
```

Let's unpack that:

- **`docker build`** is the command that reads a Dockerfile and produces an image.
- **`-t my-abyss-pi`** gives the finished image a name (a "tag"). You can call it whatever you like;
  `my-abyss-pi` is just an example. This is the name you'll tell abyss to use later.
- **`.`** (the dot at the end) means "build from the Dockerfile in the current folder." Don't forget
  it!

The first build takes a minute or two because Docker has to download the base image and run the
installs. You'll see a lot of text scroll by — that's normal. When it finishes, you should see a
line like `Successfully tagged my-abyss-pi:latest`.

### Step 3: Test the Image

Before pointing abyss at your new image, it's worth a quick sanity check. Start a container from it
and look around:

```bash
docker run --rm -it my-abyss-pi bash
```

That drops you into a shell *inside* the image. Try a few things:

```bash
which abyss       # should print /usr/local/bin/abyss
which bash        # should print /usr/bin/bash
python3 --version # should print the Python version you installed
```

If all three work, your image is ready. Type `exit` to leave the container.

### Step 4: Point Abyss at Your Image

Now update your abyss configuration to use your new image instead of the pre-built one. Open your
config file and change the `image` line:

```yaml
docker:
  # Use the image you just built instead of the pre-built one.
  image: "my-abyss-pi:latest"

  agent_command:
    - pi-acp

  host_mounts:
    - source: "./"
    - source: "~/.pi"
      destination: "/root/.pi"
```

> **Heads up:** because this image only exists on *your* computer, you need to tell abyss not to try
> pulling it from the internet. Set `image_pull_policy: "Never"` (or `"IfNotPresent"`) in your
> config, otherwise Docker will look for it online and fail. See the
> [Configuration reference](../reference/configuration.md) for the full details.

That's it — you now have a custom image running your agent. Start abyss the same way you normally
would and it'll use your image instead of the published one.

### A Few More Examples

The pattern is always the same: start `FROM` one of my images, then `RUN` whatever you need. Here
are a couple more examples to give you ideas.

**Adding Node.js to the base image** (for an agent that needs JavaScript tooling):

```dockerfile
FROM ghcr.io/sethcurry/abyss-base:latest

RUN apt-get update && apt-get install -y curl ca-certificates \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && apt-get clean autoclean && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/

# Install your agent. Here we use pi-acp from npm.
RUN npm install -g @earendil-works/pi-acp
```

**Copying a config file into the image** (handy for baking in settings that never change):

```dockerfile
FROM ghcr.io/sethcurry/abyss-pi:latest

# Copy a file from your computer into the image. The file must
# exist in the same folder as your Dockerfile when you build.
COPY my-tool.toml /etc/my-tool/my-tool.toml

RUN apt-get update && apt-get install -y my-tool \
    && apt-get clean autoclean && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/
```

**Switching to a specific version of the base image** (so your image doesn't change unexpectedly
when I publish a new `latest`):

```dockerfile
FROM ghcr.io/sethcurry/abyss-pi:0.4.2

RUN apt-get update && apt-get install -y ripgrep
```

> **A note on `latest`:** the `:latest` tag moves whenever a new version is published. That means
> rebuilding your image might pull in a newer base than you expect. If you want builds to be
> reproducible, pin to a specific version tag like `:0.4.2` instead.

### Important: Don't Break the Two Rules

When you extend an image you start off on the right side of abyss' two rules, but it's possible to
accidentally undo them:

- **Don't move the `abyss` binary.** Leave it at `/usr/local/bin/abyss`. If you need to put
  something else there, add it to the `PATH` instead.
- **Don't uninstall `bash`.** Abyss needs it to run your setup scripts and agent command.

If you stick to *adding* things (which is what `RUN apt-get install ...` and `COPY` do), you'll be
fine.

## Option 2: Building an Entirely Custom Image

If extending one of my images doesn't fit — maybe you want a different base operating system, a
minimal image size, or you're targeting a platform I don't publish images for — you can build an
image from scratch. This is more work, but it gives you complete control.

You still only need to satisfy the two rules from earlier:

1. The `abyss` binary is available at `/usr/local/bin/abyss`.
2. `bash` is installed.

That means your two jobs are: install bash (usually easy, via the package manager), and get the
abyss binary into the image (which means either compiling it from source or copying in a binary you
downloaded).

### Step 1: Get the Abyss Binary

The cleanest way is to compile abyss from source right inside the image. This guarantees the binary
matches the image's operating system and processor. Abyss is written in [Go](https://go.dev/), so
you need a Go toolchain to build it.

If you'd rather skip compiling, you can download a pre-built binary from the
[release page](https://github.com/SethCurry/abyss/releases) and `COPY` it into the image instead.
We'll show the compile-from-source approach below since it's the most flexible.

### Step 2: Write the Dockerfile

Here's a complete example that builds abyss from source on top of a plain Ubuntu image. Create a
file named `Dockerfile` in your working folder:

```dockerfile
# syntax=docker/dockerfile:1

# ---- Build stage ----
# We use a separate stage to compile abyss so the Go toolchain
# doesn't end up in the final image. This keeps it small.
FROM ubuntu:latest AS build

WORKDIR /src

# Install the Go toolchain and CA certificates (needed to download
# Go modules over HTTPS).
RUN apt-get update && apt-get install -y golang ca-certificates

# Copy the abyss source code into the image. Run this command from
# the root of the abyss git repository so the source files are
# available.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile the binary. CGO_ENABLED=0 produces a static binary that
# runs anywhere, with no extra libraries required.
RUN CGO_ENABLED=0 go build -o /out/abyss ./cmd/abyss

# ---- Runtime stage ----
# This is the image that actually runs. It only carries what it
# needs: bash and the abyss binary.
FROM ubuntu:latest

# Install bash (Ubuntu already has it, but we make it explicit) and
# clean up afterwards.
RUN apt-get update && apt-get install -y bash ca-certificates \
    && apt-get clean autoclean && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/

# Copy the abyss binary we built in the first stage into the exact
# location abyss expects.
COPY --from=build /out/abyss /usr/local/bin/abyss
```

There's a lot going on here, so let's break it down.

- **Two `FROM` lines?** This is called a *multi-stage build*. The first stage (`AS build`) compiles
  abyss and produces the binary. The second stage is the image you actually keep — it only copies
  in the finished binary, so the big Go toolchain stays behind. The result is a much smaller image.

- **`WORKDIR /src`** sets the current folder for the following commands, like `cd` in a shell.

- **`COPY go.mod go.sum ./` then `RUN go mod download`** copies just the dependency list first and
  downloads dependencies. This is a caching trick: if only your source code changes (not your
  dependencies), Docker reuses the downloaded modules and the build is much faster.

- **`COPY . .`** copies the rest of the abyss source into the image. Run `docker build` from inside
  the abyss git checkout so these files exist.

- **`RUN CGO_ENABLED=0 go build -o /out/abyss ./cmd/abyss`** compiles the abyss binary. The
  `CGO_ENABLED=0` part makes a *static* binary — it doesn't depend on any system libraries, so it
  runs on any Linux system, which is exactly what you want inside a container.

- **`COPY --from=build /out/abyss /usr/local/bin/abyss`** is the magic line that satisfies rule #1.
  It reaches back into the first stage, grabs the binary, and drops it at the path abyss expects.

> **No source code handy?** If you don't have the abyss source checked out, swap the build stage for
> a download instead:
>
> ```dockerfile
> FROM ubuntu:latest
> RUN apt-get update && apt-get install -y bash ca-certificates curl \
>     && curl -fsSL -o /tmp/abyss.tar.gz \
>        https://github.com/SethCurry/abyss/releases/latest/download/abyss_Linux_x86_64.tar.gz \
>     && tar -xzf /tmp/abyss.tar.gz -C /tmp \
>     && mv /tmp/abyss /usr/local/bin/abyss \
>     && rm -rf /tmp/abyss* \
>     && apt-get clean autoclean && apt-get autoremove --yes \
>     && rm -rf /var/lib/{apt,dpkg,cache,log}/
> ```
>
> Pick the archive that matches your image's operating system and processor (`Linux_x86_64` above is
> for Intel/AMD Linux; use `Linux_arm64` for ARM).

### Step 3: Build the Image

If you're compiling from source, run the build from inside the abyss git checkout so the source
files are available to `COPY`:

```bash
cd /path/to/abyss
docker build -t my-abyss-custom -f /path/to/your/Dockerfile .
```

- **`-t my-abyss-custom`** names the image. Pick whatever you like.
- **`-f /path/to/your/Dockerfile`** points at your Dockerfile if it isn't in the current folder.
- **`.`** is the *build context* — the set of files Docker can see when it runs `COPY`. Building
  from the abyss checkout means `COPY . .` can reach the source code.

If you used the download-from-releases approach instead, you can build from any folder that contains
your Dockerfile:

```bash
docker build -t my-abyss-custom .
```

The build takes a few minutes the first time (compiling Go from scratch is the slow part).
Subsequent builds reuse the cached layers and are much faster.

### Step 4: Verify the Two Rules

Just like before, start a shell inside your new image and confirm the two rules hold:

```bash
docker run --rm -it my-abyss-custom bash
```

Then, inside the container:

```bash
which abyss        # must print /usr/local/bin/abyss
which bash         # must print /usr/bin/bash
abyss --version    # should print a version number
```

If `which abyss` comes back empty, the binary didn't land at the right path — double-check your
`COPY` line. If `which bash` is empty, you forgot to install bash. Fix the Dockerfile, rebuild, and
try again. Type `exit` when you're done.

### Step 5: Add Your Agent

So far the image has abyss and bash, but no agent. Add whatever your agent needs with more `RUN` and
`COPY` lines. For example, to install Pi:

```dockerfile
# ...after the runtime stage above...

RUN apt-get update && apt-get install -y nodejs npm ca-certificates \
    && npm install -g @earendil-works/pi-coding-agent pi-acp \
    && apt-get clean autoclean && apt-get autoremove --yes \
    && rm -rf /var/lib/{apt,dpkg,cache,log}/

# Make pi and pi-acp easy to find on the PATH.
RUN ln -s /usr/local/bin/pi /usr/bin/pi \
    && ln -s /usr/local/bin/pi-acp /usr/bin/pi-acp

ENV PATH=/usr/local/bin:/usr/bin:/usr/sbin
```

### Step 6: Point Abyss at Your Image

Finally, update your config file exactly like in Option 1. Set the image name to the one you built,
and remember to set `image_pull_policy` so abyss doesn't try to download it:

```yaml
docker:
  image: "my-abyss-custom:latest"
  image_pull_policy: "Never"

  agent_command:
    - pi-acp

  host_mounts:
    - source: "./"
    - source: "~/.pi"
      destination: "/root/.pi"
```

Fire up abyss and your fully custom image is now running your agent.

## Pushing Your Image Somewhere (Optional)

So far your image lives only on *your* computer. That's perfectly fine for personal use. But if you
want to use the same image on another machine, or share it with your team, you need to publish it to
a *registry*. A registry is just a place that stores images, the same way a code repository stores
source code.

Common choices are [Docker Hub](https://hub.docker.com/), the
[GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry),
or a self-hosted registry. The exact steps vary by registry, but the overall flow is always:

1. **Log in to the registry** from your terminal:
   ```bash
   docker login ghcr.io
   ```
   (Use `docker login` by itself for Docker Hub.) You'll be asked for your username and a token.

2. **Tag your image with the registry's address.** The name has to include where it's going. For
   example, to push to the GitHub Container Registry under the username `alice`:
   ```bash
   docker tag my-abyss-pi ghcr.io/alice/my-abyss-pi:latest
   ```

3. **Push it up:**
   ```bash
   docker push ghcr.io/alice/my-abyss-pi:latest
   ```

Once it's published, anyone (or just you, depending on the visibility you set) can pull it by name.
Point abyss at the full name and set `image_pull_policy: "Always"` (or `"IfNotPresent"`) so abyss
fetches it from the registry:

```yaml
docker:
  image: "ghcr.io/alice/my-abyss-pi:latest"
  image_pull_policy: "Always"
```

## Which Option Should I Pick?

Still not sure which path to take? Here's a quick rule of thumb:

- **Extend an existing image** if one of my published images is *almost* what you want and you just
  need to add a few tools or files. It's faster, simpler, and you automatically pick up fixes when
  the base image is updated.
- **Build a custom image** if you need a different base operating system, a very small image, a
  platform I don't publish for, or you simply want full control over every byte inside.

Both paths produce an image abyss is happy to run. When in doubt, start with Option 1 — you can
always switch to Option 2 later.

## Where to Go Next

- [Docker Images](../reference/docker-images.md) — the full list of images you can build on.
- [Configuration](../reference/configuration.md) — every config option, including `image_pull_policy`
  and `setup_scripts` for run-time customization.
- [Troubleshooting](04-troubleshooting.md) — what to do when your image won't build or abyss can't
  start it.
