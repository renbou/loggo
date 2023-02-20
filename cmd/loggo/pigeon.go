package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/renbou/loggo/internal/api/pigeoneer"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var pigeonFlags struct {
	hqAddr         string
	bufferCapacity int
}

var pigeonCmd = &cobra.Command{
	Use:   "pigeon [flags]",
	Short: "Run a single pigeon to deliver logs to the HQ",
	RunE: func(cmd *cobra.Command, args []string) error {
		// :TODO: read in config using viper, just like in HQ. This should probably refactored using generics or something.
		return runPigeon(cmd.Context())
	},
}

func init() {
	pigeonCmd.Flags().StringVar(&pigeonFlags.hqAddr, "hq.addr", "localhost:20081",
		"Address of the HQ gRPC server to connect to")
	pigeonCmd.Flags().IntVar(&pigeonFlags.bufferCapacity, "buffer.capacity", 500,
		"Maximum number messages to buffer while HQ is unavailable")

	rootCmd.AddCommand(pigeonCmd)
}

func runPigeon(ctx context.Context) error {
	// :TODO: add authorization using tokens
	cc, err := grpc.DialContext(ctx, pigeonFlags.hqAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dialing HQ via gRPC: %w", err)
	}

	// Start the pigeoneer client. All dispatches are handled by it in the background, all that is left for us to do
	// is to read in logs from stdin and send them off.
	pigeoneerClient := pigeoneer.NewClient(cc, uint(pigeonFlags.bufferCapacity))

	var readErr error
	lineCh := make(chan []byte)
	go func() {
		defer close(lineCh)

		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadBytes('\n')

			if len(line) > 0 {
				lineCh <- line
			}

			// An EOF occurs once the actual program is closed and so have its stdout/stderr channels.
			if errors.Is(err, io.EOF) {
				return
			} else if err != nil {
				readErr = err
				return
			}
		}
	}()

	// :TODO: move graceful handling to an internal library.. rn its copy+pasted between pigeon and hq
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case line, ok := <-lineCh:
			if ok {
				pigeoneerClient.Dispatch(time.Now(), line)
				continue
			}
		case <-exitCh:
		}

		break
	}

	// Wait for the dispatching to stop. This happens near-instantly and is mostly needed to avoid
	// losing some log messages which haven't been fully sent yet.
	// :TODO: this doesn't actually wait for the previously bufferred messages to be sent...
	pigeoneerClient.Stop()
	return readErr
}
