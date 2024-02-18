package telemetry

import (
	"context"
	"time"

	"github.com/renbou/loggo/internal/logger"
	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/telemetry"
	"github.com/renbou/loggo/pkg/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// ListLogMessages provides a paginated list of log messages in the specified interval,
// filtered by the optional LogFilter.
func (s *Service) ListLogMessages(ctx context.Context, req *pb.ListLogMessagesRequest) (*pb.ListLogMessagesResponse, error) {
	limit, after := pagination.FromRequest(req)

	batch, err := s.ms.ListMessages(req.From.AsTime(), req.To.AsTime(), compileFilter(req.Filter), after, uint(limit))
	if err != nil {
		logger.Errorw("error while listing messages in storage",
			"component", "api.telemetry",
			"filter", req.Filter.String(),
			"error", err,
		)
		return nil, status.Error(codes.Internal, "storage error")
	}

	return &pb.ListLogMessagesResponse{Batch: s.batchToPB(batch)}, nil
}

// StreamLogMessages runs a long-running stream over which new log messages are sent.
// The first thing sent over the stream, however, will be a batch of the latest suitable messages. If it contains
// more elements than the limit, it will contain a next_page_token suitable for use with ListLogMessages.
func (s *Service) StreamLogMessages(req *pb.StreamLogMessagesRequest, stream pb.Telemetry_StreamLogMessagesServer) error {
	limit := pagination.CalculateLimit(req.GetPageSize())

	batch, ch, err := s.ms.StreamMessages(stream.Context(), req.From.AsTime(), compileFilter(req.Filter), uint(limit))
	if err != nil {
		logger.Errorw("error while streaming messages from storage",
			"component", "api.telemetry",
			"filter", req.Filter.String(),
			"error", err,
		)
		return status.Error(codes.Internal, "storage error")
	}

	if err = stream.Send(&pb.StreamLogMessagesResponse{
		Response: &pb.StreamLogMessagesResponse_Batch{Batch: s.batchToPB(batch)},
	}); err != nil {
		logger.Errorw("error while sending first batch during stream",
			"component", "api.telemetry",
			"error", err,
		)
		return status.Error(codes.Internal, "sending batch")
	}

	for m := range ch {
		if err = stream.Send(&pb.StreamLogMessagesResponse{
			Response: &pb.StreamLogMessagesResponse_Message{Message: m},
		}); err != nil {
			logger.Errorw("error while sending message during stream",
				"component", "api.telemetry",
				"error", err,
			)
			return status.Error(codes.Internal, "sending message")
		}
	}

	return nil
}

func (s *Service) batchToPB(batch storage.Batch) *pb.LogBatch {
	pbBatch := &pb.LogBatch{
		Messages:      make([]*pb.LogMessage, 0, len(batch.Messages)),
		NextPageToken: pagination.EncodePageToken(batch.Next),
	}

	for _, message := range batch.Messages {
		pbBatch.Messages = append(pbBatch.Messages, &pb.LogMessage{
			Message: message.M,
			Id:      message.ID,
		})
	}
	return pbBatch
}
