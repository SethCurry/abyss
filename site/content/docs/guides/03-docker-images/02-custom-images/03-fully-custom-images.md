---
title: "Fully Custom Images"
description: "Documentation on how to build a Docker image that will work with Abyss from scratch."
summary: ""
date: 2023-09-07T16:05:48+02:00
lastmod: 2026-09-22T16:04:48+02:00
draft: false
weight: 3
toc: true
---


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
