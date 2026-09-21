This is an example plugin that uses the AcpPlugin interface directly (i.e. it receives protobuf messages and has to do its own unmarshalling).

If you want to receive pre-unmarshalled and typed messages, take a look at [prompt_filter](../prompt_filter/README.md).

## Build

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go
```
