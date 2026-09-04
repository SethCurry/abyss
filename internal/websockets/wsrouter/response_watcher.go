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

type ResponseWatcher struct {
	handlers map[string]*Promise[*protobyss.Container]
	mut      sync.Mutex
}

func (r *ResponseWatcher) Handle(router *ACPRouter, msg *protobyss.Container) {
	r.mut.Lock()
	defer r.mut.Unlock()
	if handler, ok := r.handlers[msg.ResponseFor]; ok {
		handler.Resolve(msg)
		delete(r.handlers, msg.ResponseFor)
	}
}

func (r *ResponseWatcher) Register(requestID string) *Promise[*protobyss.Container] {
	r.mut.Lock()
	defer r.mut.Unlock()
	prom := &Promise[*protobyss.Container]{
		resolveChan: make(chan *protobyss.Container),
	}
	r.handlers[requestID] = prom

	return prom
}

type Promise[T any] struct {
	resolveChan chan T
}

func (p *Promise[T]) Resolve(value T) {
	p.resolveChan <- value
	close(p.resolveChan)
}

func (p *Promise[T]) Wait() T {
	return <-p.resolveChan
}
