package commentsRepository

import (
	"database/sql"
	"errors"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
)

type CommentRepository interface {
	Obtener() ([]models.Comment, error)
}

type Comments struct {
	db *sql.DB
}

func NewComments(db *sql.DB) *Comments {
	return &Comments{db: db}
} 

func (c *Comments) Obtener() ([]models.Comment, error) {
	query := `SELECT c.ID, u.UserName, c.Content, c.Date, c.Rating
				FROM Comments c
				JOIN User_Registration u ON c.UserID = u.UserID;`
	
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, errors.New("error executing SQL query")
    }
	defer rows.Close()

		var comments []models.Comment

		for rows.Next() {
            var comment models.Comment
            err := rows.Scan(&comment.ID, &comment.UserName, &comment.Content, &comment.Date, &comment.Rating)
            if err != nil {
                return nil, errors.New("unexpected error")
            }
            comments = append(comments, comment)
        }

		if err = rows.Err(); err != nil {
			return nil, errors.New("unexpected error")
		}
		return comments, nil
}