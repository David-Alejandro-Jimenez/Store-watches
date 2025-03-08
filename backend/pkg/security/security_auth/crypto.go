package securityAuth

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

//This HashPassword function is responsible for hashing a password using the bcrypt algorithm and an additional salt. 
// 1. Concatenate the password with the salt to increase security.
// 2. Hash the password using bcrypt with a predefined cost.
// 3. Returns the resulting hash as a string.
// 4. If an error occurs, it returns it so it can be handled externally.
//This approach makes brute force and rainbow table attacks difficult, since bcrypt is resistant to these attacks and the use of a single salt improves the security of password storage.
func HashPassword(password string, salt string) (string, error) {
	var saltePassword = append([]byte(password), salt...)

	var hashPassword, err = bcrypt.GenerateFromPassword(saltePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

//The GenerateSalt function generates a random 32-byte salt, encodes it in Base64, and returns it as a text string.
// 1. Generate 32 bytes of random data with crypto/rand.
// 2. Handles possible errors in generation.
// 3. Encodes the result in Base64 for easy storage.
// 4. Returns the salt as a string along with nil if there are no errors.
// This salt is used in conjunction with a password before hashing it, preventing attacks such as rainbow tables and ensuring that identical passwords generate different hashes.
func GenerateSalt() (string, error) {
	var salt = make([]byte, 32)
	var _, err = rand.Read(salt)
	if err != nil {	
		return "", err
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}