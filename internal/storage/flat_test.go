package storage

import (
	"testing"

	"github.com/renbou/loggo/internal/storage/models"
	"github.com/stretchr/testify/assert"
)

func Test_Flatten_Positive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		fields  []*models.FlatMessage_KV
	}{
		{
			name:    "object",
			message: `{"key1": "string", "key2": 1337, "key3": 1.333344447777, "key4": false, "key5": true, "key6": null}`,
			fields: []*models.FlatMessage_KV{
				{Key: "key1", Value: "string"},
				{Key: "key2", Value: "1337"},
				{Key: "key3", Value: "1.333344447777"},
				{Key: "key4", Value: "false"},
				{Key: "key5", Value: "true"},
				{Key: "key6", Value: "null"},
			},
		},
		{
			name:    "array",
			message: `["string", 1337, 1.333344447777, false, true, null]`,
			fields: []*models.FlatMessage_KV{
				{Key: "0", Value: "string"},
				{Key: "1", Value: "1337"},
				{Key: "2", Value: "1.333344447777"},
				{Key: "3", Value: "false"},
				{Key: "4", Value: "true"},
				{Key: "5", Value: "null"},
			},
		},
		{
			name:    "top-level string",
			message: `"string"`,
			fields:  nil,
		},
		{
			name:    "top-level number",
			message: `1337`,
			fields:  nil,
		},
		{
			name:    "top-level boolean",
			message: `false`,
			fields:  nil,
		},
		{
			name:    "top-level null",
			message: `null`,
			fields:  nil,
		},
		{
			name:    "nested object",
			message: `{"empty object": {}, "non-empty object'": {"key1": "value", "key2": [0, null]}}`,
			fields: []*models.FlatMessage_KV{
				{Key: "non_empty_object_.key1", Value: "value"},
				{Key: "non_empty_object_.key2.0", Value: "0"},
				{Key: "non_empty_object_.key2.1", Value: "null"},
			},
		},
		{
			name:    "nested array",
			message: `[{}, 1.337, "test", {"key1": {"key2": 3}}]`,
			fields: []*models.FlatMessage_KV{
				{Key: "1", Value: "1.337"},
				{Key: "2", Value: "test"},
				{Key: "3.key1.key2", Value: "3"},
			},
		},
		{
			name:    "sorted keys",
			message: `{"z": 1, "a": "x", "b": false}`,
			fields: []*models.FlatMessage_KV{
				{Key: "a", Value: "x"},
				{Key: "b", Value: "false"},
				{Key: "z", Value: "1"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectFields := tt.fields

			gotFields := flatten([]byte(tt.message)).Fields

			assert.Equal(t, expectFields, gotFields)
		})
	}
}

func Test_Flatten_Negative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "left-over input",
			message: `{"key": "value"} ["extra stuff"]`,
		},
		{
			name:    "invalid JSON start",
			message: `invalid`,
		},
		{
			name:    "invalid object keys",
			message: `{1: 2}`,
		},
		{
			name:    "abrupt object end",
			message: `{"key": "value", "a": 1`,
		},
		{
			name:    "invalid object end delimiter",
			message: `{"key": "value"]`,
		},
		{
			name:    "invalid nested structure",
			message: `{"key": value}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotFields := flatten([]byte(tt.message)).Fields

			// Any error should result in no parsed message fields. This will still allow the log message to
			// be found using a non-scoped search, but will help us avoid incorrect lookups.
			assert.Nil(t, gotFields)
		})
	}
}

func Test_FlatMessageToMapping(t *testing.T) {
	t.Parallel()

	const message = `{"key": "value", "a": 1337, "x": false, "nested": {"a": 1}}`

	// Both the flatMessage and the flatMapping should be reusable concurrently,
	// since only reads are performed once a flatMessage is constructed
	flatMessage := flatten([]byte(message))
	flatMapping := flatMessageToMapping(flatMessage)

	tests := []struct {
		key   string
		value string
		found bool
	}{
		{
			key:   "key",
			value: "value",
			found: true,
		},
		{
			key:   "a",
			value: "1337",
			found: true,
		},
		{
			key:   "x",
			value: "false",
			found: true,
		},
		{
			key:   "nested",
			found: false,
		},
		{
			key:   "nested.a",
			value: "1",
			found: true,
		},
		{
			key:   "nonexistent",
			found: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.key, func(t *testing.T) {
			t.Parallel()

			expectedValue := tt.value

			gotValue, gotOk := flatMapping(tt.key)

			if !tt.found {
				assert.False(t, gotOk)
				return
			}

			assert.True(t, gotOk)
			assert.Equal(t, expectedValue, gotValue)
		})
	}
}
