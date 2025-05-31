// Package repository provides SQL-based implementations of output ports for data persistence.
// This file contains SqlCommentRepository, which implements CommentRepository using a MySQL database via sqlx.
package repository

import (
	"log"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/domain/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/output"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// SqlCommentRepository implements output.CommentRepository using a SQL database.
// It uses sqlx for database interactions and expects a valid *sqlx.DB connection.
//
// Fields:
//   - db: *sqlx.DB instance for executing queries.
type SqlCommentRepository struct {
	db *sqlx.DB
}

// NewSqlCommentRepository creates a new SqlCommentRepository.
// It fatally logs and exits if the provided db is nil, indicating a critical configuration error.

// Parameters:
//   - db: *sqlx.DB connection to the comments database.

// Returns:
//   - output.CommentRepository: initialized repository instance.
func NewSqlCommentRepository(db *sqlx.DB) output.CommentRepository {
	if db == nil {
		log.Fatal(errors.NewInternalError(errors.ErrDatabaseConnection).Error())
	}

	return &SqlCommentRepository{
		db: db,
	}
}

// GetComments retrieves all comments from the database, ordered by date descending.
// It performs a JOIN with the user_registration table to include the commenter's username.

// Returns:
//   - []models.Comment: slice of Comment models containing ID, Date, Content, UserID, UserName, and Rating.
//   - error: non-nil if the query fails, wrapped as an InternalError.
func(r *SqlCommentRepository) GetComments() ([]models.Comment, error) {
	var comment []models.Comment
	// Define SQL query to select comments and join with user table.
	const sqlQuery = `
	SELECT 
		c.ID,
		c.Date,
		c.Content,
		c.UserID,
		u.UserName AS UserName,
		c.Rating
	FROM comments c
	JOIN user_registration u
		ON c.UserID = u.UserID
		ORDER BY c.Date DESC
	`

	// Execute the query and scan results into comments slice.
	err := r.db.Select(&comment, sqlQuery)
	if err != nil {
		// Wrap low-level DB error in a domain-friendly InternalError.
		return nil, errors.NewInternalError(errors.ErrDatabaseQuery).WithError(err)
	}
	return comment, nil
} 

// SaveComment inserts a new comment into the database with the current timestamp.
// It uses parameterized queries to prevent SQL injection.

// Parameters:
//   - userID: ID of the authenticated user adding the comment.
//   - content: text content of the comment.
//   - rating: numerical rating associated with the comment.

// Returns:
//   - error: non-nil if the insert fails, wrapped as an InternalError.
func(r *SqlCommentRepository) SaveComment(userID int, content string, rating int) error {
	const query = `INSERT INTO comments (UserID, Content, Rating, Date)
	VALUES (?, ?, ?, NOW())`

	// Execute the insert query with provided parameters.
	_, err := r.db.Exec(query, userID, content, rating)
	if err != nil {
		// Return a generic InternalError on failure.
		return errors.NewInternalError("Error querying the database")
	}
	return nil
}