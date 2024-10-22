package post

import (
	"database/sql"
	"log/slog"
)

type Repository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{db, logger}
}

func (m *Repository) PostCreate(title, content string, user_id any) error {
	m.Logger.Debug("Обработка запроса на добавление поста в БД")

	stmt := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(stmt, title, content, user_id)
	if err != nil {
		m.Logger.Error("Ошибка при добавлении нового поста")
		return err
	}

	m.Logger.Debug("Добавление поста в БД прошло успешно")
	return nil
}

func (m *Repository) PostGet(title, content string, user_id any) error {
	stmt := `SELECT FROM Posts WHERE title = $1 `
}
