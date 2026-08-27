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
    # Host mounts without a destination set get mounted inside the
    # container at the same path as on the host.
    - source: "./src"
    
    # You can also set a destination path.
    # NOTE: Tilde (~) is only evaluated on the host.  If you use
    # a ~ in the destination path, it gets evaluated as your user
    # on your PC, not inside the container as the container user.
    - source: "~/.pi"
      destination: "/root/.pi"

  # The command to run inside the container
  agent_command:
    - pi-acp

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

## `docker`

### `image`

`image` sets the Docker image to use.

If you are creating your own image, make sure you have the `abyss` binary installed at `/usr/local/bin/abyss`.
Also make sure that you have `bash` installed; some containers are slim and only include `sh`.

### `agent_command`

`agent_command` sets the command to start the agent.  This is only the command to start the agent, like `pi-acp`.
Abyss handles translating this command to a proper `abyss server` invocation for you.

### `host_mounts`

`host_mounts` is a list of directories you want to bind mount inside the container (i.e. the container has access to
modify the files on your host from inside the container).

Each host mount takes these fields:

| Field | Required | Description |
| ----- | -------- | ----------- |
| `source` | Yes | The path to the directory on your PC.  It can be absolute or relative, and can include `~`. |
| `destination`| No | The path to mount the direcory inside the container.  If not specified, the default behavior is to mount it at the same path as on your PC. |

{{< admonition type="note" title="On Tilde (~) Expansion" >}}
Tildes (~) are _only_ expanded on the host.  If you use a `~` as part of the `destination` path, it will be expanded as if you ran it as your user on the host.  I would heavily encourage using absolute paths for `destination`.

{{< /admonition >}}

## `setup_scripts`

`setup_scripts` are files that are copied into the container and executed before the Abyss server and your agent are started.

There are 2 primary uses for `setup_scripts`:

1. To pull things in that are too dynamic to bake into an image, like a git repository that you want to be always up-to-date
2. To avoid building your own Docker image

Each script will be executed in series.

There are no restrictions on language; you're welcome to write them in Python, Ruby, Perl, bash, whatever, as long as the image
you're using has the interpreter installed (or it was installed by a prior setup_script).

These are the options for a `setup_script`:

| Field | Required | Description |
| ----- | -------- | ----------- |
| `type` | Yes | Specifies the source for the script.  Currently only `file` and `inline` are supported.  If the `type` is `"file"`, then `source` is interpreted as a path to a file.  If the type is `"inline"`, then `source` is assumed to contain a script as a string. |
| `source` | Yes | The source of the file.  If type is "file", this should be a path to a file.  If the type is "inline", then put the script contents directly in this field. |

{{< admonition type="warning" title="Startup Time" >}}
These scripts have to finish executing before the Abyss server or your agent start, so they impose a "startup time" tax.
Abyss doesn't care how long it takes to start, but it may frustrate you if you jam a ton of stuff in.

If startup times get too long, I would suggest looking at moving things into the image you're using if at all possible.
{{< /admonition >}}

## `copy_files`

`copy_files` specifies a set of files or directories that should be copied from the host into the container before Abyss starts.

These are the fields for a `copy_files` block:

| Field | Required | Description |
| ----- | -------- | ----------- |
| `type` | Yes | The type of the source for the file, either `path` or `inline`. |
| `source` | Yes | The source of the file.  If `type` is `path` then `source` should be a path to the file.  If `type` is `inline`, then the contents of the file go directly into the YAML. |
| `target` | Yes | Where to copy the files inside the container.  If any intermediary directories don't exist, they are automatically created with mode `0755`. |
