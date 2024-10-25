package post

import (
	"database/sql"
	"log/slog"
)

type Repository struct {
	DB     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{db, logger}
}

func (m *Repository) PostCreate(title, content string, user_id any) error {
	m.logger.Debug("Обработка запроса на добавление поста в БД")

	stmt := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(stmt, title, content, user_id)
	if err != nil {
		m.logger.Error("Ошибка при добавлении нового поста")
		return err
	}

	m.logger.Debug("Добавление поста в БД прошло успешно")
	return nil
}

func (m *Repository) PostGet() ([]Post, error) {
	stmt := `SELECT * FROM posts ORDER BY created_at DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		m.logger.Error("Не удалось обработать запрос БД")
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreateAt)
		if err != nil {
			m.logger.Error("Ошибка при сканировании строки", "error", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		m.logger.Error("Ошибка при итерации по строкам", "error", err)
		return nil, err
	}

	m.logger.Debug("Последние посты успешно получены")
	return posts, nil
}
