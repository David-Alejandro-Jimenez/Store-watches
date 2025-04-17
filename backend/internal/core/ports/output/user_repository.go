// Package output defines interfaces for user persistence operations.
// It provides contracts for storage implementations handling user data management and security-sensitive operations.	
package output

// UserRepository defines the interface for user storage operations.
// Implementations should handle user data persistence, password security, and credential verification.
type UserRepository interface {
	// UserExists checks if a username is already registered in the system.
	// Returns true if the username exists in storage, false otherwise.
	// May return an error for database connectivity issues or storage failures.
	UserExists(username string) (bool, error)

	// GetHashPassword retrieves the hashed password for a registered user.
	// Primarily used during authentication processes. Returns an error if:
		// - User doesn't exist
		// - Password record is corrupted
		// - Storage system failure occurs
	GetHashPassword(username string) (string, error)

	// GetSalt retrieves the cryptographic salt used in password hashing.
	// The salt should be uniquely generated per user during registration.
	// Returns an error if the user doesn't exist or salt retrieval fails.
	GetSalt(username string) (string, error)

	// SaveUser persists a new user record with secure credential storage.
	// Implementations should:
		// - Generate unique salt per user
		// - Hash password using cryptographically secure methods
		// - Prevent duplicate usernames

	// Returns error for:
		// - Duplicate username
		// - Invalid credentials
		// - Storage persistence failures
	SaveUser(username, password string) error
}
