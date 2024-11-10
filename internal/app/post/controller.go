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

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		m.logger.Error(
			"Create",
			"ошибка декодирования json",
			"error:", err,
		)

		return
	}

	m.logger.Info(
		"Create",
		"обработка запроса в БД",
	)

	userID := r.Context().Value("user_id")
	if userID == nil {
		m.logger.Warn(
			"Create",
			"userID не найден в контексте",
		)

		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		m.logger.Warn(
			"Create",
			"userID не является строкой",
		)

		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	err = m.repo.PostCreate(post.Title, post.Content, userIDStr)
	if err != nil {
		m.logger.Error(
			"Create",
			"ошибка при обработки запроса",
		)

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (m *ControllerPost) Get(w http.ResponseWriter, r *http.Request) {
	m.logger.Info(
		"Get",
		"Обработка запроса на получение последних 10 постов",
	)

	posts, err := m.repo.PostGet()
	if err != nil {
		m.logger.Error(
			"Get",
			"ошибка при обработке запроса к БД",
			"error:", err,
		)
		
		http.Error(w, "Ошибка при получении постов", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}
