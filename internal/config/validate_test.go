package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateAuthConfig(t *testing.T) {
	t.Parallel()

	const (
		username = "username"
		// bcrypt hash of "password"
		passwordHash = "$2y$10$00iDz1fSEs.RrWniiJRP1ul.zyY8MsZM2rcBlXlkccOYfctpHNfMq"
		serviceToken = "xxxx-yyyy-zzzz"
		pigeonName   = "pigeon"
		pigeonToken  = "0000-1111-2222"
	)

	tests := []struct {
		name        string
		config      AuthConfig
		expectedErr error
	}{
		{
			name: "empty username",
			config: AuthConfig{
				Users:        []AuthUser{{Username: "", PasswordHash: passwordHash}},
				Pigeons:      []AuthPigeon{},
				ServiceToken: serviceToken,
			},
			expectedErr: errEmptyUserDetails,
		},
		{
			name: "empty password hash",
			config: AuthConfig{
				Users:        []AuthUser{{Username: username, PasswordHash: ""}},
				Pigeons:      []AuthPigeon{},
				ServiceToken: serviceToken,
			},
			expectedErr: errEmptyUserDetails,
		},
		{
			name: "invalid password hash",
			config: AuthConfig{
				Users:        []AuthUser{{Username: username, PasswordHash: "bad hash"}},
				Pigeons:      []AuthPigeon{},
				ServiceToken: serviceToken,
			},
			expectedErr: errInvalidPasswordHash,
		},
		{
			name: "empty pigeon name",
			config: AuthConfig{
				Users:        []AuthUser{{Username: username, PasswordHash: passwordHash}},
				Pigeons:      []AuthPigeon{{Name: "", Token: pigeonToken}},
				ServiceToken: serviceToken,
			},
			expectedErr: errEmptyPigeonDetails,
		},
		{
			name: "empty pigeon token",
			config: AuthConfig{
				Users:        []AuthUser{{Username: username, PasswordHash: passwordHash}},
				Pigeons:      []AuthPigeon{{Name: pigeonName, Token: ""}},
				ServiceToken: serviceToken,
			},
			expectedErr: errEmptyPigeonDetails,
		},
		{
			name: "empty service token",
			config: AuthConfig{
				Users:        []AuthUser{{Username: username, PasswordHash: passwordHash}},
				Pigeons:      []AuthPigeon{{Name: pigeonName, Token: pigeonToken}},
				ServiceToken: "",
			},
			expectedErr: errEmptyServiceToken,
		},
		{
			name: "valid config",
			config: AuthConfig{
				Users:        []AuthUser{{Username: username, PasswordHash: passwordHash}},
				Pigeons:      []AuthPigeon{{Name: pigeonName, Token: pigeonToken}},
				ServiceToken: serviceToken,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotErr := validateAuthConfig(&tt.config)

			assert.ErrorIs(t, tt.expectedErr, gotErr)
		})
	}
}
