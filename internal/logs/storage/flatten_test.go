package storage

import (
	"testing"

	"github.com/renbou/obzerva/internal/logs/storage/models"
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectFields := tt.fields

			gotFields := flatten(tt.message).Fields

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

			gotFields := flatten(tt.message).Fields

			// Any error should result in no parsed message fields. This will still allow the log message to
			// be found using a non-scoped search, but will help us avoid incorrect lookups.
			assert.Nil(t, gotFields)
		})
	}
}
