---
title: "Custom Docker Images"
description: "Build your own Docker images to use for running your agents."
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

If you want to use a custom image with abyss there are largely 2 options:

- Extending one of the images abyss provides
- Creating an entirely new and custom image

It is easier to extend one of my images than create your own, but neither is hard.

## Extending An Existing Image

You can find a list of images to base yours on [here](../reference/docker-images.md).

The only requirements for extending an image are:

- Do not move the abyss binary from /usr/local/bin/abyss
- Do not uninstall bash

abyss doesn't care about anything else.

E.g. for a custom Pi image, all you need to do is start your Dockerfile with:

```
FROM ghcr.io/sethcurry/abyss:latest

RUN my-script.sh
```

## Creating an entirely custom image

If you want an entirely custom image, you only need to meet 2 requirements:

- The abyss binary is available at /usr/local/bin/abyss
- bash is installed

There are no restrictions other than that, you're free to run it on x86, ARM, whatever you can get it to
compile on.

You can install bash via the system package manager; for `abyss` you will likely need to install a Go toolchain
and compile it from scratch.
