package securityAuth

import (
	"fmt"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var jwtSecret = []byte(viper.GetString("JWT_SECRET_KEY"))

// The GenerateJWT function is responsible for creating a JWT token (JSON Web Token) using the provided username and assigning it an expiration time.
// 1. Entry: Receives a username.
// 2. Claims: Configure the token's claims, including the username and an expiration time of 1 hour.
// 3. Creation and Signing: Create and sign the JWT token using the HS256 method and a secret key.
// 4. Output: Returns the signed JWT token or an error if any problem occurs.
// This implementation allows you to authenticate users and validate their identity in subsequent requests using the generated token.
func GenerateJWT(userName string) (string, error) {
	var claims = models.Claims{
		UserName: userName, 
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// The ValidateToken function is responsible for verifying that a JWT token is valid.
// 1. Parsing and Verification: The function uses jwt.Parse to parse the token and a callback function that validates the signing method and provides the secret key.
// 2. Error Handling: If there are problems during parsing or if the token is invalid, an error is returned.
// 3. Valid Token: If everything is correct, nil is returned.
// This feature is essential to ensure that the token presented in a request is legitimate and has not been tampered with, allowing the identity of the authenticated user to be trusted.
func ValidateToken(tokenString string) (error) {
	var token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

	return jwtSecret, nil
	})

	if err != nil {
		return err
	}
	
	if !token.Valid {
		return err
	}

	return nil
}