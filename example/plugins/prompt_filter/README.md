This is an example plugin that uses abyss.ACPPluginRouter to get pre-unmarshalled ACP objects on object-specific methods.

## Build

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go
```
