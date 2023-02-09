package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// StorageConfig contains the storage configuration.
type StorageConfig struct {
	Path string
}

// GRPCConfig contains configuration for the gRPC server.
type GRPCConfig struct {
	Addr string
}

// WebConfig contains configuration for the web (HTTP) server.
type WebConfig struct {
	Addr string
}

// Immutable contains immutable settings which are set only once at launch time.
type Immutable struct {
	Storage StorageConfig
	GRPC    GRPCConfig
	Web     WebConfig
}

// ReadImmutable reads the immutable config from the specified flags and environment variables.
func ReadImmutable(pflags *pflag.FlagSet) Immutable {
	v := viper.New()

	// This never actually returns an error
	_ = v.BindPFlags(pflags)

	v.SetEnvPrefix("loggo")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Sadly, this is needed, because AutomaticEnv doesn't load env variables during Unmarshal
	var c Immutable
	c.Storage.Path = v.GetString("storage.path")
	c.GRPC.Addr = v.GetString("grpc.addr")
	c.Web.Addr = v.GetString("web.addr")

	return c
}
