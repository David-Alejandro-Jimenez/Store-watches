package repository

import (
	"fmt"
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/models"
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/database"
)

// The GetComments function is responsible for retrieving all the comments stored in the database, combining information from the comments table and the users table to include the name of the user who made the comment.
// 1. Query: A JOIN is performed between the Comments and User_Registration tables to obtain both the comment details and the user name.
// 2. Execution and Error Handling: The query is executed and possible errors are checked.
// 3. Iteration: The results are traversed, scanned and stored in a slice.
// 4. Result: The list of comments obtained is returned, ready to be used in the application.
// This feature is essential for displaying or processing comments in the application, combining data from multiple tables to get a complete view of the information.
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