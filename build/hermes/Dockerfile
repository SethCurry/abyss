# syntax=docker/dockerfile:1

FROM ghcr.io/sethcurry/abyss-base:dev

RUN apt-get update && apt-get install -y ca-certificates python3 curl

RUN curl -fsSL https://hermes-agent.nousresearch.com/install.sh | bash
RUN cd ~/.hermes/hermes-agent && uv pip install -e '.[acp]'

ENV PATH=/usr/local/bin:/usr/bin:/usr/sbin
