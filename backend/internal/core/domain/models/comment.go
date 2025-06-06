// Package models defines core domain entities for the sale‑watches application.

// This file declares the Comment type, representing user feedback with rating.
package models

// Comment represents a user’s feedback on a product or service.

// Fields:
//   - ID:        unique identifier of the comment.
//   - Date:      timestamp when the comment was posted, usually in ISO 8601 format.
//   - UserName:  identifier of the user who posted the comment.
//   - Content:   textual body of the comment.
//   - Rating:    numeric score given by the user (e.g., 1–5).
type Comment struct {
	ID int `db:"ID"`
	Date string `db:"Date"`
	UserID int `db:"UserID"`
	UserName string `db:"UserName"`
	Content string `db:"Content"`
	Rating int `db:"Rating"`
} 