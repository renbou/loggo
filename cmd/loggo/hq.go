package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/renbou/loggo/internal/api/pigeoneer"
	"github.com/renbou/loggo/internal/api/telemetry"
	"github.com/renbou/loggo/internal/config"
	"github.com/renbou/loggo/internal/logger"
	"github.com/renbou/loggo/internal/storage"
	"github.com/renbou/loggo/internal/web"
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
	hqCmd.Flags().StringVarP(&hqFlags.configPath, "config", "c", "loggo.yaml", "Loggo configuration file path")
	hqCmd.Flags().StringVar(&hqFlags.storagePath, "storage.path", "data/", "Base directory path for log storage (BadgerDB)")
	hqCmd.Flags().StringVar(&hqFlags.grpcAddr, "grpc.addr", ":20081", "Listen address for the gRPC server")
	hqCmd.Flags().StringVar(&hqFlags.webAddr, "web.addr", ":20080", "Listen address for the Web (HTTP) server")

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

	// Register telemetry service first to be served over grpc-web
	telemetryService := telemetry.NewService(db)
	telemetryPb.RegisterTelemetryServer(grpcServer, telemetryService)

	// Need to setup the grpc-web server now to avoid any of the over gRPC services being served over grpc-web
	httpServer := setupHTTP(immutable.Web.Addr, grpcServer)

	// Pigeoneer service works over pure gRPC
	pigeoneerService := pigeoneer.NewService(db)
	pigeoneerPb.RegisterPigeoneerServer(grpcServer, pigeoneerService)

	// If an error arrives on these channels, we will need to shut down the whole program
	grpcCh := startGRPCServer(grpcServer, grpcListener)
	httpCh := startHTTPServer(httpServer)

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
	case err := <-httpCh:
		logger.Errorw("critical error while serving HTTP, will perform shutdown", "error", err)
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

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), gracefulStopTimeout)
		defer cancel()

		_ = httpServer.Shutdown(ctx)
	}()

	shutdownDone := make(chan struct{})
	go func() {
		shutdownWg.Wait()
		close(shutdownDone)
	}()

	// Signal the shutdown to the service, so that all active streams are gracefully ended.
	pigeoneerService.Stop()

	select {
	case <-shutdownDone:
	case <-time.After(gracefulStopTimeout):
	}

	// Now, actually force the shutdown
	grpcServer.Stop()

	<-grpcCh
	<-httpCh
	<-shutdownDone

	return nil
}

func setupBadger(path string) (*storage.Badger, error) {
	// :TODO: add configuration for some Badger options
	return storage.NewBadger(
		badger.DefaultOptions(path).
			WithValueLogFileSize(1 << 27). // 128 MB instead of 1G so that it works on lower-memory machines
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

func setupHTTP(addr string, grpcServer *grpc.Server) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Register all POST grpc-web resources. Avoid doing this using a middleware so
	// that only the telemetry API methods are available.
	grpcWeb := grpcweb.WrapServer(grpcServer)
	for _, route := range grpcweb.ListGRPCResources(grpcServer) {
		engine.POST(route, gin.WrapH(grpcWeb))
	}

	// All static content, including index.html and JS/CSS
	engine.StaticFS("/", web.Content)

	server := &http.Server{
		Addr:        addr,
		Handler:     engine,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 2 * time.Minute,
	}
	return server
}

func startGRPCServer(server *grpc.Server, listener net.Listener) chan error {
	ch := make(chan error)
	go func() {
		ch <- server.Serve(listener)
		close(ch)
	}()

	return ch
}

func startHTTPServer(server *http.Server) chan error {
	ch := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		} else {
			ch <- nil
		}
		close(ch)
	}()

	return ch
}
