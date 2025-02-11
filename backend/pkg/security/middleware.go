package security

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/services"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "No autorizado", http.StatusForbidden)
			return
		}

		tokenString := cookie.Value
		err = services.ValidateToken(tokenString) 
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusForbidden)
			return
		}

	next.ServeHTTP(w, r)
	})
}