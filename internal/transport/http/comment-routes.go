package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KlintonLee/go-rest-api-course/internal/database/models"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// GetComment - retrieve a comment by id
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	comment, err := h.CommentsRepositoryDb.GetComment(id)
	if err != nil {
		sendErrorResponse(w, "Error retrieving Comment by ID", err)
		return
	}

	if sendOkResponse(w, comment, http.StatusOK); err != nil {
		panic(err)
	}
}

// GetAllComments - retrieve all created comments
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.CommentsRepositoryDb.GetAllComments()
	if err != nil {
		sendErrorResponse(w, "Error retrieving comments", err)
		return
	}

	if err := sendOkResponse(w, comments, http.StatusOK); err != nil {
		panic(err)
	}
}

// PostComment - create a new comment
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	comment.ID = uuid.NewV4().String()
	comment.CreatedAt = time.Now()

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body", err)
		return
	}

	_, err := h.CommentsRepositoryDb.PostComment(&comment)
	if err != nil {
		sendErrorResponse(w, "Failed to post new comment", err)
		return
	}

	if err := sendOkResponse(w, comment, http.StatusCreated); err != nil {
		panic(err)
	}
}

// UpdateComment - updates an existing comment by id
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode JSON Body", err)
		return
	}

	updatedComment, err := h.CommentsRepositoryDb.UpdateComment(id, &comment)
	if err != nil {
		sendErrorResponse(w, "Failed to update the comment", err)
		return
	}

	if err := sendOkResponse(w, updatedComment, http.StatusOK); err != nil {
		panic(err)
	}
}

// DeleteComment - delete an existing comment by id
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)

	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.CommentsRepositoryDb.DeleteComment(id); err != nil {
		sendErrorResponse(w, "Failed to delete the comment", err)
		return
	}
}
