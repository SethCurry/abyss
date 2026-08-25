#!/bin/bash

go test ./... -coverprofile=/tmp/abyss.cover.out

go tool cover -html=/tmp/abyss.cover.out
