package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/KlintonLee/go-rest-api-course/internal/database/repositories"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - stores the pointer to our comments service
type Handler struct {
	Router               *mux.Router
	CommentsRepositoryDb *repositories.CommentsRepositoryDB
}

type Response struct {
	Message string
	Error   string
}

// NewHandler - returns a pointer to a Handler
func NewHandler(commentsRepositoryDb *repositories.CommentsRepositoryDB) *Handler {
	return &Handler{
		CommentsRepositoryDb: commentsRepositoryDb,
	}
}

// LoggingMiddleware - adds middleware around endpoints
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL,
			}).Info("Handled Request")

		next.ServeHTTP(w, r)
	})
}

// BasicAuth - a handy middleware function that will provide basic auth around
// specific endpoints
func BasicAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Basic auth endpoint hit")
		user, pass, ok := r.BasicAuth()
		if user == "admin" && pass == "admin" && ok {
			original(w, r)
		} else {
			sendErrorResponse(w, "Not authorized", errors.New("not authorized"))
			return
		}
	}
}

func validateJWTToken(accessToken string) bool {
	var mySigningKey = []byte("missionimpossible")
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there has been an error")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

// JWTAuth - a decorator function for jwt validation for endpoints
func JWTAuth(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("JWT authentication hit")
		authHeader := r.Header["Authorization"]
		if authHeader == nil {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		// Bearer jwt-token
		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			sendErrorResponse(w, "not authorized", errors.New("not authorized"))
			return
		}

		if validateJWTToken(authHeaderParts[1]) {
			original(w, r)
		} else {
			sendErrorResponse(w, "Not authorized", errors.New("not authorized"))
			return
		}
	}
}

// SetupRoutes - sets up all routes for our application
func (h *Handler) SetupRoutes() {
	log.Info("Setting up routes")
	h.Router = mux.NewRouter()
	h.Router.Use(LoggingMiddleware)

	h.Router.HandleFunc("/api/comments", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comments", JWTAuth(h.PostComment)).Methods("POST")
	h.Router.HandleFunc("/api/comments/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comments/{id}", BasicAuth(h.UpdateComment)).Methods("PUT")
	h.Router.HandleFunc("/api/comments/{id}", BasicAuth(h.DeleteComment)).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive"}); err != nil {
			panic(err)
		}
	})
}

func sendOkResponse(w http.ResponseWriter, resp interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(resp)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{Message: message, Error: err.Error()}); err != nil {
		log.Error(err)
	}
}
