package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createFlags(t *testing.T, flags map[string]string, args []string) *pflag.FlagSet {
	pflags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	for name, value := range flags {
		pflags.String(name, value, "usage")
	}

	require.NoError(t, pflags.Parse(args))
	return pflags
}

func Test_ReadImmutable(t *testing.T) {
	// Not parallel because we need to use Setenv

	flags := map[string]string{
		"storage.path": "data",
		"grpc.addr":    ":20081",
		"web.addr":     ":20080",
	}

	tests := []struct {
		name              string
		args              []string
		env               map[string]string
		expectedImmutable Immutable
		wantErr           bool
	}{
		{
			name: "default flags",
			expectedImmutable: Immutable{
				Storage: StorageConfig{Path: "data"},
				GRPC:    GRPCConfig{Addr: ":20081"},
				Web:     WebConfig{Addr: ":20080"},
			},
		},
		{
			name: "env override",
			env: map[string]string{
				"LOGGO_STORAGE_PATH": "other-data",
				"LOGGO_GRPC_ADDR":    ":1337",
				"LOGGO_WEB_ADDR":     ":1338",
			},
			expectedImmutable: Immutable{
				Storage: StorageConfig{Path: "other-data"},
				GRPC:    GRPCConfig{Addr: ":1337"},
				Web:     WebConfig{Addr: ":1338"},
			},
		},
		{
			name: "arg override",
			args: []string{"--storage.path=arg-data", "--grpc.addr=:81", "--web.addr=:80"},
			expectedImmutable: Immutable{
				Storage: StorageConfig{Path: "arg-data"},
				GRPC:    GRPCConfig{Addr: ":81"},
				Web:     WebConfig{Addr: ":80"},
			},
		},
		{
			name: "priority check",
			args: []string{"--storage.path=arg-data", "--web.addr=:80"},
			env: map[string]string{
				"LOGGO_STORAGE_PATH": "env-data",
				"LOGGO_GRPC_ADDR":    ":1337",
			},
			expectedImmutable: Immutable{
				Storage: StorageConfig{Path: "arg-data"},
				GRPC:    GRPCConfig{Addr: ":1337"},
				Web:     WebConfig{Addr: ":80"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pflags := createFlags(t, flags, tt.args)

			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			gotImmutable := ReadImmutable(pflags)

			assert.Equal(t, tt.expectedImmutable, gotImmutable)
		})
	}
}
