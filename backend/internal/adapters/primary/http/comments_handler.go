// Package http implements HTTP handlers for the sale-watches application.
// This file, which contains the CommentsHandler, is planned to be further enhanced in a future development phase. Note that only this file's functionality will be completed later.
package http

import (
	"net/http"

	"github.com/David-Alejandro-Jimenez/sale-watches/internal/core/ports/input"
	"github.com/David-Alejandro-Jimenez/sale-watches/pkg/errors"
	httpUtil "github.com/David-Alejandro-Jimenez/sale-watches/pkg/http"
)

// CommentsHandler handles HTTP requests related to comments.

// It acts as an adapter between HTTP requests and the business logic provided by the CommentService interface defined in the core domain. This handler currently supports retrieving comments.
type CommentsHandler struct {
	commentService input.CommentService
}

// NewCommentsHandler creates and returns a new instance of CommentsHandler.

// It receives an implementation of the CommentService interface, which contains the business logic for managing comments.
func NewCommentsHandler(commentService input.CommentService) *CommentsHandler {
	return &CommentsHandler{
		commentService: commentService,
	}
}

// Handle processes incoming HTTP requests to retrieve comments.

// It calls the GetComments method of the commentService to fetch comments. If an error occurs during the retrieval, it sends an HTTP error response with a 500 (Internal Server Error) status using a utility function. If successful, it returns the comments in JSON format with an HTTP 200 (OK) status.
func (h *CommentsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	comments, err := h.commentService.GetComments()
	if err != nil {
		httpUtil.HandleError(w, errors.NewInternalError("Error getting feedback"))
		return
	}

	httpUtil.SendJSONResponse(w, http.StatusOK, comments)
}
