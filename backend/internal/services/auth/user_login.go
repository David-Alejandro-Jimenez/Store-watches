package auth

import (
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/auth"
	securityAuth "github.com/David-Alejandro-Jimenez/sale-watches/pkg/security/security_auth"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceLogin defines the interface that exposes the login operation.
// Its Login method receives an account and returns a JWT token on successful authentication.
type UserServiceLogin interface {
	// Login authenticates a user based on the data contained in models.Account.
	// Returns a JWT token and a possible error.
	Login(account models.Account) (string, error)
}

// userServiceLogin is the concrete implementation of UserServiceLogin.
// It is responsible for performing user authentication using a user repository.
type userServiceLogin struct{
	userRepo authRepository.UserRepository // Repository to access user data.
}

// NewUsersServiceLogin creates a new instance of UserServiceLogin by injecting the repository dependency into it.
// userRepo: Repository instance to be used to access user data.
func NewUsersServiceLogin(userRepo authRepository.UserRepository) UserServiceLogin {
	return &userServiceLogin{userRepo: userRepo}
}

// Login implements the user authentication logic.
// Performs the following steps:
// 1. Validates the username using ValidateUserName.
// 2. Verifies the existence of the user in the database.
// 3. Retrieves the stored salt and password hash.
// 4. Compares the provided password (concatenated with the salt) to the stored hash using bcrypt.
// 5. Generates and returns a JWT token using securityAuth.GenerateJWT.
// If any error occurs in any of these steps, return an error.
func (a *userServiceLogin) Login(account models.Account) (string, error) {
	// Validate the username.
	if err := ValidateUserName(account.UserName); err != nil {
		return "", errors.NewBadRequestError("Username cannot be empty")
	}

	// Verify that the user exists.
	exists, err := a.userRepo.UserExists(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error checking if the user exists")
	}
	if !exists {
		return "", errors.NewConflictError("User does not exist")
	}

	// Get the salt stored for the user.
	salt, err := a.userRepo.GetSalt(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error getting user salt")
	}

	// Get the hash of the password stored for the user.
	storedHash, err := a.userRepo.GetHashPassword(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error getting user password hash")
	}

	// Combine the provided password with the salt.
	passwordWithSalt := append([]byte(account.Password), salt...)

	// Compare the stored hash with the hash of the provided password.
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), passwordWithSalt)
	if err != nil {
		return "", errors.NewAuthError("Invalid username or password")
	}

	// Generate the JWT token for the authenticated user.
	token, err := securityAuth.GenerateJWT(account.UserName)
	if err != nil {
		return "", errors.NewInternalError("Error generating JWT token")
	}

	return token, nil
}