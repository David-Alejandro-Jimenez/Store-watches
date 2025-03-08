package securityAuth

import (
	"net/http"
)

//This AuthMiddleware function is an authentication middleware in Go. Its purpose is to intercept HTTP requests and check if the user has a valid cookie with an authentication token before allowing the request to continue. This middleware protects routes that require authentication. If the token is invalid or does not exist, the middleware responds with a 403 (Forbidden) error. If the token is valid, the request continues.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		tokenString := cookie.Value
		err = ValidateToken(tokenString) 
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusForbidden)
			return
		}

	next.ServeHTTP(w, r)
	})
}