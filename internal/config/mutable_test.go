package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testdataConfig(name string) string {
	return filepath.Join("testdata", name)
}

func assertMutable(t *testing.T, expect *Mutable, got *Mutable) {
	assert.Equal(t, expect.userMap, got.AuthUsers())
	assert.Equal(t, expect.pigeonMap, got.AuthPigeons())
	assert.Equal(t, expect.c.Auth.ServiceToken, got.AuthServiceToken())
}

func Test_ReadMutable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		path          string
		expectMutable *Mutable
		wantErr       bool
	}{
		{
			name: "valid config",
			path: testdataConfig("loggo-valid.yaml"),
			expectMutable: &Mutable{
				c: MutableConfig{Auth: AuthConfig{ServiceToken: "xxxx-yyyy-zzzz"}},
				userMap: map[string]string{
					"username": "$2y$10$00iDz1fSEs.RrWniiJRP1ul.zyY8MsZM2rcBlXlkccOYfctpHNfMq",
				},
				pigeonMap: map[string]string{
					"0000-1111-2222": "pigeon",
				},
			},
		},
		{
			name:    "bad YAML",
			path:    testdataConfig("loggo-bad.yaml"),
			wantErr: true,
		},
		{
			name:    "invalid config values",
			path:    testdataConfig("loggo-invalid.yaml"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotMutable, err := ReadMutable(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assertMutable(t, tt.expectMutable, gotMutable)
			}
		})
	}
}

func Test_Mutable_Modifications(t *testing.T) {
	t.Parallel()

	// A new config path should lead to a generated config file
	nonExistentConfigPath := filepath.Join(t.TempDir(), "loggo.yaml")

	m1, err := ReadMutable(nonExistentConfigPath)

	require.NoError(t, err)
	assert.Empty(t, m1.AuthUsers())
	assert.Empty(t, m1.AuthPigeons())
	assert.NotEmpty(t, m1.AuthServiceToken())

	// Loading another config from the same path should load the same settings
	m2, err := ReadMutable(nonExistentConfigPath)

	require.NoError(t, err)
	assertMutable(t, m1, m2)

	// Modifying the original config should also write it to the file,
	// so the new config can pick the changes up during Reload
	setErr := m1.SetMutableConfig(MutableConfig{Auth: AuthConfig{ServiceToken: "custom token"}})
	reloadErr := m2.Reload()

	assert.NoError(t, setErr)
	assert.NoError(t, reloadErr)
	assertMutable(t, m1, m2)

	// Changing to an invalid config shouldn't change anything
	setErr = m1.SetMutableConfig(MutableConfig{})

	assert.Error(t, setErr)
}

func Test_ReadMutable_InvalidPath(t *testing.T) {
	t.Parallel()

	m := &Mutable{path: ""}

	// Reading/generating a config for an invalid path should always fail
	_, readErr := ReadMutable("")
	setErr := m.SetMutableConfig(MutableConfig{Auth: AuthConfig{ServiceToken: "valid"}})

	assert.Error(t, readErr)
	assert.Error(t, setErr)
}
