package logger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func assertingConsumer(t *testing.T, expect string, err error) LogConsumer {
	var called bool
	return func(tm time.Time, message []byte) error {
		assert.False(t, called)
		assert.Equal(t, expect, string(message))
		called = true
		return err
	}
}

func Test_ConsumerSink_Write(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		consumer LogConsumer
		chunks   []string
		wantErr  bool
	}{
		{
			name:     "single chunk",
			consumer: assertingConsumer(t, "message", nil),
			chunks:   []string{"message\n"},
		},
		{
			name:     "multiple chunks",
			consumer: assertingConsumer(t, "split message", nil),
			chunks:   []string{"split ", "message\n"},
		},
		{
			name:     "trailing data",
			consumer: assertingConsumer(t, "no trailing", nil),
			chunks:   []string{"no trailing", "\ndata"},
		},
		{
			name:     "error during send",
			consumer: assertingConsumer(t, "failed message", assert.AnError),
			chunks:   []string{"failed message\n"},
			wantErr:  true,
		},
		{
			name:     "nil consumer",
			consumer: nil,
			chunks:   []string{"datadata", "more data\n"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sink := &consumerSink{consumer: tt.consumer, lineEnding: []byte(newline)}

			for _, chunk := range tt.chunks {
				n, err := sink.Write([]byte(chunk))

				assert.Equal(t, len(chunk), n)

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func Test_ConsumerSink_Sync(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		consumer LogConsumer
		chunks   []string
		wantErr  bool
	}{
		{
			name:     "leftover data",
			consumer: assertingConsumer(t, "during sync", nil),
			chunks:   []string{"during sync"},
		},
		{
			name:     "no leftover data",
			consumer: assertingConsumer(t, "all data already sent", nil),
			chunks:   []string{"all", " data", " already sent\n"},
		},
		{
			name:     "error during send",
			consumer: assertingConsumer(t, "unsynced", assert.AnError),
			chunks:   []string{"un", "synced"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sink := &consumerSink{consumer: tt.consumer, lineEnding: []byte(newline)}

			for _, chunk := range tt.chunks {
				n, err := sink.Write([]byte(chunk))

				assert.Equal(t, len(chunk), n)
				assert.NoError(t, err)
			}

			err := sink.Sync()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
