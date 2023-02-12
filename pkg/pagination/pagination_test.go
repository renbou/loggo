package pagination

import (
	"testing"

	"github.com/renbou/loggo/pkg/pagination/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCalculateLimit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pageSize   int32
		calculated int32
	}{
		{name: "default", pageSize: 0, calculated: defaultPageSize},
		{name: "negative resets to default", pageSize: -1, calculated: defaultPageSize},
		{name: "limited to max", pageSize: 100000, calculated: maxPageSize},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := CalculateLimit(tt.pageSize)

			assert.Equal(t, tt.calculated, got)
		})
	}
}

func TestDecodePageToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pageToken string
		value     []byte
	}{
		{name: "empty page token", pageToken: "", value: []byte{}},
		{name: "base64 corrupt page token", pageToken: "-=+%^", value: nil},
		{name: "valid page token", pageToken: "dGVzdA", value: []byte("test")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := DecodePageToken(tt.pageToken)

			assert.Equal(t, tt.value, got)
		})
	}
}

func TestEncodePageToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   []byte
		encoded string
	}{
		{name: "nil token", value: nil, encoded: ""},
		{name: "empty token", value: []byte(""), encoded: ""},
		{name: "non-empty token token", value: []byte("test"), encoded: "dGVzdA"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := EncodePageToken(tt.value)

			assert.Equal(t, tt.encoded, got)
		})
	}
}

func TestFromRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request *mocks.PaginatedMock
		limit   int32
		token   []byte
	}{
		{
			name: "default",
			request: mocks.NewPaginatedMock(t).
				GetPageSizeMock.Return(0).
				GetPageTokenMock.Return(""),
			limit: defaultPageSize,
			token: []byte{},
		},
		{
			name: "non-empty page_size, empty page_token",
			request: mocks.NewPaginatedMock(t).
				GetPageSizeMock.Return(1).
				GetPageTokenMock.Return(""),
			limit: 1,
			token: []byte{},
		},
		{
			name: "empty page_size, non-empty page_token",
			request: mocks.NewPaginatedMock(t).
				GetPageSizeMock.Return(0).
				GetPageTokenMock.Return("dGVzdA"),
			limit: defaultPageSize,
			token: []byte("test"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotLimit, gotToken := FromRequest(tt.request)

			assert.Equal(t, tt.limit, gotLimit)
			assert.Equal(t, tt.token, gotToken)
		})
	}
}
