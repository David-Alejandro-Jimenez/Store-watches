package services

import (
	"fmt"
	"unicode"
)

const (
	minUserNameLength = 5
	minPasswordLength = 10
)

//The ValidateUserName function is responsible for validating whether the way the user enters their username is correct
func ValidateUserName(userName string) error {
	if userName == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}

	if len(userName) < minUserNameLength {
		return fmt.Errorf("you cannot enter a name that is less than 5 characters")
	}

	return nil
}

//The ValidatePassword function es la responsable de verificar que la contraseña que el usuario introduce es lo más seguro posible
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