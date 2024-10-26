package auth

import (
	"database/sql"
	"fmt"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/user"
	"log/slog"
)

type Repository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{
		DB:     db,
		Logger: logger,
	}
}

func (m *Repository) RegistrationUser(user user.User) error {
	m.Logger.Debug("starting registration user")

	stmt, err := m.DB.Prepare("INSERT INTO users (name, username, email, password, gender, dob) VALUES ($1,$2,$3,$4,$5,$6)")
	if err != nil {
		m.Logger.Error("error preparing statement", "err:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Username, user.Email, user.Password, user.Gender, user.Dob)
	if err != nil {
		m.Logger.Error("error executing statement", "err:", err)
		return err
	}

	m.Logger.Debug("registration user finished", "username", user.Username)
	return nil
}

func (m *Repository) LoginUser(username, password string) (string, error) {
	m.Logger.Debug("starting login user")

	var UserID string
	var passwordHash string

	stmt := `SELECT id, password FROM users WHERE username = $1`
	err := m.DB.QueryRow(stmt, username).Scan(&UserID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			m.Logger.Warn("пользователь не найден", "username", username)
			return "", fmt.Errorf("пользователь не найден")
		}
		m.Logger.Error("ошибка выполнения SQL-запроса при логине", "err", err)
		return "", err
	}

	if !CheckPasswordHash(password, passwordHash) {
		m.Logger.Info("Пользователь успешно аутентифицирован", "username", username)
		return UserID, nil
	}

	m.Logger.Info("пользователь успешно аутентифицирован", "username", username)
	return UserID, nil
}
