package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Repository interface {
	Search(UserID, Username string) ([]Search, error)
}

type UserController struct {
	repo   Repository
	logger *slog.Logger
}

func NewUserController(repo Repository, logger *slog.Logger) *UserController {
	return &UserController{
		repo,
		logger,
	}
}

func (c *UserController) GetSearch(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Missing 'username' query parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "localhost")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cookie")

	users, err := c.repo.Search("", username)
	if err != nil {
		c.logger.Error("Error fetching users", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
