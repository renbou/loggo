package storage

import (
	"fmt"
	"time"
)

const logsPrefix = "logs"

var logsSequenceKey = []byte("sequence:logs")

// messagePrefix returns a common for all log messages in the same nanosecond
// of course this will always be unique in a normal use case, but just in case,
// all keys will be prefixed using integers from the badger db sequence
func messagePrefix(t time.Time) []byte {
	return []byte(fmt.Sprintf("%s:%020d", logsPrefix, t.UnixNano()))
}

// messageKey is like messagePrefix but with an additional sequence number to guarantee uniqueness
func messageKey(t time.Time, seq uint64) []byte {
	return []byte(fmt.Sprintf("%s:%020d:%d", logsPrefix, t.UnixNano(), seq))
}
