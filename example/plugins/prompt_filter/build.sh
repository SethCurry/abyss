#!/bin/bash

GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go

