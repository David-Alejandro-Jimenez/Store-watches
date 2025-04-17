// Package input defines interfaces for comment management operations.
// It provides service contracts for implementations handling comment-related business logic.
package input

import "github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"

// CommentService defines the interface for comment management operations.
// Implementations should handle business logic for comment retrieval and creation.
type CommentService interface {
	// GetComments retrieves all available comments from the system.
	// Returns a slice of Comment models or an error if the operation fails.
	GetComments() ([]models.Comment, error)

	// AddComment creates a new comment entry in the system.
	// Accepts a Comment model as parameter and returns an error if the operation fails.
	// Potential errors may include invalid comment data or storage system failures.
	AddComment(comment models.Comment) error
}
