// Package securityAuth provides interfaces and implementations for password hashing
// and salt generation, supporting secure authentication workflows in the sale-watches application.
package securityAuth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Hasher defines methods for hashing password byte slices.
type Hasher interface {
	// Hash takes a password as a byte slice and returns its hashed representation or an error if hashing fails.
	Hash(password []byte) (string, error)
}

// Combined appends the given salt string to the password bytes and returns the combined slice, ready to be passed to a Hasher implementation.
func Combined(password, salt string) ([]byte) {
	combined := append([]byte(password), salt...)
	return combined
}

// BcryptHasher implements the Hasher interface using bcrypt with the DefaultCost.
type BcryptHasher struct {}

// Hash generates a bcrypt hash of the provided password bytes using DefaultCost.
// It returns the resulting hash as a string, or an InternalError if hashing fails.
func (h BcryptHasher) Hash(password []byte) (string, error) {
	var hashPassword, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalError("error hashing the password")
	}
	return string(hashPassword), nil
}

// Generator defines methods for generating random salts.
type Generator interface {
	// Generate returns a new base64-encoded salt string or an error if the underlying random source fails.
	Generate() (string, error)
}

// RandomSaltGenerator implements Generator by reading cryptographically secure bytes from crypto/rand and encoding them in base64.
type RandomSaltGenerator struct {}

// Generate creates a 32-byte random salt, encodes it using StdEncoding, and returns the resulting string. On failure, it returns an InternalError.
func (g RandomSaltGenerator) Generate() (string, error) {
	var salt = make([]byte, 32)
	var _, err = rand.Read(salt)
	if err != nil {	
		return "", errors.NewInternalError("error generating the salt")
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}