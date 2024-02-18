package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/renbou/loggo/internal/api/telemetry/mocks"
	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/telemetry"
	"github.com/renbou/loggo/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func setupTestService(t *testing.T) (*Service, *mocks.MessageStoreMock) {
	mockStore := mocks.NewMessageStoreMock(t)
	service := NewService(mockStore)

	return service, mockStore
}

func Test_Service_ListLogMessages_Works(t *testing.T) {
	t.Parallel()

	message := storage.Message(`{"message": "example"}`)
	fm := flatMappingFromMap(map[string]string{"message": "example"})
	after := []byte("test")
	limit := 13

	service, mockStore := setupTestService(t)
	mockStore.ListMessagesMock.Set(func(from, to time.Time, filter storage.Filter, gotAfter []byte, gotLimit uint,
	) (storage.Batch, error) {
		assert.NotZero(t, from)
		assert.NotZero(t, to)
		assert.Equal(t, after, gotAfter)
		assert.EqualValues(t, limit, gotLimit)
		assert.True(t, filter(message, fm))
		return storage.Batch{Messages: []storage.StoredMessage{{M: message, ID: nil}}, Next: []byte("test")}, nil
	})

	expectResp := &pb.ListLogMessagesResponse{
		Batch: &pb.LogBatch{Messages: []*pb.LogMessage{{Message: message, Id: nil}}, NextPageToken: "dGVzdA"},
	}

	gotResp, err := service.ListLogMessages(context.Background(), &pb.ListLogMessagesRequest{
		From: timestamppb.Now(),
		To:   timestamppb.Now(),
		Filter: &pb.LogFilter{Filter: &pb.LogFilter_Scoped_{Scoped: &pb.LogFilter_Scoped{
			Field: "message",
			Value: "example",
		}}},
		PageSize:  int32(limit),
		PageToken: pagination.EncodePageToken(after),
	})

	assert.NoError(t, err)
	assert.Equal(t, expectResp.Batch.Messages, gotResp.Batch.Messages)
	assert.Equal(t, expectResp.Batch.NextPageToken, gotResp.Batch.NextPageToken)
}

func Test_Service_ListLogMessages_StoreError(t *testing.T) {
	t.Parallel()

	service, mockStore := setupTestService(t)
	mockStore.ListMessagesMock.Return(storage.Batch{}, assert.AnError)

	_, err := service.ListLogMessages(context.Background(), &pb.ListLogMessagesRequest{
		From:   timestamppb.Now(),
		To:     timestamppb.Now(),
		Filter: nil,
	})

	assert.Error(t, err)
}

func Test_Service_StreamLogMessages_Works(t *testing.T) {
	t.Parallel()

	messages := []storage.Message{
		storage.Message(`{"message": "example"}`),
		storage.Message(`{"also": "example"}`),
	}
	fm := flatMappingFromMap(map[string]string{"message": "example"})
	limit := 5

	service, mockStore := setupTestService(t)
	mockStore.StreamMessagesMock.Set(func(ctx context.Context, from time.Time, filter storage.Filter, gotLimit uint,
	) (storage.Batch, chan storage.Message, error) {
		assert.NotZero(t, from)
		assert.EqualValues(t, limit, gotLimit)
		assert.True(t, filter(messages[0], fm))

		// Simulate a single message being sent while the context is alive
		assert.True(t, filter(messages[1], fm))
		ch := make(chan storage.Message, 1)
		ch <- messages[1]
		close(ch)

		return storage.Batch{Messages: []storage.StoredMessage{{M: messages[0], ID: nil}}}, ch, nil
	})

	var sentN int
	mockStream := mocks.NewTelemetry_StreamLogMessagesServerMock(t)
	mockStream.SendMock.Set(func(sp1 *pb.StreamLogMessagesResponse) (err error) {
		if sentN == 0 {
			// First response must always be a batch
			batch, ok := sp1.Response.(*pb.StreamLogMessagesResponse_Batch)
			assert.True(t, ok)
			assert.Equal(t, []*pb.LogMessage{{Message: messages[0], Id: nil}}, batch.Batch.Messages)
			sentN++
			return nil
		}

		// Next, we expect a single response with the second message from the channel
		assert.Equal(t, 1, sentN)

		message, ok := sp1.Response.(*pb.StreamLogMessagesResponse_Message)
		assert.True(t, ok)
		assert.EqualValues(t, messages[1], message.Message)
		return nil
	})
	mockStream.ContextMock.Expect().Return(context.Background())

	err := service.StreamLogMessages(&pb.StreamLogMessagesRequest{
		From:     timestamppb.Now(),
		Filter:   &pb.LogFilter{Filter: &pb.LogFilter_Text_{Text: &pb.LogFilter_Text{Value: "example"}}},
		PageSize: int32(limit),
	}, mockStream)
	assert.NoError(t, err)
}

func Test_Service_StreamLogMessages_StoreError(t *testing.T) {
	t.Parallel()

	service, mockStore := setupTestService(t)
	mockStore.StreamMessagesMock.Return(storage.Batch{}, nil, assert.AnError)

	mockStream := mocks.NewTelemetry_StreamLogMessagesServerMock(t)
	mockStream.ContextMock.Expect().Return(context.Background())

	err := service.StreamLogMessages(&pb.StreamLogMessagesRequest{
		From:   timestamppb.Now(),
		Filter: nil,
	}, mockStream)

	assert.Error(t, err)
}

func Test_Service_StreamLogMessages_SendBatchError(t *testing.T) {
	t.Parallel()

	service, mockStore := setupTestService(t)
	mockStore.StreamMessagesMock.Return(storage.Batch{}, make(chan storage.Message), nil)

	mockStream := mocks.NewTelemetry_StreamLogMessagesServerMock(t)
	mockStream.ContextMock.Expect().Return(context.Background())
	mockStream.SendMock.Return(assert.AnError)

	err := service.StreamLogMessages(&pb.StreamLogMessagesRequest{
		From:   timestamppb.Now(),
		Filter: nil,
	}, mockStream)
	assert.Error(t, err)
}

func Test_Service_StreamLogMessages_SendMessageError(t *testing.T) {
	t.Parallel()

	ch := make(chan storage.Message, 1)
	ch <- storage.Message{}
	close(ch)

	service, mockStore := setupTestService(t)
	mockStore.StreamMessagesMock.Return(storage.Batch{}, ch, nil)

	mockStream := mocks.NewTelemetry_StreamLogMessagesServerMock(t)
	mockStream.ContextMock.Expect().Return(context.Background())

	var sentN int
	mockStream.SendMock.Set(func(_ *pb.StreamLogMessagesResponse) (err error) {
		if sentN == 0 {
			sentN++
			return nil
		}
		return assert.AnError
	})

	err := service.StreamLogMessages(&pb.StreamLogMessagesRequest{
		From:   timestamppb.Now(),
		Filter: nil,
	}, mockStream)
	assert.Error(t, err)
}
