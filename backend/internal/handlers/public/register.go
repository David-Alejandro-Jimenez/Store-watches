package public

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	authConfig "github.com/David-Alejandro-Jimenez/sale-watches/internal/config/auth_config"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/sale-watches/pkg/http"
)

func PostRegister(w http.ResponseWriter, r *http.Request) {
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

	isProduction := os.Getenv("ENV") == "production"

	httpUtil.SetAuthCookie(w, token, isProduction)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User created successfully",
		"redirect": "/",
	})
}
