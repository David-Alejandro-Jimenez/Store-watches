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