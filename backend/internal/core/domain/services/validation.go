// Package services provides input‑port implementations for validating user credentials.

// It includes username and password validators enforcing basic length and character rules.
package services

import (
	"fmt"
	"unicode"
)

const (
	minUserNameLength = 5 // minimum number of characters for a valid username

	minPasswordLength = 10 // minimum number of characters for a valid password
)

// UserNameValidator checks that usernames are non‑empty and meet a minimum length.
//
// It returns an error if the username is empty or shorter than minUserNameLength.
type UserNameValidator struct{}

// Validate enforces the username rules:
//   • non‑empty
//   • at least minUserNameLength characters

// Returns a formatted error describing the violation.
func (c *UserNameValidator) Validate(username string) error {
	if username == "" {
		return fmt.Errorf("you cannot enter empty fields")
	}

	if len(username) < minUserNameLength {
		return fmt.Errorf("you cannot enter a name that is less than 5 characters")
	}

	return nil
}

// PasswordValidator checks that passwords are non‑empty, of sufficient length, and contain at least one uppercase letter, one digit, and one special character.

// It returns an error describing the first rule violation encountered.
type PasswordValidator struct{}

// Validate enforces the password rules:
//   • non‑empty
//   • at least minPasswordLength characters
//   • contains at least one uppercase letter (unicode.IsUpper) :contentReference[oaicite:1]{index=1}
//   • contains at least one digit (unicode.IsDigit)
//   • contains at least one punctuation or symbol (unicode.IsPunct or unicode.IsSymbol)

// Returns a formatted error for the first unmet requirement.
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
