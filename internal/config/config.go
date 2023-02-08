package config

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/google/renameio"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// storageConfig contains options related to the log storage used by loggo.
type storageConfig struct {
	Directory string `mapstructure:"directory" yaml:"directory,omitempty"`
}

// grpcConfig contains details regarding the gRPC listener. Currently, it is used only by pigeons.
type grpcConfig struct {
	// Address of the gRPC listener
	Addr string `mapstructure:"addr" yaml:"addr,omitempty"`
}

// webConfig contains details regarding the web UI and service endpoints.
type webConfig struct {
	// Address of the HTTP listener
	Addr string `mapstructure:"addr" yaml:"addr,omitempty"`
}

// authConfig stores all of the auth credentials used during runtime.
// All options in this config support hot reloading and modifications.
type authConfig struct {
	// Credentials used by frontend-users
	Users []struct {
		Username     string `mapstructure:"username" yaml:"username,omitempty"`
		PasswordHash string `mapstructure:"password_hash" yaml:"username,omitempty"`
	} `mapstructure:"users" yaml:"users,omitempty"`

	// Username -> PasswordHash
	userMap map[string]string

	// Credentials used for authenticating and authorizing pigeon log deliveries
	Pigeons []struct {
		Name  string `mapstructure:"name" yaml:"name,omitempty"`
		Token string `mapstructure:"token" yaml:"token,omitempty"`
	} `mapstructure:"pigeons" yaml:"pigeons,omitempty"`

	// Token -> Name
	pigeonMap map[string]string

	// Token used for service-level authorization. Currently used for /metrics only
	ServiceToken string `mapstructure:"service_token" yaml:"service_token,omitempty"`
}

// loggoConfig contains the whole configuration
type loggoConfig struct {
	Storage storageConfig `mapstructure:"storage" yaml:"storage,omitempty"`
	GRPC    grpcConfig    `mapstructure:"grpc" yaml:"grpc,omitempty"`
	Web     webConfig     `mapstructure:"web" yaml:"web,omitempty"`
	Auth    authConfig    `mapstructure:"auth" yaml:"auth,omitempty"`
}

// Config stores the whole configuration used by Loggo, additionally wrapping it with modifications and hot-reloading.
type Config struct {
	mu sync.RWMutex
	v  *viper.Viper
	c  loggoConfig
}

// Read reads the configuration file, using the binded pflags for additional info.
func Read(pflags *pflag.FlagSet) (*Config, error) {
	v := viper.New()

	if pflags != nil {
		if err := v.BindPFlags(pflags); err != nil {
			return nil, fmt.Errorf("binding command line flags: %w", err)
		}
	}

	// Storage directory and listener addresses are additionally pulled from the env
	v.SetEnvPrefix("loggo")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.MustBindEnv("storage.directory")
	v.MustBindEnv("grpc.addr")
	v.MustBindEnv("web.addr")

	v.SetConfigType("yaml")
	v.SetConfigFile(v.GetString("config"))

	config := Config{v: v}
	if err := config.ReloadFromFile(); err != nil {
		return nil, err
	}

	return nil, nil
}

// Reload reloads the config using the specified YAML data and saves it to the file.
func (c *Config) Reload(config []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.v.ReadConfig(bytes.NewReader(config)); err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	if err := c.unmarshal(); err != nil {
		return err
	}

	if err := c.writeToFile(); err != nil {
		return fmt.Errorf("saving config to file: %w", err)
	}
	return nil
}

// ReloadFromFile reloads then config from the file.
func (c *Config) ReloadFromFile() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.v.ReadInConfig(); err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	return c.unmarshal()
}

// Marshal marshals the current config into its YAML representation.
func (c *Config) Marshal() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.marshal()
}

// StorageDirectory returns the storage.directory setting.
func (c *Config) StorageDirectory() string {
	return readConfigLocked(c, func(cfg *loggoConfig) string {
		return cfg.Storage.Directory
	})
}

// GRPCAddr returns the grpc.addr setting.
func (c *Config) GRPCAddr() string {
	return readConfigLocked(c, func(cfg *loggoConfig) string {
		return cfg.GRPC.Addr
	})
}

// WebAddr returns the web.addr setting.
func (c *Config) WebAddr() string {
	return readConfigLocked(c, func(cfg *loggoConfig) string {
		return cfg.Web.Addr
	})
}

// AuthServiceToken returns the auth.service_token setting.
func (c *Config) AuthServiceToken() string {
	return readConfigLocked(c, func(cfg *loggoConfig) string {
		return cfg.Auth.ServiceToken
	})
}

// AuthUsers returns the auth.users setting represented as a map from username to password hash.
func (c *Config) AuthUsers() map[string]string {
	return readConfigLocked(c, func(cfg *loggoConfig) map[string]string {
		return cfg.Auth.userMap
	})
}

// AuthPigeons returns the auth.pigeons setting represented as a map from token to pigeon name.
func (c *Config) AuthPigeons() map[string]string {
	return readConfigLocked(c, func(cfg *loggoConfig) map[string]string {
		return cfg.Auth.pigeonMap
	})
}

func (c *Config) unmarshal() error {
	var newCfg loggoConfig
	if err := c.v.UnmarshalExact(&newCfg); err != nil {
		return fmt.Errorf("unmarshaling config file: %w", err)
	}

	// Rebuild internal structures after unmarshaling
	userMap := make(map[string]string)
	for _, user := range newCfg.Auth.Users {
		userMap[user.Username] = user.PasswordHash
	}
	newCfg.Auth.userMap = userMap

	pigeonMap := make(map[string]string)
	for _, pigeon := range newCfg.Auth.Pigeons {
		pigeonMap[pigeon.Token] = pigeon.Name
	}
	newCfg.Auth.pigeonMap = pigeonMap

	c.c = newCfg
	return nil
}

func (c *Config) marshal() ([]byte, error) {
	data, err := yaml.Marshal(c.c)
	if err != nil {
		return nil, fmt.Errorf("marshaling config: %w", err)
	}
	return data, nil
}

func (c *Config) writeToFile() error {
	data, err := c.marshal()
	if err != nil {
		return err
	}

	// Write the config file atomically to avoid possible errors.
	file := c.v.ConfigFileUsed()
	if err := renameio.WriteFile(file, data, 0o644); err != nil {
		return fmt.Errorf("atomically writing %s: %w", file, err)
	}
	return nil
}

func readConfigLocked[T any](c *Config, f func(cfg *loggoConfig) T) T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return f(&c.c)
}
