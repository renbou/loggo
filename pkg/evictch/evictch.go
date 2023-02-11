// Package evictch contains an implementation of a buffered channel with the ability to evict old
// entries once the channel is filled entirely.
//
// Note: the current implementation supports only 1 reader and 1 writer running concurrently.
package evictch

import (
	"container/list"
	"sync"
)

// Chan is an implementation of a buffered channel with automatic eviction of oldest entries.
type Chan[T any] struct {
	capacity uint

	list   *list.List
	closed bool

	mu   sync.Mutex
	cond sync.Cond
}

// NewChan constructs a new buffered channel with eviction. A maximum of capacity elements will be retained.
// Note that capacity must be greater than 0, as the whole idea of such a channel is to allow write operations
// which will not wait for a reader to arrive.
func NewChan[T any](capacity uint) *Chan[T] {
	ch := &Chan[T]{
		capacity: capacity,
		list:     list.New(),
	}
	ch.cond = sync.Cond{L: &ch.mu}

	return ch
}

// Write writes a new value to the buffered channel. If it is already full, the oldest (so, the head of the queue)
// element is removed first, to meet the expected channel capacity.
func (ch *Chan[T]) Write(v T) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.list.Len() == int(ch.capacity) {
		// Remove the head first
		_ = ch.popFront()
	}

	_ = ch.list.PushBack(v)
	ch.cond.Signal()
}

// Read reads a value from the channel. If one is available in the buffer,
// it is instantly chosen, otherwise, Read waits for a Write using sync.Cond.
func (ch *Chan[T]) Read() (T, bool) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if ch.list.Len() == 0 {
		// Wait for a write. No need to check any conditions here, since we have only 1 reader and 1 writer.
		ch.cond.Wait()
	}

	var v T
	if !ch.closed {
		v = ch.popFront()
	}

	return v, !ch.closed
}

// Close closes the channel. This should be called after the writer is finished writing.
func (ch *Chan[T]) Close() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.closed = true
	ch.cond.Signal()
}

func (ch *Chan[T]) popFront() T {
	return ch.list.Remove(ch.list.Front()).(T)
}
