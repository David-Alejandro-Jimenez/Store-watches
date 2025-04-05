// Package middleware provides HTTP middleware components for request processing.
// It includes authentication, logging, and other cross-cutting concerns.
package middleware

import (
	"net/http"

	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
)

// AuthOptions contains configuration options for the authentication middleware.
// It allows for customization of authentication behavior, particularly which
// paths should be excluded from authentication requirements.
type AuthOptions struct {
	// ExcludedPaths is a slice of URL patterns that will not require authentication.
	// Both exact matches and directory prefixes (ending with '/') are supported.
	ExcludedPaths []string
}

// DefaultAuthOptions creates and returns a new AuthOptions instance with default values.
// The default configuration excludes common public paths like home, authentication pages,
// and static asset directories from requiring authentication.
// Returns a pointer to the newly created AuthOptions.
func DefaultAuthOptions() *AuthOptions {
	return &AuthOptions{
		ExcludedPaths: []string{"/", "/login", "/register", "/css/", "/js/", "/assets/"},
	}
}

// AuthMiddleware creates a middleware that verifies user authentication through JWT tokens.
// It intercepts HTTP requests, checks if the path requires authentication, and validates
// the authentication token stored in cookies.
//
// The middleware follows these steps:
//  1. Checks if the requested path is in the excluded paths list.
//  2. For paths requiring authentication, verifies the presence of a "token" cookie.
//  3. Validates the token using the security package.
//  4. If authentication succeeds, passes control to the next handler.
//  5. If authentication fails, returns a 401 Unauthorized response.
//
// Parameters:
//   - options: Configuration options for authentication, including paths to exclude
//
// Returns a Middleware function that can be used in an HTTP handler chain.
func AuthMiddleware(options *AuthOptions) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the path is excluded from authentication
			for _, path := range options.ExcludedPaths {
				if r.URL.Path == path || (path[len(path)-1] == '/' && len(r.URL.Path) >= len(path) && r.URL.Path[:len(path)] == path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Verify authentication by checking for token cookie
			cookie, err := r.Cookie("token")
			if err != nil || cookie.Value == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate the token using security package
			tokenString := cookie.Value
			err = securityAuth.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Authentication successful, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
