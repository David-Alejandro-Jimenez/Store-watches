package public

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/config/auth_config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

// The RegisterPOST function handles the registration of a new user via an HTTP POST request. Its purpose is to receive the registration data, validate it, save the user to the database, generate a JWT token, and set a cookie with that token.
// Perform the following steps:
// 1. HTTP method verification: Only accept POST.
// 2. JSON Decoding: Extracts log data.
// 3. Validation: Check that the username and password meet the requirements.
// 4. User existence: Check if the user is already registered.
// 5. Saved: Inserts the new user into the database.
// 6. JWT Generation: Create a token for authentication.
// 7. Cookie Set: Send the token to the client using a cookie.
// 8. Successful response: Informs the client that the registration was successful.
// This process ensures that user registration is handled securely and consistently, applying appropriate validations and protections at every step.
func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != http.MethodPost {
		http.Error(w, "Disallowed method", http.StatusMethodNotAllowed)
		return
	}

	var application models.Account
	err = json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	token, err := authConfig.UserServiceRegister.Register(application)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			http.Error(w, appErr.Message, appErr.Code)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	
	cookie := http.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(12 * time.Hour),
		HttpOnly: false,
		Path: "/",
		Secure: false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie) 

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"redirect": "/",
	})
}