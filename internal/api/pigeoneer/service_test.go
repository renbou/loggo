package pigeoneer

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/renbou/loggo/internal/storage"
	desc "github.com/renbou/loggo/pkg/api/pigeoneer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bufSize = 1 << 16

var testMessages = []storage.Message{
	storage.Message(`{"valid":"json"}`),
	storage.Message("invalid json message"),
	storage.Message(`{}`),
}

type messageStoreMock struct {
	gotMessages []storage.Message
}

func (s *messageStoreMock) AddMessage(_ time.Time, m storage.Message) error {
	s.gotMessages = append(s.gotMessages, m)
	return nil
}

func setupTestServer(wg *sync.WaitGroup) (*messageStoreMock, *Service, *bufconn.Listener, *grpc.Server) {
	ms := &messageStoreMock{}
	service := NewService(ms)

	l := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	desc.RegisterPigeoneerServer(server, service)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = server.Serve(l)
	}()

	return ms, service, l, server
}

func setupTestClient(ctx context.Context, t *testing.T, listener *bufconn.Listener) (
	*grpc.ClientConn, desc.PigeoneerClient,
) {
	t.Helper()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)

	return conn, desc.NewPigeoneerClient(conn)
}

func Test_Service_Dispatch_ClientStop(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	// Initialize and start server
	store, service, listener, server := setupTestServer(&wg)

	// Connect the client
	ctx := context.Background()
	conn, client := setupTestClient(ctx, t, listener)

	// Send messages via the client, then end the stream on the client's side.
	// Since the server isn't getting closed, no errors should occur anywhere here
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	for _, message := range testMessages {
		err := dispatchClient.Send(&desc.DispatchRequest{
			Timestamp: timestamppb.Now(),
			Message:   message,
		})

		assert.NoError(t, err)
	}

	_, err = dispatchClient.CloseAndRecv()
	assert.NoError(t, err)

	assert.NoError(t, conn.Close())

	// Close the service. Since no streams are running, this should happen instantly
	service.Stop()
	server.Stop()

	// Validate that we received the correct messages
	assert.Equal(t, testMessages, store.gotMessages)

	wg.Wait()
}

func Test_Service_Dispatch_ClientCancel(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	// Initialize and start server
	_, service, listener, server := setupTestServer(&wg)

	// Connect the client
	ctx, cancel := context.WithCancel(context.Background())
	conn, client := setupTestClient(ctx, t, listener)

	// Send single message via client, then cancel. The server should handle the cancel and stop the stream.
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	message := testMessages[0]
	err = dispatchClient.Send(&desc.DispatchRequest{
		Timestamp: timestamppb.Now(),
		Message:   message,
	})
	assert.NoError(t, err)

	cancel()

	_, err = dispatchClient.CloseAndRecv()
	assert.Equal(t, codes.Canceled, status.Code(err))

	assert.NoError(t, conn.Close())

	// Close the service. Since no streams are running, this should happen instantly
	service.Stop()
	server.Stop()

	wg.Wait()
}

func Test_Service_Dispatch_AlreadyStopped(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	// Initialize and start server
	_, service, listener, server := setupTestServer(&wg)

	// Connect the client
	ctx := context.Background()
	conn, client := setupTestClient(ctx, t, listener)

	// Service gets closed, but server hasn't been closed yet
	service.Stop()

	// Client should immediately get an Unavailable error
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	_, err = dispatchClient.CloseAndRecv()
	assert.Equal(t, codes.Unavailable, status.Code(err))

	assert.NoError(t, conn.Close())

	// Finally, close the actual server
	server.Stop()
	wg.Wait()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, goleak.IgnoreTopFunction("github.com/golang/glog.(*loggingT).flushDaemon"))
}
