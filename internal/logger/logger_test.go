package logger

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Logger(t *testing.T) {
	// Not parallel because it uses the unsynchronized global logger

	// Undo everything before running other tests which might use the global logger
	baseLogger := globalLogger
	defer func() {
		globalLogger = baseLogger
		globalSink.consumer = nil
	}()

	var gotMessages [][]byte
	consumer := func(t time.Time, message []byte) error {
		gotMessages = append(gotMessages, message)
		return nil
	}

	// Initialize the global logger once
	assert.NoError(t, InitGlobal(consumer))

	// Log some messages to it
	Infow("test informational message", "field", "value", "number", 1337)
	Errorw("test error message", "err", errors.New("fake error"))

	// After Sync all messages should have arrived at the consumer
	assert.NoError(t, Sync())

	expectMessages := []string{
		`^{"level":"info","ts":".+","msg":"test informational message","service":"loggo","field":"value","number":1337}$`,
		`^{"level":"error","ts":".+","msg":"test error message","service":"loggo","err":"fake error"}$`,
	}

	for i := range expectMessages {
		assert.Regexp(t, expectMessages[i], string(gotMessages[i]))
	}
}
