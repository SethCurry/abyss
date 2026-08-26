# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM ubuntu:latest AS build

WORKDIR /src

RUN apt-get update && apt-get install -y golang ca-certificates

# Cache dependencies first for faster incremental builds.
COPY go.mod go.sum ./

RUN go mod download

# Copy the source and build the binary.
COPY . .

RUN CGO_ENABLED=0 go build -o /out/abyss ./cmd/abyss

# ---- Runtime stage ----
FROM ubuntu:latest

COPY --from=build /out/abyss /usr/local/bin/abyss
