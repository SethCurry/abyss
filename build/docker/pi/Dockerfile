# syntax=docker/dockerfile:1

FROM abyss-base-js:latest

RUN npm install -g @earendil-works/pi-coding-agent pi-acp

RUN ln -s /usr/local/bin/pi /usr/bin/pi
RUN ln -s /usr/local/bin/pi-acp /usr/bin/pi-acp

ENV PATH=/usr/local/bin:/usr/bin:/usr/sbin
