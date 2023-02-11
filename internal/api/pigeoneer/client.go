package pigeoneer

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/renbou/loggo/internal/storage"
	desc "github.com/renbou/loggo/pkg/api/pigeoneer"
	"github.com/renbou/loggo/pkg/evictch"
)

// Parameters for retrying stream connections if they fail.
// :TODO: does this need configuring?
const (
	jitterCoef = 0.1
	backoffExp = 2
	backoffMin = time.Millisecond * 50
	backoffMax = time.Second * 5
)

// Client provides an API to access the pigeoneer service with automatic buffering of messages and retries.
type Client struct {
	grpc      desc.PigeoneerClient
	queue     *evictch.Chan[*desc.DispatchRequest]
	jitterRnd *rand.Rand

	wg     sync.WaitGroup
	cancel context.CancelFunc
}

// NewClient initializes a new pigeoneer client using the given connection.
// Up to bufferCapacity requests will be buffered in memory, the oldest one being evicted, if needed.
// Actual dispatch of messages will happen in a background goroutine with automatic reconnects and retries.
func NewClient(cc grpc.ClientConnInterface, bufferCapacity uint) *Client {
	client := &Client{
		grpc:      desc.NewPigeoneerClient(cc),
		queue:     evictch.NewChan[*desc.DispatchRequest](bufferCapacity),
		jitterRnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client.cancel = cancel
	client.wg.Add(1)
	go client.runDispatches(ctx)

	return client
}

// Stop gracefully stops the client by closing the buffer and waiting for the dispatch runner to exit.
// Note that new messages will still be accepted (due to the usage of evictch and its implementation), but they won't be sent out.
func (c *Client) Stop() {
	c.cancel()
	c.queue.Close()
	c.wg.Wait()
}

// Dispatch adds a new request to the dispatch queue, performing the actual dispatch in the background runner.
func (c *Client) Dispatch(t time.Time, m storage.Message) {
	c.queue.Write(&desc.DispatchRequest{Timestamp: timestamppb.New(t), Message: m})
}

func (c *Client) runDispatches(ctx context.Context) {
	defer c.wg.Done()

	stream := c.connect(ctx)
	defer func() {
		if stream != nil {
			_ = stream.CloseSend()
		}
	}()

	var request *desc.DispatchRequest
	for {
		// Stream might be nil even on the first iteration if the context is canceled before we managed to connect
		if stream == nil {
			return
		}

		// If we've successfully sent the previous request, get a new one
		if request == nil {
			var ok bool
			request, ok = c.queue.Read()
			if !ok {
				// Channel has been closed, meaning we'll get no more reads
				return
			}
		}

		// Send the request. This doesn't guarantee that it will actually be written,
		// since gRPC send doesn't wait for anything.
		sendErr := stream.Send(request)
		// This, on the other hand, guarantees that the message was successfully persisted by the service
		_, recvErr := stream.Recv()

		if sendErr == nil && recvErr == nil {
			// Prepare for next request
			request = nil
			continue
		}

		// The error might be the context getting closed, in which case we should quit
		if c.shouldStop(ctx, sendErr) || c.shouldStop(ctx, recvErr) {
			return
		}

		// Otherwise, try to reconnect and resend the message
		stream = c.connect(ctx)
	}
}

// connect connects to the service using retries with backoff
// :TODO: maybe this needs logging? but in debug mode, for example...
func (c *Client) connect(ctx context.Context) desc.Pigeoneer_DispatchClient {
	stream, err := c.grpc.Dispatch(ctx)

	var wait time.Duration
	nextWait := backoffMin
	for !(err == nil || c.shouldStop(ctx, err)) {
		wait, nextWait = c.nextConnectWait(nextWait)
		time.Sleep(wait)

		stream, err = c.grpc.Dispatch(ctx)
	}

	return stream
}

func (c *Client) nextConnectWait(prev time.Duration) (wait, next time.Duration) {
	next = prev * backoffExp
	if next > backoffMax {
		next = backoffMax
	}

	// jitter in (-1, 1), scale it down to needed amount, and then scale wait using it
	jitter := c.jitterRnd.Float64()*2 - 1
	jitter *= jitterCoef

	wait = prev + time.Duration(float64(prev)*jitter)
	return wait, next
}

func (c *Client) shouldStop(ctx context.Context, err error) bool {
	code := status.Code(err)
	return code == codes.Canceled || code == codes.DeadlineExceeded
}
