package wsrouter

import (
	"strconv"
	"testing"
	"time"

	"github.com/SethCurry/abyss/internal/websockets/protobyss"
)

// TestNewResponseWatcher initializes an empty watcher with no registered handlers.
func TestNewResponseWatcher(t *testing.T) {
	rw := NewResponseWatcher()
	if rw == nil {
		t.Fatal("expected non-nil ResponseWatcher")
	}
	if len(rw.handlers) != 0 {
		t.Fatalf("expected no handlers, got %d", len(rw.handlers))
	}
}

// TestRegisterReturnsWaitingPromise verifies Register records a promise that
// blocks until a matching response is handled.
func TestRegisterReturnsWaitingPromise(t *testing.T) {
	rw := NewResponseWatcher()
	prom := rw.Register("req-1")

	if prom == nil {
		t.Fatal("expected non-nil promise")
	}
	if _, ok := rw.handlers["req-1"]; !ok {
		t.Fatal("expected handler registered for req-1")
	}

	select {
	case <-prom.resolveChan:
		t.Fatal("promise should not be resolved before Handle")
	default:
	}
}

// TestHandleResolvesPromise verifies Handle resolves the registered promise with
// the matching response and removes the handler.
func TestHandleResolvesPromise(t *testing.T) {
	rw := NewResponseWatcher()
	prom := rw.Register("req-1")

	msg := &protobyss.ACPContainer{ResponseFor: "req-1"}
	go func() {
		time.Sleep(10 * time.Millisecond)
		rw.Handle(nil, msg)
	}()

	got := prom.Wait()
	if got != msg {
		t.Fatalf("expected resolved message %v, got %v", msg, got)
	}

	rw.mut.Lock()
	defer rw.mut.Unlock()
	if _, ok := rw.handlers["req-1"]; ok {
		t.Fatal("expected handler to be removed after Handle")
	}
}

// TestHandleUnknownResponseFor verifies Handle is a no-op when no handler matches.
func TestHandleUnknownResponseFor(t *testing.T) {
	rw := NewResponseWatcher()

	// Should not panic or block.
	rw.Handle(nil, &protobyss.ACPContainer{ResponseFor: "missing"})

	if len(rw.handlers) != 0 {
		t.Fatalf("expected no handlers registered, got %d", len(rw.handlers))
	}
}

// TestRegisterOverwritesExisting verifies a second Register for the same ID replaces
// the prior promise, leaving the old one unresolved.
func TestRegisterOverwritesExisting(t *testing.T) {
	rw := NewResponseWatcher()
	first := rw.Register("req-1")
	second := rw.Register("req-1")

	rw.mut.Lock()
	defer rw.mut.Unlock()
	if rw.handlers["req-1"] != second {
		t.Fatal("expected second promise to replace first")
	}

	select {
	case <-first.resolveChan:
		t.Fatal("first promise should remain unresolved after overwrite")
	default:
	}
}

// TestResolveDeliversValue verifies Promise.Resolve sends the value on the channel.
func TestResolveDeliversValue(t *testing.T) {
	prom := &Promise[*protobyss.ACPContainer]{
		resolveChan: make(chan *protobyss.ACPContainer, 1),
	}
	msg := &protobyss.ACPContainer{}
	prom.Resolve(msg)

	got := prom.Wait()
	if got != msg {
		t.Fatalf("expected %v, got %v", msg, got)
	}
}

// TestWaitBlocksUntilResolved verifies Wait blocks until Resolve is called.
func TestWaitBlocksUntilResolved(t *testing.T) {
	prom := &Promise[*protobyss.ACPContainer]{
		resolveChan: make(chan *protobyss.ACPContainer),
	}

	select {
	case <-time.After(20 * time.Millisecond):
	case got := <-prom.resolveChan:
		t.Fatalf("expected Wait to block, got %v", got)
	}

	msg := &protobyss.ACPContainer{}
	go func() {
		time.Sleep(10 * time.Millisecond)
		prom.Resolve(msg)
	}()

	got := prom.Wait()
	if got != msg {
		t.Fatalf("expected %v, got %v", msg, got)
	}
}

// TestHandleConcurrent verifies concurrent Register/Handle operations do not race.
func TestHandleConcurrent(t *testing.T) {
	rw := NewResponseWatcher()
	done := make(chan struct{})

	for i := 0; i < 50; i++ {
		go func(i int) {
			id := "req-" + strconv.Itoa(i)
			prom := rw.Register(id)
			go func() {
				rw.Handle(nil, &protobyss.ACPContainer{ResponseFor: id})
			}()
			_ = prom.Wait()
		}(i)
	}

	go func() {
		<-time.After(100 * time.Millisecond)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for concurrent operations")
	}
}
