---
title: "Extending An Abyss Image"
description: "Documentation on how to extend an Abyss image with your own tools, libraries, etc."
summary: ""
date: 2023-09-07T16:05:48+02:00
lastmod: 2026-09-22T16:04:48+02:00
draft: false
weight: 2
toc: true
---


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
03-fully-custom-image
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
