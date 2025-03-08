package public

import (
	"encoding/json"
	"net/http"

	
	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository"
)

// The Comments function is responsible for obtaining and sending the list of comments stored in the database in JSON format.
// 1. Obtaining: The comments are retrieved through a query in the database.
// 2. Error Handling: Any error is handled by returning a 500 error response.
// 3. Response Format: The response is configured to send JSON.
// 4. Encoding: Comments encoded in JSON are sent to the client.
// This feature is essential for displaying comments on the frontend of the application or for any client that needs to access this information.
func Comments(w http.ResponseWriter, r *http.Request) {
	comments, err := repository.GetComments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}