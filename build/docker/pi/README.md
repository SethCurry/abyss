# Overview

This is an agent image for abyss based on [Pi](https://pi.dev/), using
[pi-acp](https://github.com/svkozak/pi-acp) for the ACP connection.

Here is an example config using it:

```yaml
docker:
  # The name of the image to use
  image: "ghcr.io/sethcurry/abyss-pi:latest"
  agent_command: 
    - pi-acp
```
