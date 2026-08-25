# syntax=docker/dockerfile:1

FROM ghcr.io/sethcurry/abyss-base:latest

RUN apt-get update && apt-get install -y ca-certificates nodejs npm
RUN npm install -g @earendil-works/pi-coding-agent pi-acp
RUN apt-get clean autoclean && apt-get autoremove --yes && rm -rf /var/lib/{apt,dpkg,cache,log}/

ENTRYPOINT ["/usr/local/bin/abyss"]
