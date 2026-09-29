#!/bin/bash

GOOS=wasip1 GOARCH=wasm go build -o secrets_filter.wasm -buildmode=c-shared secrets_filter.go
