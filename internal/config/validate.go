package config

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// A validation error is returned during a reload of the config if the new configuration is invalid.
type ValidationError struct {
	part    string
	message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validating %s: %s", e.part, e.message)
}

var (
	errEmptyUserDetails    = &ValidationError{part: "auth config", message: "usernames and password hashes must not be empty"}
	errInvalidPasswordHash = &ValidationError{part: "auth config", message: "password hashes must be valid bcrypt hashes"}
	errEmptyPigeonDetails  = &ValidationError{part: "auth config", message: "pigeon names and tokens must not be empty"}
	errEmptyServiceToken   = &ValidationError{part: "auth config", message: "service token must not be empty"}
)

// Used for validating that a password hash is a valid bcrypt hash during config reloads.
var examplePassword = []byte("\x42")

// validateAuthConfig validates that all credentials are valid and usable by the system
func validateAuthConfig(c *AuthConfig) error {
	for _, user := range c.Users {
		if len(user.Username) < 1 || len(user.PasswordHash) < 1 {
			return errEmptyUserDetails
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), examplePassword)
		if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
			return errInvalidPasswordHash
		}
	}

	for _, pigeon := range c.Pigeons {
		if len(pigeon.Name) < 1 || len(pigeon.Token) < 1 {
			return errEmptyPigeonDetails
		}
	}

	if len(c.ServiceToken) < 1 {
		return errEmptyServiceToken
	}

	return nil
}

func validateMutableConfig(c *MutableConfig) error {
	return validateAuthConfig(&c.Auth)
}
