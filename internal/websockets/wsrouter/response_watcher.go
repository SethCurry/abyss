package wsrouter

import (
	"sync"

	"github.com/SethCurry/abyss/internal/timber"
	"github.com/SethCurry/abyss/pkg/protobyss"
	"github.com/rs/zerolog"
)

// NewResponseWatcher creates a new ResponseWatcher with an empty handler map.
func NewResponseWatcher() *ResponseWatcher {
	return &ResponseWatcher{
		handlers: make(map[string]*Promise[*protobyss.ACPContainer]),
		logger:   timber.ComponentLogger("ResponseWatcher"),
	}
}

// ResponseWatcher manages the lifetime of RPC-via-websocket requests.
// protobyss.Container messages with a non-empty ResponseFor field are routed
// here so the response can be fed to the Promise.
type ResponseWatcher struct {
	handlers map[string]*Promise[*protobyss.ACPContainer]
	mut      sync.Mutex
	logger   zerolog.Logger
}

// Handle dispatches the message to the Promise that is waiting for it.
// This is a no-op if the ResponseFor field doesn't match a waiting Promise.
func (r *ResponseWatcher) Handle(router *ACPRouter, msg *protobyss.ACPContainer) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.logger.Debug().
		Str("message_id", msg.GetMessageId()).
		Msg("handling message")

	// This will just drop messages that have no listening responder
	if handler, ok := r.handlers[msg.GetResponseFor()]; ok {
		handler.Resolve(msg)
		delete(r.handlers, msg.GetResponseFor())
	} else {
		r.logger.Error().
			Str("message_id", msg.GetMessageId()).
			Msg("no handler registered for message")
	}
}

// Register creates a new *Promise and registers it as waiting for a response.
func (r *ResponseWatcher) Register(requestID string) *Promise[*protobyss.ACPContainer] {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.logger.Debug().
		Str("request_id", requestID).
		Msg("registering handler for response")
	prom := &Promise[*protobyss.ACPContainer]{
		resolveChan: make(chan *protobyss.ACPContainer),
	}

	// TODO this will leak for requests that never get answered
	r.handlers[requestID] = prom

	return prom
}

// Promise is a basic re-implementation of NodeJS promises backed by a
// single-use channel.  It avoids repeating the rigamarole of "wait for value,
// close channel, move on".
type Promise[T any] struct {
	resolveChan chan T
}

// Resolve resolves the promise with the provided value.
func (p *Promise[T]) Resolve(value T) {
	p.resolveChan <- value
	close(p.resolveChan)
}

// Wait blocks until the promise is resolved with a value.
func (p *Promise[T]) Wait() T {
	return <-p.resolveChan
}
