# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies first for faster incremental builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source and build the binary.
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/abyss ./cmd/abyss

# ---- Runtime stage ----
FROM alpine:3.21

# CA certificates are needed for TLS websocket connections.
# NodeJS is required to run the Pi LLM coding agent (https://pi.dev).
RUN apk add --no-cache ca-certificates nodejs npm bash

COPY --from=build /out/abyss /usr/local/bin/abyss

# Install the Pi LLM coding agent globally.
RUN npm install -g --ignore-scripts @earendil-works/pi-coding-agent
RUN npm install -g --ignore-scripts pi-acp
RUN ln -s /usr/local/bin/pi /usr/bin/pi
RUN ln -s /usr/local/bin/pi-acp /usr/bin/pi-acp

ENTRYPOINT ["/usr/local/bin/abyss"]
