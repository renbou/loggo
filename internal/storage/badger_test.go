package storage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func exampleTime(s int) time.Time {
	return time.Date(2000, time.January, 1, 0, 0, s, 0, time.UTC)
}

func openTestBadger(t *testing.T) *Badger {
	t.Helper()

	testOpts := badger.DefaultOptions("").WithInMemory(true)
	b, err := NewBadger(testOpts)

	assert.NoError(t, err)
	return b
}

func Test_Badger_OpenClose(t *testing.T) {
	t.Parallel()

	b := openTestBadger(t)

	// Sequence must be initialized and return a proper ID
	expectSequenceID := 0

	gotSequenceID, err := b.sequence.Next()
	require.NoError(t, err)

	assert.EqualValues(t, expectSequenceID, gotSequenceID)

	// DB must be functioning
	err = b.db.View(func(txn *badger.Txn) error {
		return nil
	})
	assert.NoError(t, err)

	err = b.db.Update(func(txn *badger.Txn) error {
		return nil
	})
	assert.NoError(t, err)

	// DB shouldn't error on expected Close
	assert.NoError(t, b.Close())
}

func Test_Badger_AddMessage(t *testing.T) {
	t.Parallel()

	b := openTestBadger(t)
	defer b.Close()

	message := Message(`{"test": "message"}`)
	messageT := exampleTime(0)
	expectedSequenceID := uint64(0)

	// Addition of a message should never fail
	err := b.AddMessage(messageT, message)
	require.NoError(t, err)

	// After AddMessage returns, the message should be added with the expected (logs:...:0) key and be accessible from it
	err = b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(messageKey(messageT, expectedSequenceID))
		assert.NoError(t, err)

		err = item.Value(func(val []byte) error {
			gotMessage, _, err := unprepareMessage(val)
			assert.NoError(t, err)
			assert.EqualValues(t, message, gotMessage)
			return nil
		})
		assert.NoError(t, err)

		return nil
	})
	assert.NoError(t, err)
}

func Test_Badger_ListMessages(t *testing.T) {
	t.Parallel()

	b := openTestBadger(t)

	// Add a couple of messages to the DB
	add := []struct {
		m Message
		t time.Time
	}{
		{Message(`{"a": "b"}`), exampleTime(1)},
		{Message(`{"c": "d"}`), exampleTime(3)},
		{Message(`{"e": "f"}`), exampleTime(5)},
	}

	for _, data := range add {
		err := b.AddMessage(data.t, data.m)
		require.NoError(t, err)
	}

	// Run different batch listing tests
	tests := []struct {
		name        string
		from        time.Time
		to          time.Time
		filter      Filter
		after       []byte
		limit       uint
		expectBatch Batch
	}{
		{
			name: "query all",
			from: exampleTime(1),
			to:   exampleTime(5),
			expectBatch: Batch{
				Messages: []Message{add[2].m, add[1].m, add[0].m},
			},
		},
		{
			name: "query single item",
			from: exampleTime(2),
			to:   exampleTime(4),
			expectBatch: Batch{
				Messages: []Message{add[1].m},
			},
		},
		{
			name: "query exact time",
			from: exampleTime(3),
			to:   exampleTime(3),
			expectBatch: Batch{
				Messages: []Message{add[1].m},
			},
		},
		{
			name:  "query all with limit",
			from:  exampleTime(1),
			to:    exampleTime(5),
			limit: 1,
			expectBatch: Batch{
				Messages: []Message{add[2].m},
				Next:     []byte("logs:00946684803000000000:1"),
			},
		},
		{
			name:  "query all with limit and after",
			from:  exampleTime(1),
			to:    exampleTime(5),
			after: []byte("logs:00946684803000000000:1"),
			limit: 2,
			expectBatch: Batch{
				Messages: []Message{add[1].m, add[0].m},
			},
		},
		{
			name: "query with filter",
			from: exampleTime(2),
			to:   exampleTime(5),
			filter: func(_ Message, fm FlatMapping) bool {
				v, ok := fm("c")
				return ok && v == "d"
			},
			expectBatch: Batch{
				Messages: []Message{add[1].m},
			},
		},
	}

	// Need to wait for all tests to complete before closing the DB
	var wg sync.WaitGroup
	wg.Add(len(tests))
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			defer wg.Done()

			gotBatch, err := b.ListMessages(tt.from, tt.to, tt.filter, tt.after, tt.limit)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectBatch, gotBatch)
		})
	}

	go func() {
		wg.Wait()
		b.Close()
	}()
}

func Test_Badger_StreamMessages(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	add := []struct {
		m Message
		t time.Time
	}{
		{Message(`{"a": "b", "field": 1337}`), exampleTime(1)},
		{Message(`{"c": "d"}`), exampleTime(3)},
		{Message(`{"e": "f", "field": 1337}`), exampleTime(5)},
	}
	expectBatch := Batch{Messages: []Message{add[0].m}}
	expectChMessages := []Message{add[2].m}

	b := openTestBadger(t)
	defer b.Close()

	// Add a single message
	err := b.AddMessage(add[0].t, add[0].m)
	require.NoError(t, err)

	// Batch should contain a single message right now
	gotBatch, ch, err := b.StreamMessages(ctx, exampleTime(0), func(_ Message, fm FlatMapping) bool {
		v, ok := fm("field")
		return ok && v == "1337"
	}, 0)

	require.NoError(t, err)
	require.Equal(t, expectBatch, gotBatch)

	// A separate goroutine should receive only the valid filtered messages
	var wg sync.WaitGroup
	var gotChMessages []Message
	wg.Add(1)
	go func() {
		defer wg.Done()

		for message := range ch {
			gotChMessages = append(gotChMessages, message)
		}
	}()

	// Add the other two messages
	for i := 1; i < len(add); i++ {
		err := b.AddMessage(add[i].t, add[i].m)
		require.NoError(t, err)
	}
	cancel()

	wg.Wait()
	assert.Equal(t, expectChMessages, gotChMessages)
	assert.Equal(t, 0, b.queue.list.Len())
}
