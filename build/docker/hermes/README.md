# Overview

This is an agent image for abyss based on [hermes](https://hermes-agent.nousresearch.com/).

Here is an example config using it:

```yaml
docker:
  # The name of the image to use
  image: "ghcr.io/sethcurry/abyss-hermes:latest"
  agent_command: 
    - hermes-acp
```
