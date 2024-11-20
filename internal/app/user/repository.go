package user

import (
	"database/sql"
	"log/slog"
)

type UserRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewUserRepository(db *sql.DB, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db,
		logger,
	}
}

func (r *UserRepository) Search(UserID, Username string) ([]Search, error) {
	stmt := `SELECT id, username FROM users WHERE username ILIKE $1`
	rows, err := r.DB.Query(stmt, Username)
	if err != nil {
		r.Logger.Error("Ошибки во время запроса в БД")
		return nil, err
	}
	defer rows.Close()

	var users []Search
	for rows.Next() {
		var search Search
		if err = rows.Scan(&search.Username); err != nil {
			r.Logger.Error("Ошибка")
			return nil, err
		}
		users = append(users, search)
	}

	return users, nil
}
