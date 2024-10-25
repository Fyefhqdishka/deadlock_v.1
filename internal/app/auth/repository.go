package auth

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Repository struct {
	DB     *sql.DB
	logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{
		DB:     db,
		logger: logger,
	}
}

func (m *Repository) RegistrationUser(auth *Auth) error {
	m.logger.Debug("starting registration user")

	stmt, err := m.DB.Prepare("INSERT INTO users (name, username, email, password, gender, dob) VALUES ($1,$2,$3,$4,$5,$6)")
	if err != nil {
		m.logger.Error("error preparing statement")
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(auth.UserAuth.Name, auth.UserAuth.Username, auth.UserAuth.Email, auth.UserAuth.Password, auth.UserAuth.Gender, auth.UserAuth.Dob)
	if err != nil {
		m.logger.Error("error executing statement")
		return err
	}

	m.logger.Debug("registration user finished")
	return nil
}

func (m *Repository) LoginUser(username, password string) (string, error) {
	m.logger.Debug("starting login user")

	var UserID string
	var passwordHash string

	stmt := `SELECT id, password FROM users WHERE username = $1`
	err := m.DB.QueryRow(stmt, username).Scan(&UserID, &passwordHash)
	if err != nil {
		m.logger.Error("error executing statement")
		return "", err
	}

	if CheckPasswordHash(password, passwordHash) {
		m.logger.Info("Пользователь успешно аутентифицирован", "username", username)
		return UserID, nil
	} else {
		m.logger.Warn("Неверный пароль при аутентификации", "username", username)
		return "", fmt.Errorf("неверный пароль")
	}
}
