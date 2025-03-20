package public

import (
	"encoding/json"
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/repository/comments"
)

type HandlerComment struct {
	CommentRepo commentsRepository.CommentRepository
}

func NewHandlerComment(commentRepo commentsRepository.CommentRepository) *HandlerComment {
	return &HandlerComment{CommentRepo: commentRepo}
}

func (h *HandlerComment) Comments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.CommentRepo.Obtener()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")	
	json.NewEncoder(w).Encode(comments)
}