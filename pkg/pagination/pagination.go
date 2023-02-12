package pagination

import (
	"encoding/base64"
)

const (
	defaultPageSize = 100
	maxPageSize     = 1000
)

type paginated interface {
	GetPageSize() int32
	GetPageToken() string
}

// CalculateLimit calculates the actual limit to use based on the page size received.
// It is equal to min(pageSize < 1 ? defaultPageSize : pageSize, maxPageSize)
func CalculateLimit(pageSize int32) int32 {
	limit := pageSize
	if limit < 1 {
		limit = defaultPageSize
	} else if limit > maxPageSize {
		limit = maxPageSize
	}
	return limit
}

// DecodePageToken tries to decode the page token using base64, otherwise returning a nil slice.
func DecodePageToken(pageToken string) []byte {
	data, err := base64.RawStdEncoding.DecodeString(pageToken)
	if err != nil {
		return nil
	}
	return data
}

// EncodePageToken encodes a page token using base64.
func EncodePageToken(pageToken []byte) string {
	return base64.RawStdEncoding.EncodeToString(pageToken)
}

// FromRequest calculates the limit and page token from an API request containing page_size and page_token fields.
func FromRequest(request paginated) (int32, []byte) {
	return CalculateLimit(request.GetPageSize()), DecodePageToken(request.GetPageToken())
}
