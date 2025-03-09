package auth

import (
	"fmt"
	"unicode"
)

const (
	minUserNameLength = 5
	minPasswordLength = 10
)

type Validator interface {
	Validate(input string) error
}

type UserNameValidator struct{}

func (c *UserNameValidator) Validate(username string) error {
	if username == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}

	if len(username) < minUserNameLength {
		return fmt.Errorf("you cannot enter a name that is less than 5 characters")
	}

	return nil
}

type PasswordValidator struct{}

func (p *PasswordValidator) Validate(password string) error {
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