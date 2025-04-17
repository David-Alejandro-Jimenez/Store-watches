// Package services provides implementations of input port interfaces for authentication services.
package services

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserLoginService implements the input.UserServiceLogin interface.

// It handles user authentication by validating input, checking user existence, verifying credentials, and issuing JWT tokens.
type UserLoginService struct {
	BaseAuthService
}

// NewUserLoginService constructs a UserLoginService with necessary dependencies.

// Parameters:
//   - userRepo: repository for user data access (output.UserRepository)
//   - userNameValidator: validator for username input (input.Validator)
//   - passwordValidator: validator for password input (input.Validator)

// Returns:
//   - input.UserServiceLogin: ready-to-use login service.
func NewUserLoginService(userRepo output.UserRepository, userNameValidator, passwordValidator input.Validator) input.UserServiceLogin {
	return &UserLoginService{
		BaseAuthService: BaseAuthService{
			UserRepo:          userRepo,
			UserNameValidator: userNameValidator,
			PasswordValidator: passwordValidator,
		},
	}
}

// Login authenticates a user and returns a signed JWT token.

// Steps:
//   1. Validate username format.
//   2. Check that the user exists in the repository.
//   3. Retrieve stored salt and password hash for the username.
//   4. Combine provided password with salt and compare hash.
//   5. Generate and return a JWT token if credentials are valid.

// Parameters:
//   - account: models.Account containing Username and Password.

// Returns:
//   - token string: a signed JWT token on success.
//   - error: non-nil if validation, lookup, or authentication fails.
func (l *UserLoginService) Login(account models.Account) (string, error) {
	// 1. Validate username
	if err := l.ValidateUserName(account.UserName); err != nil {
		return "", errors.NewValidationError(errors.ErrInvalidUsername)
	}

	// 2. Check user existence
	exists, err := l.CheckUserExists(account.UserName)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.NewNotFoundError(errors.ErrUserNotFound)
	}

	// 3. Retrieve salt and hashed password
	salt, err := l.UserRepo.GetSalt(account.UserName)
	if err != nil {
		return "", err
	}
	storedHash, err := l.UserRepo.GetHashPassword(account.UserName)
	if err != nil {
		return "", err
	}

	// 4. Verify password by hashing provided password with salt
	passwordWithSalt := append([]byte(account.Password), salt...)
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), passwordWithSalt)
	if err != nil {
		return "", errors.NewAuthError(errors.ErrInvalidCredentials)
	}

	// 5. Generate JWT token
	return l.GenerateToken(account.UserName)
}
