package public

import (
	"encoding/json"
	"net/http"

	//"os"
	"time"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/models"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/repository"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/services"
	"golang.org/x/crypto/bcrypt"
)

var err error

//func LoginGET(w http.ResponseWriter, r *http.Request) {
	//filePath := "./../frontend/pages/login.html"
	//if _, err := os.Stat(filePath); err != nil {
	//	http.Error(w, "Archivo no encontrado", http.StatusNotFound)
	//	return
	//}
	//http.ServeFile(w, r, filePath)
//}

func LoginPOST(w http.ResponseWriter, r *http.Request) {
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

	exists, err := repository.GetUser(application.UserName)
	if err != nil {
		http.Error(w, "Server error while validating user", http.StatusInternalServerError)
		return
		} 
		
	if !exists {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	salt, err := repository.GetSalt(application.UserName)
	if err != nil {
		http.Error(w, "Server error retrieving salt", http.StatusInternalServerError)
		return
	}
	
	storeHash, err :=  repository.GetHashPassword(application.UserName)
	if err != nil {
		http.Error(w, "Server error retrieving hash", http.StatusInternalServerError)
		return
	}

	var passwordWithSalt = append([]byte(application.Password), salt...)
	err = bcrypt.CompareHashAndPassword([]byte(storeHash), passwordWithSalt)
	if err != nil {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successful login",
		"redirect": "/",
	})
}