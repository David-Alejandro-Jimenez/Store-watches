// Package http implements HTTP handlers for the sale-watches application.
// This file contains the LoginHandler, which is responsible for processing login requests.
package http

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/sale-watches/pkg/http"
)

// LoginHandler handles HTTP requests related to user login.

// It acts as an adapter between HTTP requests and the core domain's login functionality, using the UserServiceLogin interface to process login operations.
type LoginHandler struct {
	userServiceLogin input.UserServiceLogin
}

// NewLoginHandler creates a new instance of LoginHandler.

// It receives an implementation of the UserServiceLogin interface, which encapsulates the business logic for authenticating users.
func NewLoginHandler(userServiceLogin input.UserServiceLogin) *LoginHandler {
	return &LoginHandler{
		userServiceLogin: userServiceLogin,
	}
}

// Handle processes HTTP login requests.

// It validates that the request method is POST, decodes the JSON body into an Account model, and calls the login service to perform authentication. If the login operation is successful, it sets an authentication cookie based on the application's environment and sends a JSON response with a success message. Otherwise, it handles errors appropriately.
func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrMethodNotAllowed))
		return
	}

	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		httpUtil.HandleError(w, errors.NewBadRequestError(errors.ErrInvalidRequest))
		return
	}

	token, err := h.userServiceLogin.Login(account)
	if err != nil {
		httpUtil.HandleError(w, err)
		return
	}

	isProduction := os.Getenv("ENV") == "production"
	httpUtil.SetAuthCookie(w, token, isProduction)
	httpUtil.SendJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Successful login",
	})
}
