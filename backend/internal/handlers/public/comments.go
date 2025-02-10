package public

import (
	"encoding/json"
	"net/http"

	
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/repository"
)

func Comments(w http.ResponseWriter, r *http.Request) {
	comments, err := repository.GetComments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}