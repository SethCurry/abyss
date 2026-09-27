# Overview

This is an agent image for abyss based on [codex-acp](https://github.com/agentclientprotocol/codex-acp).

Here is an example config using it:

```yaml
docker:
  # The name of the image to use
  image: "ghcr.io/sethcurry/abyss-codex:latest"
  agent_command: 
    - codex-acp
```
