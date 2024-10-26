package auth

import (
	"encoding/json"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/user"
	"log/slog"
	"net/http"
)

type RepositoryAuth interface {
	RegistrationUser(user user.User) error
	LoginUser(username, password string) (string, error)
}

type ControllerAuth struct {
	Repo   RepositoryAuth
	Logger *slog.Logger
}

func NewControllerAuth(repo RepositoryAuth, logger *slog.Logger) *ControllerAuth {
	return &ControllerAuth{
		Repo:   repo,
		Logger: logger,
	}
}

func (c *ControllerAuth) Register(w http.ResponseWriter, r *http.Request) {
	c.Logger.Debug("начала обработки регистрации пользователя")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user *Auth
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		c.Logger.Error("ошибка при декодировании JSON", "err", err)
		return
	}

	if err = user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.Logger.Error("ошибка валидации данных пользователя", "err", err)
		return
	}

	hashedPassword, err := HashPassword(user.UserAuth.Password)
	if err != nil {
		c.Logger.Error("ошибка при хэшировании пароля", "err:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.UserAuth.Password = hashedPassword

	if err = c.Repo.RegistrationUser(user.UserAuth); err != nil {
		c.Logger.Error("ошибка при сохранении пользователя в репозитории", "err", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	c.Logger.Debug("User successfully registered", "username", user.UserAuth.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set status code before sending response
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
