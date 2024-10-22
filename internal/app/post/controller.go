package post

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type RepositoryPost interface {
	PostCreate(title, content string, user_id any) error
}

type ControllerPost struct {
	repo   RepositoryPost
	Logger *slog.Logger
}

func NewControllerPost(repo RepositoryPost, logger *slog.Logger) *ControllerPost {
	return &ControllerPost{
		repo,
		logger,
	}
}

func (m *ControllerPost) Create(w http.ResponseWriter, r *http.Request) {
	var post Post
	UserId := "2ec3fddd-e9d9-4511-a5c4-e01d2536e398"

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		m.Logger.Error("ошибка декодирования json", "error:", err)
		return
	}

	m.Logger.Info("обработка запроса в БД")

	err = m.repo.PostCreate(post.Title, post.Content, UserId)
	if err != nil {
		m.Logger.Error("ошикба при обработки запроса")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
