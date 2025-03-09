package auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth_repository"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
)

type UserServiceRegister interface {
	Register(account models.Account) (string, error)
}

type userServiceRegister struct {
	userRepo authRepository.UserRepository
	userNameValidator Validator
	passwordValidator Validator
}

func NewUsersServiceRegister(userRepo authRepository.UserRepository, userNameValidator, passwordValidator 	Validator) UserServiceRegister {
	return &userServiceRegister{
		userRepo: userRepo,
		userNameValidator: userNameValidator,
		passwordValidator: passwordValidator,
	}
}

func (r *userServiceRegister) Register(account models.Account) (string, error) {
	if err := r.userNameValidator.Validate(account.UserName); err != nil {
		return "", errors.NewBadRequestError("Username cannot be empty or must not have less than 5 characters")
	}

	if err :=  r.passwordValidator.Validate(account.Password); err != nil {
		return "", errors.NewBadRequestError("Password cannot be empty, , must not have less than 10 characters, must have a number, a capital letter and a special character")
	}

	exists, err := r.userRepo.UserExists(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("error checking if the user exists")
	}
	if exists {
		return "", errors.NewConflictError("User already exists")
	}

	if err := r.userRepo.SaveUser(account.UserName, account.Password); err != nil {
		return "", errors.NewInternalError("error saving the user")
	}

	token, err := securityAuth.GenerateJWT(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("error generating the token")
	}

	return token, nil
}