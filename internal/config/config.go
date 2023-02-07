package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// storageConfig contains options related to the log storage used by loggo.
type storageConfig struct {
	Directory string `mapstructure:"directory"`
}

// grpcConfig contains details regarding the gRPC listener. Currently, it is used only by pigeons.
type grpcConfig struct {
	// Address of the gRPC listener
	Addr string `mapstructure:"addr"`
}

// webConfig contains details regarding the web UI and service endpoints.
type webConfig struct {
	// Address of the HTTP listener
	Addr string `mapstructure:"addr"`
}

// authConfig stores all of the auth credentials used during runtime.
// All options in this config support hot reloading and modifications.
type authConfig struct {
	// Credentials used by frontend-users
	Users []struct {
		Username     string `mapstructure:"username"`
		PasswordHash string `mapstructure:"password_hash"`
	} `mapstructure:"users"`

	// Username -> PasswordHash
	userMap map[string]string

	// Credentials used for authenticating and authorizing pigeon log deliveries
	Pigeons []struct {
		Name  string `mapstructure:"name"`
		Token string `mapstructure:"token"`
	} `mapstructure:"pigeons"`

	// Token -> Name
	pigeonMap map[string]string

	// Token used for service-level authorization. Currently used for /metrics only
	ServiceToken string `mapstructure:"service_token"`
}

// loggoConfig contains the whole configuration
type loggoConfig struct {
	Storage storageConfig `mapstructure:"storage"`
	GRPC    grpcConfig    `mapstructure:"grpc"`
	Web     webConfig     `mapstructure:"web"`
	Auth    authConfig    `mapstructure:"auth"`
}

// Config stores the whole configuration used by Loggo, additionally wrapping it with modifications and hot-reloading.
type Config struct {
	c loggoConfig
}

// Read reads the configuration file, using the binded pflags for additional info.
func Read(pflags *pflag.FlagSet) (*Config, error) {
	if pflags != nil {
		if err := viper.BindPFlags(pflags); err != nil {
			return nil, fmt.Errorf("binding command line flags: %w", err)
		}
	}

	// Storage directory and listener addresses are additionally pulled from the env
	viper.SetEnvPrefix("loggo")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(viper.GetString("config"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return nil, nil
}

func (c *Config) reload() error {
	var newCfg Config
	if err := viper.UnmarshalExact(&newCfg); err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	return nil
}
