package wsrouter

import (
	"sync"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
)

func NewResponseWatcher() *ResponseWatcher {
	return &ResponseWatcher{
		handlers: make(map[string]*Promise[*protobyss.Container]),
	}
}

// ResponseWatcher is responsible for managing the lifetime of RPC-via-websocket requests.
// protobyss.Container messages that have a non-empty ResponseFor field will get routed
// here so that the response can be fed to the Promise.
type ResponseWatcher struct {
	handlers map[string]*Promise[*protobyss.Container]
	mut      sync.Mutex
}

// Handle dispatches the message to the Promise that is waiting for it.  This is a no-op
// if the ResponseFor field doesn't match with a waiting Promise.
func (r *ResponseWatcher) Handle(router *ACPRouter, msg *protobyss.Container) {
	r.mut.Lock()
	defer r.mut.Unlock()
	if handler, ok := r.handlers[msg.ResponseFor]; ok {
		handler.Resolve(msg)
		delete(r.handlers, msg.ResponseFor)
	}
}

// Register creates a new *Promise and registers it as waiting for a response.
func (r *ResponseWatcher) Register(requestID string) *Promise[*protobyss.Container] {
	r.mut.Lock()
	defer r.mut.Unlock()
	prom := &Promise[*protobyss.Container]{
		resolveChan: make(chan *protobyss.Container),
	}
	r.handlers[requestID] = prom

	return prom
}

// Promise is a very basic re-implementation of NodeJS promises powered by a single-use channel.
// It just avoids repeating the rigamarole of "wait for value, close channel, move on".
type Promise[T any] struct {
	resolveChan chan T
}

// Resolves the promise with the provided value.
func (p *Promise[T]) Resolve(value T) {
	p.resolveChan <- value
	close(p.resolveChan)
}

// Blocks until the promise is resolved with a value.
func (p *Promise[T]) Wait() T {
	return <-p.resolveChan
}
