# Overview

This is an agent image for abyss based on [claude-agent-acp](https://github.com/agentclientprotocol/claude-agent-acp).

Here is an example config using it:

```yaml
docker:
  # The name of the image to use
  image: "ghcr.io/sethcurry/abyss-claude-agent:latest"
  agent_command: 
    - claude-agent-acp
```
