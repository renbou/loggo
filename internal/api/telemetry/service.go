package telemetry

import (
	"context"
	"time"

	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/telemetry"
)

type messageStore interface {
	ListMessages(from, to time.Time, filter storage.Filter, after []byte, limit uint,
	) (storage.Batch, error)
	StreamMessages(ctx context.Context, from time.Time, filter storage.Filter, limit uint,
	) (storage.Batch, chan storage.Message, error)
}

// Service provides the API for reading various collected telemetry.
type Service struct {
	pb.UnimplementedTelemetryServer
	ms messageStore
}

// NewService initializes a new service using the given message store.
func NewService(ms messageStore) *Service {
	return &Service{
		ms: ms,
	}
}

func (s *Service) ListLogMessages(context.Context, *pb.ListLogMessagesRequest) (*pb.ListLogMessagesResponse, error) {
	return &pb.ListLogMessagesResponse{}, nil
}

func (s *Service) StreamLogMessages(*pb.StreamLogMessagesRequest, pb.Telemetry_StreamLogMessagesServer) error {
	return nil
}
