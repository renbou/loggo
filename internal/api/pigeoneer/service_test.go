package pigeoneer

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/renbou/loggo/internal/storage"
	pb "github.com/renbou/loggo/pkg/api/pigeoneer"
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
	mu          sync.Mutex
	gotMessages []storage.Message
}

func (s *messageStoreMock) AddMessage(_ time.Time, m storage.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gotMessages = append(s.gotMessages, m)
	return nil
}

func (s *messageStoreMock) WaitFor(n int) {
	for {
		s.mu.Lock()
		l := len(s.gotMessages)
		s.mu.Unlock()

		if l == n {
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}

type multiBufconnListener struct {
	listener atomic.Pointer[bufconn.Listener]
	mu       sync.RWMutex
}

func (l *multiBufconnListener) Accept() (net.Conn, error) {
	return l.listener.Load().Accept()
}

func (l *multiBufconnListener) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	_ = l.listener.Load().Close()
	l.listener.Store(bufconn.Listen(bufSize))
	return nil
}

func (l *multiBufconnListener) Addr() net.Addr {
	return l.listener.Load().Addr()
}

func (l *multiBufconnListener) Dial() (net.Conn, error) {
	return l.listener.Load().Dial()
}

func (l *multiBufconnListener) finalClose() error {
	return l.listener.Load().Close()
}

func setupTestServer(wg *sync.WaitGroup, listener net.Listener, service *Service) *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterPigeoneerServer(server, service)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = server.Serve(listener)
	}()

	return server
}

func setupTestService(wg *sync.WaitGroup) (*messageStoreMock, *Service, *multiBufconnListener, *grpc.Server) {
	ms := &messageStoreMock{}
	service := NewService(ms)

	listener := &multiBufconnListener{}
	listener.listener.Store(bufconn.Listen(bufSize))

	server := setupTestServer(wg, listener, service)
	return ms, service, listener, server
}

func setupTestClient(ctx context.Context, t *testing.T, listener *multiBufconnListener) (
	*grpc.ClientConn, pb.PigeoneerClient,
) {
	t.Helper()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)

	return conn, pb.NewPigeoneerClient(conn)
}

func Test_Service_Dispatch_ClientStop(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	// Initialize and start server
	store, service, listener, server := setupTestService(&wg)
	defer listener.finalClose()

	// Connect the client
	ctx := context.Background()
	conn, client := setupTestClient(ctx, t, listener)

	// Send messages via the client, then end the stream on the client's side.
	// Since the server isn't getting closed, no errors should occur anywhere here
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	for _, message := range testMessages {
		err := dispatchClient.Send(&pb.DispatchRequest{
			Timestamp: timestamppb.Now(),
			Message:   message,
		})
		assert.NoError(t, err)

		_, err = dispatchClient.Recv()
		assert.NoError(t, err)
	}

	err = dispatchClient.CloseSend()
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
	_, service, listener, server := setupTestService(&wg)
	defer listener.finalClose()

	// Connect the client
	ctx, cancel := context.WithCancel(context.Background())
	conn, client := setupTestClient(ctx, t, listener)

	// Send single message via client, then cancel. The server should handle the cancel and stop the stream.
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	message := testMessages[0]
	err = dispatchClient.Send(&pb.DispatchRequest{
		Timestamp: timestamppb.Now(),
		Message:   message,
	})
	assert.NoError(t, err)

	cancel()

	_, err = dispatchClient.Recv()
	assert.Equal(t, codes.Canceled, status.Code(err))

	assert.NoError(t, dispatchClient.CloseSend())
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
	_, service, listener, server := setupTestService(&wg)
	defer listener.finalClose()

	// Connect the client
	ctx := context.Background()
	conn, client := setupTestClient(ctx, t, listener)

	// Service gets closed, but server hasn't been closed yet
	service.Stop()

	// Client should get an Unavailable error as the response
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	_, err = dispatchClient.Recv()
	assert.Equal(t, codes.Unavailable, status.Code(err))

	assert.NoError(t, dispatchClient.CloseSend())
	assert.NoError(t, conn.Close())

	// Finally, close the actual server
	server.Stop()
	wg.Wait()
}

func Test_Service_Dispatch_StopDuringStream(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	// Initialize and start server
	_, service, listener, server := setupTestService(&wg)
	defer listener.finalClose()

	// Connect the client
	ctx := context.Background()
	conn, client := setupTestClient(ctx, t, listener)

	// Client connects now, sends a message
	dispatchClient, err := client.Dispatch(ctx)
	require.NoError(t, err)

	request := &pb.DispatchRequest{
		Timestamp: timestamppb.Now(),
		Message:   testMessages[0],
	}
	err = dispatchClient.Send(request)
	assert.NoError(t, err)

	// Then, the service gets stopped
	service.Stop()

	// Client should get an Unavailable error as the response
	_, err = dispatchClient.Recv()
	assert.Equal(t, codes.Unavailable, status.Code(err))

	// Close the server and the connection
	assert.NoError(t, dispatchClient.CloseSend())
	assert.NoError(t, conn.Close())
	server.Stop()
	wg.Wait()
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, goleak.IgnoreTopFunction("github.com/golang/glog.(*fileSink).flushDaemon"))
}
