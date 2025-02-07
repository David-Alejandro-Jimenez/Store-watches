package public

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/models"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/repository"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/services"
)

var err error

func RegisterPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Disallowed method", http.StatusMethodNotAllowed)
		return
	}

	var application models.Account
	err = json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err = services.ValidateUserName(application.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	err = services.ValidatePassword(application.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := repository.GetUser(application.UserName) 
	if err != nil {
    	http.Error(w, "Server error while validating user", http.StatusInternalServerError)
   		return
	}

	if exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	err = repository.SaveUser(application.UserName, application.Password)
	if err != nil {
		http.Error(w, "Error saving user in database", http.StatusInternalServerError)
		return
	}

	token, err := services.GenerateJWT(application.UserName)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
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