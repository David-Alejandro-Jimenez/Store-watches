// Package input defines service contracts for comment-related business logic, user operations, and input validation.
package input

// Validator defines a generic interface for input validation.
type Validator interface {
	// Validate enforces rules on the given input.
    // Returns a ValidationError if input is invalid.
	Validate(input interface{}) error
}