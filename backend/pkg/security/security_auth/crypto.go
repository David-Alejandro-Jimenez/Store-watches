package securityAuth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Hash(password []byte) (string, error)
}

func Combined(password, salt string) ([]byte) {
	combined := append([]byte(password), salt...)
	return combined
}

type BcryptHasher struct {}

func (h BcryptHasher) Hash(password []byte) (string, error) {
	var hashPassword, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInternalError("error hashing the password")
	}
	return string(hashPassword), nil
}

type Generator interface {
	Generate() (string, error)
}

type RandomSaltGenerator struct {}

func (g RandomSaltGenerator) Generate() (string, error) {
	var salt = make([]byte, 32)
	var _, err = rand.Read(salt)
	if err != nil {	
		return "", errors.NewInternalError("error generating the salt")
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}