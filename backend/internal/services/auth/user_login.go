package auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth_repository"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceLogin interface {
	Login(account models.Account) (string, error)
}

type userServiceLogin struct{
	userRepo authRepository.UserRepository
	userNameValidator Validator
	passwordValidator Validator
}

func NewUsersServiceLogin(userRepo authRepository.UserRepository,  userNameValidator, passwordValidator 	Validator) UserServiceLogin {
	return &userServiceLogin{
		userRepo: userRepo,
		userNameValidator: userNameValidator,
		passwordValidator: passwordValidator,
	}
}

func (l *userServiceLogin) Login(account models.Account) (string, error) {
	if err := l.userNameValidator.Validate(account.UserName); err != nil {
		return "", errors.NewBadRequestError("Username cannot be empty")
	}

	exists, err := l.userRepo.UserExists(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error checking if the user exists")
	}
	if !exists {
		return "", errors.NewConflictError("User does not exist")
	}

	salt, err := l.userRepo.GetSalt(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error getting user salt")
	}

	storedHash, err := l.userRepo.GetHashPassword(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error getting user password hash")
	}

	passwordWithSalt := append([]byte(account.Password), salt...)

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), passwordWithSalt)
	if err != nil {
		return "", errors.NewAuthError("Invalid username or password")
	}

	token, err := securityAuth.GenerateJWT(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error generating JWT token")
	}

	return token, nil
}