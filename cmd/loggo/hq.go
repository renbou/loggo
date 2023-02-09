package main

import (
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/renbou/loggo/internal/config"
	"github.com/renbou/loggo/internal/logger"
	"github.com/renbou/loggo/internal/storage"
	"github.com/spf13/cobra"
)

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
	// :TODO: add configuration for some Badger options
	db, err := storage.NewBadger(
		badger.DefaultOptions(immutable.Storage.Path).
			WithLogger(nil),
	)
	if err != nil {
		return err
	}

	logger.InitGlobal(func(t time.Time, message []byte) error {
		return db.AddMessage(t, message)
	}, "stderr")

	logger.Infow("all components initialized and started")

	return nil
}
