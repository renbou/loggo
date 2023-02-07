package logger

import (
	"bytes"
	"fmt"
	"time"
)

// consumerSink is a special sink for zap which routes log messages to the specified consumer.
// Log messages are buffered and sent to the consumer once a line ending is encountered.
type consumerSink struct {
	consumer   LogConsumer
	lineEnding []byte
	buffer     []byte
}

func (s *consumerSink) Write(b []byte) (int, error) {
	s.buffer = append(s.buffer, b...)

	i := bytes.Index(s.buffer, s.lineEnding)
	if i == -1 {
		return len(b), nil
	}

	message := s.buffer[:i]
	if err := s.send(message); err != nil {
		return len(b), err
	}

	s.buffer = s.buffer[:copy(s.buffer, s.buffer[i+1:])]
	return len(b), nil
}

func (s *consumerSink) Sync() error {
	if len(s.buffer) == 0 {
		return nil
	}

	if err := s.send(s.buffer); err != nil {
		return err
	}

	s.buffer = s.buffer[:0]
	return nil
}

func (s *consumerSink) Close() error {
	return s.Sync()
}

func (s *consumerSink) send(m []byte) error {
	// Allows lazy-initializing the consumer, which simplifies tests
	if s.consumer == nil {
		return nil
	}

	if err := s.consumer(time.Now(), m); err != nil {
		return fmt.Errorf("sending log message to consumer (storage): %w", err)
	}
	return nil
}
