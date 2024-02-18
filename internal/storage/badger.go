package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

const (
	// 100 should be a good default value for our non-high-throughput usecase
	// :TODO: make this configurable
	logSequenceBandwidth = 100
	badgerGCInterval     = time.Minute * 5
	badgerGCRatio        = 0.5
)

// Badger is the realtime storage implementation using BadgerDB.
type Badger struct {
	db       *badger.DB
	sequence *badger.Sequence
	queue    *realtimeQueue
	done     chan struct{}
}

// NewBadger opens a new badger DB using the specified options.
func NewBadger(opts badger.Options) (*Badger, error) {
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening badger db: %w", err)
	}

	sequence, err := db.GetSequence(logsSequenceKey, logSequenceBandwidth)
	if err != nil {
		return nil, fmt.Errorf("getting log message sequence from badger: %w", err)
	}

	b := &Badger{
		db:       db,
		sequence: sequence,
		queue:    newRealtimeQueue(),
		done:     make(chan struct{}),
	}

	// Start the GC in the background. It will be stopped gracefully on close.
	go b.runGC()

	return b, nil
}

// Close releases the ID sequence and closes the DB.
func (b *Badger) Close() error {
	seqErr := b.sequence.Release()
	dbErr := b.db.Close()

	close(b.done)

	return errors.Join(seqErr, dbErr)
}

// AddMessage adds a new message to the storage and sends it to any fitting realtime consumers.
// :TODO: batch message additions for higher throughput
func (b *Badger) AddMessage(t time.Time, m Message) error {
	// First prepare the value
	flat := flatten(m)
	value, err := prepareMessage(m, flat)
	if err != nil {
		return err
	}

	// Then expand the sequence
	id, err := b.sequence.Next()
	if err != nil {
		return fmt.Errorf("retrieving next message id from sequence: %w", err)
	}
	key := messageKey(t, id)

	// Finally, start the transaction
	err = b.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, value); err != nil {
			return fmt.Errorf("adding message with key %s: %w", key, err)
		}
		return nil
	})
	if err != nil {
		// ErrConflict and other errors don't need to be handled here because all we do is a single Set
		return fmt.Errorf("updating db: %w", err)
	}

	// Finally, send this message through the realtime requests
	b.queue.iterate(m, flat)
	return nil
}

// ListMessages returns a batched, reversed (latest to earliest) list of all of the log messages in the interval [from, to].
// Nil may be passed here for a filter, in which case no filtering occurs.
// If after is set to a non-empty sequence, the iteration is started from "after" instead of "from".
func (b *Badger) ListMessages(from, to time.Time, filter Filter, after []byte, limit uint) (Batch, error) {
	l := messagePrefix(from)
	r := after

	if len(r) == 0 {
		// Add 1 nanosecond so that the interval is ], not )
		r = messagePrefix(to.Add(time.Nanosecond))
	}

	// Run iterator in a read-only transaction
	var batch Batch
	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Reverse = true

		it := txn.NewIterator(opts)
		defer it.Close()

		// Start from the latest message and iterate while the interval left edge >= current key (since keys are timestamps)
		for it.Seek(r); it.Valid(); it.Next() {
			item := it.Item()
			if bytes.Compare(l, item.Key()) > 0 {
				break
			}

			if limit != 0 && uint(len(batch.Messages)) == limit {
				batch.Next = item.KeyCopy(nil)
				break
			}

			// Add the stored message to the batch if needed
			err := item.Value(func(val []byte) error {
				// After unprepare the message can safely be used because it is already a copy.
				message, flat, err := unprepareMessage(val)
				if err != nil {
					return fmt.Errorf("unmarshaling db item: %w", err)
				}

				// Only append the message if it actually passes any filtering.
				if filter == nil || filter(message, flatMessageToMapping(flat)) {
					batch.Messages = append(batch.Messages, message)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("getting message from db item: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return Batch{}, fmt.Errorf("iterating over items: %w", err)
	}

	return batch, nil
}

// StreamMessages should be used when messages need to be streamed in real-time starting from some
// timestamp. Like ListMessages, a Batch is returned consisting of the latest messages after "from",
// meaning that, if there are lots of such messages, only a small part is returned in the Batch
// and its Next field is set, which can then be used with ListMessages like usual.
// New messages that arrive after the batch are passed through the returned channel until the passed context is done.
func (b *Badger) StreamMessages(ctx context.Context, from time.Time, filter Filter, limit uint,
) (Batch, chan Message, error) {
	batch, err := b.ListMessages(from, time.Now(), filter, nil, limit)
	if err != nil {
		return Batch{}, nil, fmt.Errorf("retrieving first batch: %w", err)
	}

	// After getting the first batch we can safely create a channel and register it with the realtime queue.
	// Deletion is then implemented using a simple goroutine watching for ctx.Done.
	ch := make(chan Message)
	el := b.queue.add(&realtimeRequest{filter, ch})
	go func() {
		<-ctx.Done()
		b.queue.delete(el)

		// close specifically *after* the element is gone from the queue
		close(ch)
	}()

	return batch, ch, nil
}

// runGC runs the badgerDB log GC in the background, which should be good enough for our simple use case
func (b *Badger) runGC() {
	ticker := time.NewTicker(badgerGCInterval)
	defer ticker.Stop()

	for {
		select {
		case <-b.done:
			return
		case <-ticker.C:
		}

		_ = b.db.RunValueLogGC(badgerGCRatio)
	}
}
