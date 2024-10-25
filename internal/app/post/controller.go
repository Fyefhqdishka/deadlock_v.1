package post

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type RepositoryPost interface {
	PostCreate(title, content string, user_id any) error
	PostGet() ([]Post, error)
}

type ControllerPost struct {
	repo   RepositoryPost
	logger *slog.Logger
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
		m.logger.Error("ошибка декодирования json", "error:", err)
		return
	}

	m.logger.Info("обработка запроса в БД")

	err = m.repo.PostCreate(post.Title, post.Content, UserId)
	if err != nil {
		m.logger.Error("ошикба при обработки запроса")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (m *ControllerPost) Get(w http.ResponseWriter, r *http.Request) {
	m.logger.Info("Обработка запроса на получение последних 10 постов")

	posts, err := m.repo.PostGet()
	if err != nil {
		http.Error(w, "Ошибка при получении постов", http.StatusInternalServerError)
		m.logger.Error("ошибка при обработке запроса к БД", "error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
