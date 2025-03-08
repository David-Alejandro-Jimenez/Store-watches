package auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
)

type UserServiceRegister interface {
	Register(account models.Account) (string, error)
}

type userServiceRegister struct {
	userRepo authRepository.UserRepository
}

func NewUsersService(userRepo authRepository.UserRepository) UserServiceRegister {
	return &userServiceRegister{userRepo: userRepo}
}

func (a *userServiceRegister) Register(account models.Account) (string, error) {
	if err := ValidateUserName(account.UserName); err != nil {
		return "", errors.NewBadRequestError("Username cannot be empty or must not have less than 5 characters")
	}

	if err := ValidatePassword(account.Password); err != nil {
		return "", errors.NewBadRequestError("Password cannot be empty, , must not have less than 10 characters, must have a number, a capital letter and a special character")
	}

	exists, err := a.userRepo.UserExists(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("error checking if the user exists")
	}
	if exists {
		return "", errors.NewConflictError("User already exists")
	}

	if err := a.userRepo.SaveUser(account.UserName, account.Password); err != nil {
		return "", errors.NewInternalError("error saving the user")
	}

	token, err := securityAuth.GenerateJWT(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("error generating the token")
	}

	return token, nil
}