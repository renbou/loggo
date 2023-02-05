package storage

import (
	"bytes"
	"container/list"
	"strings"
	"sync"
	"testing"

	"github.com/renbou/obzerva/internal/logs/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRealtimeConsumer is run from its own goroutine and tries to satisfy
// all of the wanted messages by receiving messages from its assigned realtime channel
type mockRealtimeConsumer struct {
	queueElement *list.Element
	ch           <-chan Message
	want         map[string]int
}

func (mcr *mockRealtimeConsumer) run(t *testing.T) {
	t.Helper()

	for message := range mcr.ch {
		s := string(message)

		if !assert.Contains(t, mcr.want, s) {
			continue
		}

		// decrease want counter and delete if it's now at 0
		next := mcr.want[s] - 1
		if next == 0 {
			delete(mcr.want, s)
		} else {
			mcr.want[s] = next
		}
	}

	assert.Empty(t, mcr.want)
}

// mockMessageProducer simply spams messages to the realtime queue.
// In reality, messages would come from realtime writes to the store.
type mockMessageProducer struct {
	queue   *realtimeQueue
	message Message
	flat    *models.FlatMessage
	n       int
}

func (mmp *mockMessageProducer) run(t *testing.T) {
	t.Helper()

	for i := 0; i < mmp.n; i++ {
		mmp.queue.iterate(mmp.message, mmp.flat)
	}
}

func addMockRealtimeConsumer(q *realtimeQueue, realtimeCh chan Message,
	filter Filter, want map[string]int,
) *mockRealtimeConsumer {
	el := q.add(&realtimeRequest{filter: filter, ch: realtimeCh})
	return &mockRealtimeConsumer{
		queueElement: el,
		ch:           realtimeCh,
		want:         want,
	}
}

func Test_RealtimeQueue_Parallel(t *testing.T) {
	t.Parallel()

	// Each producer will produce the assigned message 50 times
	const producerN = 50

	queue := newRealtimeQueue()
	require.NotNil(t, queue)

	// Prepare 4 different messages which will be sent from multiple concurrent goroutines
	messages := []Message{
		Message(`{"key": "some pointless message", "a": 13337}`),
		Message(`{"nested": {"value": true}}`),
		Message(`non-JSON messages are also valid`),
		Message(`[1, 2, 3, 4]`),
	}

	var flat []models.FlatMessage
	var channels []chan Message
	for _, message := range messages {
		flat = append(flat, flatten(message))
		channels = append(channels, make(chan Message))
	}

	// Create 4 consumers each with a unique filter. They will run and wait for updates from the realtime queue.
	consumers := []*mockRealtimeConsumer{
		// Simple field filter
		addMockRealtimeConsumer(queue, channels[0],
			func(m Message, fm FlatMapping) bool {
				v, ok := fm("key")
				return ok && strings.Contains(v, "pointless")
			},
			map[string]int{
				string(messages[0]): producerN,
			}),

		// Nested filter
		addMockRealtimeConsumer(queue, channels[1],
			func(m Message, fm FlatMapping) bool {
				v, ok := fm("nested.value")
				return ok && v == "true"
			},
			map[string]int{
				string(messages[1]): producerN,
			}),

		// Non-scoped filter and field filter
		addMockRealtimeConsumer(queue, channels[2],
			func(m Message, fm FlatMapping) bool {
				if bytes.Contains(m, []byte("non-JSON")) {
					return true
				}

				v, ok := fm("a")
				return ok && v == "13337"
			},
			map[string]int{
				string(messages[0]): producerN,
				string(messages[2]): producerN,
			}),

		// No filter, should receive everything
		addMockRealtimeConsumer(queue, channels[3], nil, map[string]int{
			string(messages[0]): producerN,
			string(messages[1]): producerN,
			string(messages[2]): producerN,
			string(messages[3]): producerN,
		}),
	}

	// And start each of the consumers
	var consumerWg sync.WaitGroup
	consumerWg.Add(len(consumers))
	for _, consumer := range consumers {
		consumer := consumer
		go func() {
			defer consumerWg.Done()
			consumer.run(t)
		}()
	}

	// Create and run all of the producers
	var producerWg sync.WaitGroup
	producerWg.Add(len(messages))
	for i := range messages {
		producer := mockMessageProducer{
			queue:   queue,
			message: messages[i],
			flat:    &flat[i],
			n:       producerN,
		}

		go func() {
			defer producerWg.Done()
			producer.run(t)
		}()
	}

	// Wait for all of the producers to stop, then wait for the consumers to stop
	producerWg.Wait()
	for _, channel := range channels {
		close(channel)
	}
	consumerWg.Wait()

	// Finally, clear the queue
	for _, consumer := range consumers {
		queue.delete(consumer.queueElement)
	}

	assert.Equal(t, 0, queue.list.Len())
}
