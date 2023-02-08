package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// assertConfig asserts that the config is equal to what was expected using the exported methods
func assertConfig(t *testing.T, expect *loggoConfig, got *Config) {
	t.Helper()

	assert.Equal(t, expect.Storage.Directory, got.StorageDirectory())
	assert.Equal(t, expect.GRPC.Addr, got.GRPCAddr())
	assert.Equal(t, expect.Web.Addr, got.WebAddr())
	assert.Equal(t, expect.Auth.ServiceToken, got.AuthServiceToken())
	assert.Equal(t, expect.Auth.userMap, got.AuthUsers())
	assert.Equal(t, expect.Auth.pigeonMap, got.AuthPigeons())
}

func testFlags(t *testing.T, config string) *pflag.FlagSet {
	pflags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	_ = pflags.String("config", config, "")
	assert.NoError(t, pflags.Parse(nil))
	return pflags
}

func Test_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		pflags       *pflag.FlagSet
		expectConfig loggoConfig
		assertion    assert.ErrorAssertionFunc
	}{
		{
			name:   "full config in file",
			pflags: testFlags(t, ""),
		},
	}
}
