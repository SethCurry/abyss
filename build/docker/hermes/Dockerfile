# syntax=docker/dockerfile:1

FROM ghcr.io/sethcurry/abyss-base:latest

ENV PATH=/usr/local/bin:/usr/bin:/usr/sbin:/root/.local/bin:/bin

RUN apt-get update && apt-get install -y ca-certificates python3 curl git libatomic1
RUN curl -LsSf https://astral.sh/uv/install.sh | sh
RUN curl -fsSL https://hermes-agent.nousresearch.com/install.sh | bash
RUN cd ~/.hermes/hermes-agent && uv venv && uv pip install -e '.[acp]'
