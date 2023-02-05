package storage

import (
	"container/list"
	"sync"

	"github.com/renbou/obzerva/internal/logs/storage/models"
)

// realtimeRequest specifies a request which wants to receive logs in real-time
// as they are being written to the store. Only the List* method in which the request
// was created should close its channel, after the request was canceled.
type realtimeRequest struct {
	filter Filter
	ch     chan<- Message
}

// realtimeQueue stores all of the current realtime requests and operates on them using a single lock.
// Technically, we could construct a lock-free queue here, but the lock allows us to easily avoid
// closing a realtimeRequest channel before it's actually safe to do so.
type realtimeQueue struct {
	mu   sync.RWMutex
	list *list.List
}

func newRealtimeQueue() *realtimeQueue {
	return &realtimeQueue{list: list.New()}
}

// add adds a new request to the end of the queue and returns a reference to it so that it can later be deleted
func (q *realtimeQueue) add(r *realtimeRequest) *list.Element {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.list.PushBack(r)
}

// delete deletes a previouosly added element. Once an element is deleted from this queue, thanks to the mutex,
// no iterations will read this request, allowing its channel to be closed safely.
func (q *realtimeQueue) delete(e *list.Element) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.list.Remove(e)
}

// iterate iterates over the existing realtime requests, performing filtering on the message, if needed,
// and sending it out on the request's channel
func (q *realtimeQueue) iterate(message Message, flat *models.FlatMessage) {
	mapping := flatMessageToMapping(flat)

	q.mu.RLock()
	defer q.mu.RUnlock()

	for el := q.list.Front(); el != nil; el = el.Next() {
		r := el.Value.(*realtimeRequest)

		if r.filter == nil || r.filter(message, mapping) {
			r.ch <- message
		}
	}
}
