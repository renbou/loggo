package storage

// Message is a single log message. Log messages *SHOULD* be JSON-formatted to enable field-specific search,
// however it isn't enforced, and non-field searches are still possible.
type Message []byte

// Batch is a paginated batch of messages. Next is nil if there are no more messages to be retrieved,
// otherwise it is an internal key of the next item to be retrieved.
type Batch struct {
	Messages []Message
	Next     []byte
}

// FlatMapping describes an unnested map[string]any, allowing for fast access.
type FlatMapping func(key string) (value string, ok bool)

// Filter can (optionally) be passed to List* methods of the storage to filter out the logs which are returned.
// Note: Filters should be thread-safe since they can be called from multiple goroutines at once.
type Filter func(Message, FlatMapping) bool
