# global_logger

This is an example plugin that implements the `protobyss.ACPPlugin` interface directly: it receives raw, serialized messages and must do its own unmarshalling. In exchange, it sees *every* message that flows through abyss, in both directions, regardless of type.

The plugin logs each message's raw content to stdout and passes everything through untouched.

If you'd rather receive pre-unmarshalled, typed messages on per-type handler methods, take a look at [prompt_filter](../prompt_filter/README.md) instead. Choose this raw pattern when you want to observe or handle all traffic (logging, tracing, generic middleware); choose the router when you only care about a few specific message types.

## Writing your own plugin

Plugins are Go modules compiled to WASI (wazero-compatible) and loaded by abyss at startup. To turn this example into your own plugin:

### 1. Copy the directory

Copy this directory and rename the `ACPLoggerPlugin` struct in `main.go` to something meaningful for your plugin.

Keep the `//go:build wasip1` build tag at the top of the file; it's required for the WASM build to work.

### 2. Understand `ACPContainer`

The entire `protobyss.ACPPlugin` interface is a single method:

```go
HandleMessage(ctx context.Context, msg *protobyss.ACPContainer) (*protobyss.ACPContainerList, error)
```

Every message that passes through abyss is handed to you as an `ACPContainer`, which has four fields:

- `TypeId` — an identifier for which ACP message this is. It corresponds to the `abyss.MessageTypeID` constants in [pkg/abyss/message_type.go](../../../pkg/abyss/message_type.go) (e.g. `PromptRequestType`, `SessionNotificationType`).
- `MessageId` — a unique ID for this message.
- `ResponseFor` — on responses, the `MessageId` of the request it answers.
- `Content` — the message body, JSON-encoded. It can be unmarshalled into the `acp` structs from `github.com/coder/acp-go-sdk`.

If you want to inspect a message's type, you can use `abyss.GetMessageTypeByID(msg.TypeId)` to get a `TypedMessage`, whose `UnmarshalAny(msg.Content)` will decode the body into the concrete struct. In this example we don't bother, since we just log the raw bytes either way.

### 3. Decide what to return

Your return value determines what happens next in the message pipeline:

- **Pass through:** return the original message, as this example does by wrapping `msg` in the `ACPContainerList`. Messages you return untouched continue on as if you weren't there.
- **Modify the message:** unmarshal `Content`, change the struct, re-marshal it back into `Content`, and return the same container.
- **Drop the message:** return an empty `ACPContainerList`.
- **Emit multiple messages:** return several containers; each is processed in order.

Unlike the `ACPPluginRouter` path, nothing is filled in for you — if you synthesize new messages yourself, you are responsible for setting `MessageId` and `ResponseFor` appropriately.

### 4. Register the plugin

Wire it up in `init()` — the empty `main()` just needs to exist so the package compiles under go-plugin:

```go
func init() {
	protobyss.RegisterACPPlugin(NewMyPlugin())
}
```

A compile-time assertion like the one in `main.go` catches interface mismatches at build time instead of at plugin load:

```go
var _ protobyss.ACPPlugin = (*MyPlugin)(nil)
```

### 5. Build

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm -buildmode=c-shared main.go
```