package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	defaultServiceTokenLength = 16
	defaultConfigFilePerms    = 0o644
)

// AuthUser contains credentials used to authenticate a single user.
type AuthUser struct {
	Username     string
	PasswordHash string `yaml:"password_hash"`
}

// AuthPigeon contains credentials used to authenticate a single pigeon.
type AuthPigeon struct {
	Name  string
	Token string
}

// AuthConfig contains all settings related to authentication and authorization.
type AuthConfig struct {
	Users        []AuthUser
	Pigeons      []AuthPigeon
	ServiceToken string `yaml:"service_token"`
}

// MutableConfig contains all settings which can be changed during execution.
type MutableConfig struct {
	Auth AuthConfig
}

// Mutable wraps MutableConfig with support for hot-reloading and modifications.
type Mutable struct {
	path string

	mu sync.RWMutex
	c  MutableConfig

	// Username -> PasswordHash
	userMap map[string]string
	// Token -> Name
	pigeonMap map[string]string
}

// ReadMutable reads the mutable configuration file, creating one if it isn't present.
func ReadMutable(path string) (*Mutable, error) {
	m := &Mutable{path: path}

	if err := m.loadFromFile(true); err != nil {
		return nil, err
	}
	return m, nil
}

// Reload reloads the whole config from the file from which it was originally loaded.
func (m *Mutable) Reload() error {
	return m.loadFromFile(false)
}

// AuthServiceToken returns the auth.service_token setting.
func (m *Mutable) AuthServiceToken() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.c.Auth.ServiceToken
}

// AuthUsers returns the auth.users setting represented as a map from username to password hash.
func (m *Mutable) AuthUsers() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.userMap
}

// AuthPigeons returns the auth.pigeons setting represented as a map from token to pigeon name.
func (m *Mutable) AuthPigeons() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pigeonMap
}

// SetMutableConfig changes the used mutable config, if the one passed is valid.
func (m *Mutable) SetMutableConfig(c MutableConfig) error {
	if err := validateMutableConfig(&c); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if err := writeMutableConfig(m.path, &c); err != nil {
		return fmt.Errorf("saving modified config: %w", err)
	}

	m.c = c
	return nil
}

func (m *Mutable) loadFromFile(generate bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mcfg, err := readMutableConfig(m.path)
	if errors.Is(err, fs.ErrNotExist) && generate {
		// Try creating a default config file with a generated service token
		mcfg, err = generateMutableConfig(m.path)
		if err != nil {
			return fmt.Errorf("creating new config since one wasn't found: %w", err)
		}
	} else if err != nil {
		return err
	}

	if err := validateMutableConfig(&mcfg); err != nil {
		return err
	}

	// Actually replace the config only once we've made sure everything's ok
	m.c = mcfg
	m.rebuildInternal()
	return nil
}

// rebuildInternal rebuilds internal structures after unmarshaling.
func (m *Mutable) rebuildInternal() {
	userMap := make(map[string]string)
	for _, user := range m.c.Auth.Users {
		userMap[user.Username] = user.PasswordHash
	}
	m.userMap = userMap

	pigeonMap := make(map[string]string)
	for _, pigeon := range m.c.Auth.Pigeons {
		pigeonMap[pigeon.Token] = pigeon.Name
	}
	m.pigeonMap = pigeonMap
}

func readMutableConfig(path string) (MutableConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return MutableConfig{}, fmt.Errorf("reading config file contents: %w", err)
	}

	var newCfg MutableConfig
	if err := yaml.Unmarshal(data, &newCfg); err != nil {
		return MutableConfig{}, fmt.Errorf("unmarshaling config file as YAML: %w", err)
	}

	return newCfg, nil
}

func writeMutableConfig(path string, c *MutableConfig) error {
	data, err := yaml.Marshal(&c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, defaultConfigFilePerms); err != nil {
		return fmt.Errorf("writing config to file: %w", err)
	}

	return nil
}

func generateMutableConfig(path string) (MutableConfig, error) {
	serviceTokenBytes := make([]byte, defaultServiceTokenLength)
	if _, err := rand.Read(serviceTokenBytes); err != nil {
		return MutableConfig{}, fmt.Errorf("generating service token: %w", err)
	}

	c := MutableConfig{Auth: AuthConfig{ServiceToken: hex.EncodeToString(serviceTokenBytes)}}
	if err := writeMutableConfig(path, &c); err != nil {
		return MutableConfig{}, fmt.Errorf("saving generated config: %w", err)
	}

	return c, nil
}
