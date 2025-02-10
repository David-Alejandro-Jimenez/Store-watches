package repository

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/models"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/repository/database"
)

func GetComments() ([]models.Comment, error) {
	query := `SELECT c.ID, u.UserName, c.Content, c.Date, c.Rating
				FROM Comments c
				JOIN User_Registration u ON c.UserID = u.UserID;`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf(err.Error(), http.StatusInternalServerError)
    }
		var comments = []models.Comment{}

		for rows.Next() {
            var comment models.Comment
            err := rows.Scan(&comment.ID, &comment.UserName, &comment.Content, &comment.Date, &comment.Rating)
            if err != nil {
                return nil, fmt.Errorf(err.Error(), http.StatusInternalServerError)
            }
            comments = append(comments, comment)
        }
		return comments, nil
}