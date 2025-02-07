package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, salt string) (string, error) {
	var saltePassword = append([]byte(password), salt...)

	var hashPassword, err = bcrypt.GenerateFromPassword(saltePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func GenerateSalt() (string, error) {
	var salt = make([]byte, 16)
	var _, err = rand.Read(salt)
	if err != nil {	
		return "", err
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}