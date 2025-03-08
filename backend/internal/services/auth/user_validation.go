package auth

import (
	"fmt"
	"unicode"
)

const (
	minUserNameLength = 5
	minPasswordLength = 10
)

// The ValidateUserName function is responsible for verifying that a username meets certain conditions before being accepted.
// 1. Empty fields: Rejected if the name is empty.
// 2. Minimum length: Rejected if the name is less than 5 characters.
// 3. Successful Validation: Returns nil when the username is valid.
// This feature is useful to ensure that the data entered by the user meets basic validity criteria before proceeding with other processes in the application.
func ValidateUserName(userName string) error {
	if userName == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}

	if len(userName) < minUserNameLength {
		return fmt.Errorf("you cannot enter a name that is less than 5 characters")
	}

	return nil
}

// The ValidatePassword function is responsible for validating that a password meets certain security requirements before being accepted.
// 1. Non-empty field: The password must not be empty.
// 2. Minimum length: A minimum length is required (at least 10 characters).
// 3. Complexity requirements:
		// At least one capital letter.
		// At least one numerical digit.
		// At least one special character (punctuation or symbol).
// If the password does not meet any of these criteria, the function returns an error with the corresponding message; otherwise, it returns nil, signaling that the password is valid.
func ValidatePassword(password string) error  {
	if password == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}
	if len(password) < minPasswordLength {
		return fmt.Errorf("you cannot enter a password that is less than 10 characters")
	}

	var hasUppercase bool
	var hasDigit bool
	var hasSpecialCharacter bool
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialCharacter = true
		}

		if hasUppercase && hasDigit && hasSpecialCharacter {
			break
		}
	}

	if !hasUppercase {
		return fmt.Errorf("the password must have at least one uppercase letter")
	}
	
	if !hasDigit {
		return fmt.Errorf("the password must have at least one number")
	}

	if !hasSpecialCharacter {
		return fmt.Errorf("the password must have some special character")
	}
	return nil
}