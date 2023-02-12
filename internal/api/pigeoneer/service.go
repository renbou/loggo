package pigeoneer

import (
	"io"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/renbou/loggo/internal/logger"
	"github.com/renbou/loggo/internal/mw"
	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/pigeoneer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var dispatchedLogMessagesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "loggo",
	Subsystem: "pigeoneer",
	Name:      "dispatched_log_messages_total",
	Help:      "Number of log messages dispatched through the Pigeoneer gRPC service.",
}, []string{"pigeon"})

type messageStore interface {
	AddMessage(t time.Time, m storage.Message) error
}

type dispatchRequestWrapper struct {
	message *pb.DispatchRequest
	err     error
}

// Service is the pigeoneer service providing dispatch functionality for log messages via pigeons.
type Service struct {
	pb.UnimplementedPigeoneerServer
	ms messageStore

	done chan struct{}
	mu   sync.RWMutex
	wg   sync.WaitGroup
}

// NewService initializes a new service using the given message store.
func NewService(ms messageStore) *Service {
	return &Service{
		ms:   ms,
		done: make(chan struct{}),
	}
}

// Stop gracefully stops all of the currently active streams, waiting for them to finish.
// This is needed instead of a simple gRPC graceful stop to avoid losing any log messages.
func (s *Service) Stop() {
	// Lock needed to synchronize Add and Wait on the WaitGroup during the moment when we feel like the service
	// is already closed, but some new RPCs have popped up
	s.mu.Lock()
	defer s.mu.Unlock()

	close(s.done)
	s.wg.Wait()
}

// Dispatch dispatches incoming messages from the stream to the storage.
// All of the complex synchronization here is needed during the shutdown of the server, allowing pending messages
// to be persisted to storage, and new ones to not be accepted. If a simple grpcServer.Stop() is used, we might
// fail to send an "ACK" back to the client, resulting in duplicates.
// Note: there might still be duplicates if an ACK was failed to be sent,
// however that should happen only in cases when the actuall network connection has been broken.
func (s *Service) Dispatch(stream pb.Pigeoneer_DispatchServer) error {
	serverClosed := status.Error(codes.Unavailable, "pigeoneer stopping")
	pigeonName := mw.PigeonNameFromCtx(stream.Context())

	// If Stop is running now, this will fail, and we'll simply exit in the next if
	// Otherwise, this will successfully add a new active stream to be handled by Stop
	s.mu.RLock()
	s.wg.Add(1)
	defer s.wg.Done()
	s.mu.RUnlock()

	// Check this before actually adding to the WaitGroup. If this check isn't made, then we might get a message
	// first, instead of the s.done signal, since the order in select isn't defined when multiple channels are ready.
	if s.isClosed() {
		return serverClosed
	}

	var request dispatchRequestWrapper
	ch := s.pipeMessages(stream)
	for {
		select {
		case request = <-ch:
		case <-s.done:
			return serverClosed
		}

		// Properly handle client closing the connection
		message, err := request.message, request.err
		if err != nil {
			if err == io.EOF || status.Convert(err).Code() == codes.Canceled {
				return nil
			}

			logger.Errorw("unexpected error while receiving new message",
				"component", "pigeoneer",
				"error", err,
			)
			return status.Error(codes.Internal, "receive error")
		}

		// No errors should occur here during normal execution,
		// since we always stop for this service to finish before closing the storage
		if err := s.ms.AddMessage(message.Timestamp.AsTime(), message.Message); err != nil {
			logger.Errorw("failed to add message to storage during dispatch",
				"component", "pigeoneer",
				"message_timestamp", message.Timestamp.AsTime(),
				"message_data", message.Message,
				"error", err,
			)
			return status.Error(codes.Internal, "storage error")
		}

		dispatchedLogMessagesTotal.WithLabelValues(pigeonName).Inc()

		if err := stream.Send(&emptypb.Empty{}); err != nil {
			logger.Errorw("unexpected error while acking written message",
				"component", "pigeoneer",
				"error", err,
			)
			return status.Error(codes.Internal, "send ack error")
		}
	}
}

// pipeMessages runs in a separate goroutine so that Dispatch is never blocked during Recv
func (s *Service) pipeMessages(stream pb.Pigeoneer_DispatchServer) chan dispatchRequestWrapper {
	ch := make(chan dispatchRequestWrapper)

	go func() {
		for {
			message, err := stream.Recv()

			select {
			// Once the done channel is closed, noone will be listening on the channel, and the stream will already be closed
			case ch <- dispatchRequestWrapper{message, err}:
				if err == nil {
					continue
				}
				// Error will be handled by Dispatch, we also need to exit to avoid leaking the goroutine
			case <-s.done:
			}

			// close the channel here, since we know no more writes will happen
			close(ch)
			return
		}
	}()

	return ch
}

func (s *Service) isClosed() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}
