package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/renbou/loggo/internal/api/pigeoneer"
	"github.com/renbou/loggo/internal/api/telemetry"
	"github.com/renbou/loggo/internal/config"
	"github.com/renbou/loggo/internal/logger"
	"github.com/renbou/loggo/internal/storage"
	pigeoneerPb "github.com/renbou/loggo/pkg/api/pigeoneer"
	telemetryPb "github.com/renbou/loggo/pkg/api/telemetry"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const gracefulStopTimeout = time.Second * 5

var hqFlags struct {
	configPath  string
	storagePath string
	grpcAddr    string
	webAddr     string
}

var hqCmd = &cobra.Command{
	Use:   "hq [flags]",
	Short: "Run the headquarters of Loggo",
	RunE: func(cmd *cobra.Command, _ []string) error {
		immutable := config.ReadImmutable(cmd.LocalFlags())
		mutable, err := config.ReadMutable(hqFlags.configPath)
		if err != nil {
			return err
		}

		return runHQ(&immutable, mutable)
	},
}

func init() {
	hqCmd.LocalFlags().StringVarP(&hqFlags.configPath, "config", "c", "loggo.yaml", "Loggo configuration file path")
	hqCmd.LocalFlags().StringVar(&hqFlags.storagePath, "storage.path", "data/", "Base directory path for log storage (BadgerDB)")
	hqCmd.LocalFlags().StringVar(&hqFlags.grpcAddr, "grpc.addr", ":20081", "Listen address for the gRPC server")
	hqCmd.LocalFlags().StringVar(&hqFlags.webAddr, "web.addr", ":20080", "Listen address for the Web (HTTP) server")

	rootCmd.AddCommand(hqCmd)
}

func runHQ(immutable *config.Immutable, mutable *config.Mutable) error {
	// Open the DB first and instantly setup a defer with Close() for all our writes to be persisted no matter what
	db, err := setupBadger(immutable.Storage.Path)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorw("error while closing the database", "error", err)
		}
	}()

	// Initialize logger after DB so that all future log messages are persisted to it
	logger.InitGlobal(func(t time.Time, message []byte) error {
		return db.AddMessage(t, message)
	}, "stderr")

	grpcServer, grpcListener, err := setupGRPC(immutable.GRPC.Addr)
	if err != nil {
		return err
	}

	// Pigeoneer service works over pure gRPC
	pigeoneerService := pigeoneer.NewService(db)
	pigeoneerPb.RegisterPigeoneerServer(grpcServer, pigeoneerService)

	// While the telemetry service works over gRPC-web for user access from the frontend
	telemetryService := telemetry.NewService(db)
	telemetryPb.RegisterTelemetryServer(grpcServer, telemetryService)

	// If an error arrives on this channel, we will need to shut down the whole program
	grpcCh := startGRPCServer(grpcServer, grpcListener)

	logger.Infow("all components initialized and started",
		"config_file", hqFlags.configPath,
		"storage_path", immutable.Storage.Path,
		"grpc_addr", immutable.GRPC.Addr,
		"web_addr", immutable.Web.Addr,
	)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-grpcCh:
		// Perform an emergency shutdown
		logger.Errorw("critical error while serving gRPC, will perform shutdown", "error", err)
	case <-exitCh:
		// Perform a normal shutdown
		logger.Infow("gracefully shutting down all components")
	}

	// Launch graceful shutdowns and wait for some time
	var shutdownWg sync.WaitGroup
	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		grpcServer.GracefulStop()
	}()

	shutdownDone := make(chan struct{})
	go func() {
		shutdownWg.Wait()
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
	case <-time.After(gracefulStopTimeout):
	}

	// Now, actually force the shutdowns
	pigeoneerService.Stop()
	grpcServer.Stop()

	<-grpcCh
	<-shutdownDone

	return nil
}

func setupBadger(path string) (*storage.Badger, error) {
	// :TODO: add configuration for some Badger options
	return storage.NewBadger(
		badger.DefaultOptions(path).
			WithLogger(nil),
	)
}

func setupGRPC(addr string) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("starting listener for gRPC on %s: %w", addr, err)
	}

	// :TODO: does this need configuring?
	server := grpc.NewServer()
	reflection.Register(server)

	return server, listener, nil
}

func startGRPCServer(server *grpc.Server, listener net.Listener) chan error {
	ch := make(chan error)
	go func() {
		ch <- server.Serve(listener)
		close(ch)
	}()

	return ch
}
