package pigeoneer

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testBufferSize = 10

func Test_Client(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	// Initialize and start server
	store, service, listener, server := setupTestService(&wg)
	defer listener.finalClose()

	// Connect the client
	ctx := context.Background()
	conn, _ := setupTestClient(ctx, t, listener)
	client := NewClient(conn, testBufferSize)

	// Dispatch the first two messages
	client.Dispatch(time.Now(), testMessages[0])
	client.Dispatch(time.Now(), testMessages[1])

	// Wait for messages to arrive
	store.WaitFor(2)
	assert.Equal(t, testMessages[:2], store.gotMessages)

	// Now, close the service, but not the server. Sending a message from the client now should return Unavailable,
	// and make the client go into retry mode
	service.Stop()
	client.Dispatch(time.Now(), testMessages[2])
	time.Sleep(time.Millisecond * 20)

	// Now, stop the actual server
	server.Stop()
	wg.Wait()

	// Simulate a delay before a new server launches
	time.Sleep(time.Millisecond * 20)
	service = NewService(store)
	server = setupTestServer(&wg, listener, service)

	// Wait for the new stream to start and message to go through
	store.WaitFor(3)

	// Finally, stop the client, then the server
	client.Stop()
	assert.NoError(t, conn.Close())

	service.Stop()
	server.Stop()
	wg.Wait()
}

func Test_Client_NextConnectWait(t *testing.T) {
	t.Parallel()

	expectedWaits := []time.Duration{
		time.Millisecond * 50,
		time.Millisecond * 100,
		time.Millisecond * 200,
		time.Millisecond * 400,
		time.Millisecond * 800,
		time.Millisecond * 1600,
		time.Millisecond * 3200,
		time.Millisecond * 5000,
		time.Millisecond * 5000,
	}

	client := &Client{jitterRnd: rand.New(rand.NewSource(time.Now().UnixNano()))}

	var gotWaits []time.Duration
	var wait time.Duration
	nextWait := backoffMin
	for range expectedWaits {
		wait, nextWait = client.nextConnectWait(nextWait)
		gotWaits = append(gotWaits, wait)
	}

	// Since we have jitter, each result should be close to expected, but not equal
	for i := range expectedWaits {
		assert.InDelta(t, expectedWaits[i], gotWaits[i], float64(expectedWaits[i])*jitterCoef)
	}
}
