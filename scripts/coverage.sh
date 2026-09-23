#!/bin/bash

go test ./... -coverprofile=/tmp/abyss.cover.out

#go tool cover -html=/tmp/abyss.cover.out
go tool cover -func=/tmp/abyss.cover.out | grep -v 'abyss/pkg/protobyss' | grep -v '100.0%'
