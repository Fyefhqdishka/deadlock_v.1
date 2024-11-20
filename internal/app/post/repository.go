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

func (r *Repository) PostCreate(title, content string, user_id any) error {
	r.Logger.Debug(
		"PostCreate",
		"Обработка запроса на добавление поста в БД",
	)

	stmt := `INSERT INTO Posts (title, content, user_id) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(stmt, title, content, user_id)
	if err != nil {
		r.Logger.Error(
			"PostCreate",
			"Ошибка при добавлении нового поста",
		)

		return err
	}

	r.Logger.Debug(
		"PostCreate",
		"Добавление поста в БД прошло успешно",
	)

	return nil
}

func (m *Repository) PostGet() ([]Post, error) {
	stmt := `SELECT * FROM posts ORDER BY created_at DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		m.Logger.Error(
			"PostGet",
			"Не удалось обработать запрос БД",
		)

		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreateAt)
		if err != nil {
			m.Logger.Error(
				"PostGet",
				"Ошибка при сканировании строки",
				"error", err,
			)

			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		m.Logger.Error(
			"PostGet",
			"Ошибка при итерации по строкам",
			"error", err,
		)

		return nil, err
	}

	m.Logger.Debug(
		"PostGet",
		"Последние посты успешно получены",
	)

	return posts, nil
}
