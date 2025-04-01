package securityAuth

import (
	"net/http"
)

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
