# prompt_filter

This is an example plugin that uses `abyss.NewACPPluginRouter` to receive pre-unmarshalled ACP objects on object-specific methods, rather than handling raw protobuf messages.

When a user submits a prompt, the plugin checks the prompt contents against a list of banned regexes. If any match, the original prompt is dropped and a snarky `SessionNotification` is sent in its place. Otherwise, the prompt is passed through unmodified.

If you'd rather work with the raw protobuf messages directly, see the [global_logger](../global_logger/README.md) example, which implements the `protobyss.ACPPlugin` interface itself.

## Writing your own plugin

Plugins are Go modules compiled to WASI (wazero-compatible) and loaded by abyss at startup. To turn this example into your own plugin:

### 1. Copy the directory

Copy this directory and rename the `PromptFilter` struct in `main.go` to something meaningful for your plugin.

Keep the `//go:build wasip1` build tag at the top of the file; it's required for the WASM build to work.

### 2. Add the handler methods you care about

`ACPPluginRouter` dispatches each ACP message type to a corresponding `On...` method on your struct. The full list of supported methods lives in [pkg/abyss/plugin_router_builder.go](../../../pkg/abyss/plugin_router_builder.go) (e.g. `OnSessionNotification`, `OnWriteTextFileRequest`, `OnRequestPermissionResponse`, ...).

You only implement the methods you care about; `NewACPPluginRouter` uses interface assertions to detect them, and message types with no matching method simply pass through untouched. The handler signature is:

```go
func (p *MyPlugin) OnPromptRequest(req acp.PromptRequest) ([]*protobyss.ACPContainer, error)
```

The `acp` structs come from `github.com/coder/acp-go-sdk`.

### 3. Decide what to return

Your handler's return value determines what happens next in the message pipeline:

- **Pass through:** return the original message, as this example does with `abyss.ACPContainers(req)`. You can also modify the struct before returning it.
- **Replace the message:** return something else entirely, like the `SessionNotification` in this example, which swallows the original prompt.
- **Emit multiple messages:** return several messages at once; each is wrapped into its own `ACPContainer` and processed in order.

Use `abyss.ACPContainers(...)` to convert typed `acp` structs into the `[]*protobyss.ACPContainer` that handlers must return. Message IDs and `ResponseFor` fields are filled in automatically where appropriate, so you don't need to set them yourself.

### 4. Register the plugin

Wire it up in `init()` — the empty `main()` just needs to exist so the package compiles as a `main` package:

```go
func init() {
	myPlugin := &MyPlugin{}
	plugin := abyss.NewACPPluginRouter(myPlugin)
	protobyss.RegisterACPPlugin(plugin)
}
```

### 5. Build

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go
```
