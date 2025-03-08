package public

import (
	"encoding/json"
	"net/http"
	"time"

	authConfig "github.com/David-Alejandro-Jimenez/sale-watches/internal/config/auth_config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

// LoginPOST handles HTTP POST requests for login
func LoginPOST(w http.ResponseWriter, r *http.Request) {
	var err error
	// Check that the HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Disallowed method", http.StatusMethodNotAllowed)
		return
	}

	var application models.Account // Structure to receive user data
	// The request body is decoded into the application structure
	err = json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Attempt to log in and get a token using the authentication service.
	// The Login method is expected to return a JWT token if the credentials are valid.
	token, err := authConfig.UserServiceLogin.Login(application)
	if err != nil {
		// Check if the error is of type *AppError, which indicates a handled application-level error
		if appErr, ok := err.(*errors.AppError); ok {
			// Return an HTTP error with the specific message and status code from the application error.
			http.Error(w, appErr.Message, appErr.Code)
		} else {
			// If the error is not an *AppError, treat it as an internal server error.
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		// Exit the function to prevent further processing since login failed.
		return
	}

	// A cookie is created with the session token
	cookie := http.Cookie{
		Name: "token",
		Value: token,
		Expires: time.Now().Add(12 * time.Hour), // Cookie expires in 12 hours
		HttpOnly: false, // Cookie is accessible from JavaScript (set to security)
		Path: "/",
		Secure: false, // Not marked as Secure (adjust according to environment)
		SameSite: http.SameSiteLaxMode, // Avoid CSRF issues in normal navigation
	}

	// The cookie is sent to the client
	http.SetCookie(w, &cookie) 

	// Responds with a success message and a redirect
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successful login",
		"redirect": "/",
	})
}