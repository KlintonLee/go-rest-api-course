package http

import (
	"fmt"
	"net/http"

	"github.com/KlintonLee/go-rest-api-course/internal/database/models"
	"github.com/KlintonLee/go-rest-api-course/internal/database/repositories"
	"github.com/gorilla/mux"
)

// Handler - stores the pointer to our comments service
type Handler struct {
	Router               *mux.Router
	CommentsRepositoryDb *repositories.CommentsRepositoryDB
}

// NewHandler - returns a pointer to a Handler
func NewHandler(commentsRepositoryDb *repositories.CommentsRepositoryDB) *Handler {
	return &Handler{
		CommentsRepositoryDb: commentsRepositoryDb,
	}
}

// SetupRoutes - sets up all routes for our application
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up routes")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comments", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comments", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comments/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comments/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comments/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I am alive!")
	})
}

// GetComment - retrieve a comment by id
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	comment, err := h.CommentsRepositoryDb.GetComment(id)
	if err != nil {
		fmt.Fprintf(w, "Error retrieving Comment by ID")
	}

	fmt.Fprintf(w, "%+v", comment)
}

func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.CommentsRepositoryDb.GetAllComments()
	if err != nil {
		fmt.Fprintf(w, "Error retrieving comments")
	}

	fmt.Fprintf(w, "%+v", comments)
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	comment, err := h.CommentsRepositoryDb.PostComment(&models.Comment{
		ID:   "154cf9a4-f226-4fc6-a7f0-c244811d6064",
		Slug: "/",
	})
	if err != nil {
		fmt.Fprintf(w, "Failed to post new comment")
	}

	fmt.Fprintf(w, "%+v", comment)
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	comment, err := h.CommentsRepositoryDb.UpdateComment("154cf9a4-f226-4fc6-a7f0-c244811d6064", &models.Comment{
		Slug: "/updated",
	})
	if err != nil {
		fmt.Fprintf(w, "Failed to update the comment")
	}

	fmt.Fprintf(w, "%+v", comment)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.CommentsRepositoryDb.DeleteComment(id); err != nil {
		fmt.Fprintf(w, "Failed to delete the comment")
	}

	fmt.Fprintf(w, "Successfully deleted comment")
}
