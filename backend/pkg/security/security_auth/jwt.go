// Package securityAuth provides JWT creation and validation services for the sale‑watches application’s authentication flow.
package securityAuth

import (
	"fmt"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

// JWTService manages operations related to JSON Web Tokens.
// It uses a secret key to sign and verify tokens.
type JWTService struct {
	secretKey []byte
}

// NewJWTService creates a new JWTService with the given secret.
// secretKey: the HMAC secret used to sign and validate tokens.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

// GenerateJWT generates a signed JWT for the specified userName.
// The token embeds the username and an expiration set to one hour from now.
func (j *JWTService) GenerateJWT(userName string) (string, error) {
	var claims = models.Claims{
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken parses and verifies the given tokenString.
// It returns an error if parsing fails, the signature is invalid, or the token is expired.
func (j *JWTService) ValidateToken(tokenString string) error {
	var token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// defaultJWTService holds the globally configured JWTService for convenience.
var defaultJWTService *JWTService

// SetDefaultJWTService configures the package‑level JWTService.
// It should be called once at application startup with the secret key.
func SetDefaultJWTService(secretKey string) {
	defaultJWTService = NewJWTService(secretKey)
}

// GenerateJWT signs a token for userName using the default service.
// Returns an error if the service has not been initialized.
func GenerateJWT(userName string) (string, error) {
	if defaultJWTService == nil {
		return "", fmt.Errorf("JWT service not initialized")
	}
	return defaultJWTService.GenerateJWT(userName)
}

// ValidateToken verifies tokenString using the default service.
// Returns an error if validation fails or the service is not initialized.
func ValidateToken(tokenString string) error {
	if defaultJWTService == nil {
		return fmt.Errorf("servicio JWT no inicializado")
	}
	return defaultJWTService.ValidateToken(tokenString)
}
