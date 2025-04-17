// Package errors defines common error messages used throughout the application.
// It provides consistent error messaging for various domains including authentication,
// database operations, validation, and API rate limiting.
package errors

// Common error messages
const (
	// Authentication errors
	ErrInvalidCredentials = "Invalid credentials"
	ErrUserNotFound       = "User not found"
	ErrUserAlreadyExists  = "The user already exists"
	ErrInvalidUsername    = "Invalid username"
	ErrInvalidPassword    = "Invalid password"
	ErrTokenGeneration    = "Error generating token"
	ErrTokenValidation    = "Invalid or expired token"

	// Database errors
	ErrDatabaseConnection = "Database connection error"
	ErrDatabaseQuery      = "Error executing query"
	ErrDatabaseInsert     = "Error inserting into the database"
	ErrDatabaseUpdate     = "Error updating the database"
	ErrDatabaseDelete     = "Error deleting from database"

	// Validation errors
	ErrEmptyField        = "The field cannot be empty"
	ErrInvalidFormat     = "Invalid format"
	ErrInvalidLength     = "Invalid length"
	ErrInvalidCharacters = "Characters not allowed"
	
	// Comment operations errors
	ErrCommentNotFound = "Comment not found"
	ErrCommentCreation = "Error creating comment"
	ErrCommentUpdate   = "Error updating comment"
	ErrCommentDelete   = "Error deleting comment"
	
	// Rate limiting errors
	ErrTooManyRequests   = "Too many requests"
	ErrRateLimitExceeded = "Rate limit exceeded"

	// General API errors
	ErrInternalServer   = "Internal Server Error"
	ErrMethodNotAllowed = "Disallowed method"
	ErrInvalidRequest   = "Invalid request"
	ErrUnauthorized     = "Unauthorized"
	ErrForbidden        = "Prohibited access"
)
