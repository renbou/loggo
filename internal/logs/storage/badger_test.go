package storage

import (
	"testing"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	messageT := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
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
