package auth

import (
	"encoding/json"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/user"
	"github.com/Fyefhqdishka/deadlock_v.1/pkg/jwt"
	"io"
	"log/slog"
	"net/http"
	"time"
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
	c.Logger.Debug(
		"Register",
		"начала обработки регистрации пользователя",
	)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user *Auth
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		c.Logger.Error(
			"Register",
			"ошибка при декодировании JSON",
			"err", err,
		)
		return
	}

	if err = user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.Logger.Error(
			"Register",
			"ошибка валидации данных пользователя",
			"err", err,
		)
		return
	}

	hashedPassword, err := HashPassword(user.UserAuth.Password)
	if err != nil {
		c.Logger.Error(
			"Register",
			"ошибка при хэшировании пароля",
			"err:", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.UserAuth.Password = hashedPassword

	if err = c.Repo.RegistrationUser(*user.UserAuth); err != nil {
		c.Logger.Error(
			"Register",
			"ошибка при сохранении пользователя в репозитории",
			"err", err,
		)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	c.Logger.Debug(
		"Register",
		"User successfully registered",
		"username", user.UserAuth.Username,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (c *ControllerAuth) Login(w http.ResponseWriter, r *http.Request) {
	c.Logger.Info(
		"login",
		"начало обработки Авторизации",
	)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		c.Logger.Error(
			"Login",
			"ошибка чтения тела запроса",
			"error", err,
		)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var user *Auth
	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		c.Logger.Error(
			"Login",
			"ошибка декодирования JSON",
			"error", err,
		)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	c.Logger.Debug(
		"Login",
		"Перед вызовом LoginUser",
		"username", user.UserAuth.Username,
	)

	userID, err := c.Repo.LoginUser(user.UserAuth.Username, user.UserAuth.Password)
	if err != nil {
		c.Logger.Error(
			"Login",
			"ошибка при аутентификации пользователя",
			"username", user.UserAuth.Username,
			"password", user.UserAuth.Password,
			"error", err,
		)

		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	tokenStr, err := jwt.GenerateToken(userID)
	if err != nil {
		c.Logger.Error(
			"Login",
			"Ошибка генерации токена",
			"error", err,
		)

		http.Error(w, "Не удалось сгенерировать токен", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenStr,
		Expires:  time.Now().Add(36 * time.Hour),
		HttpOnly: true,
	})

	c.Logger.Info(
		"Login",
		"Пользователь успешно аутентифицирован",
		"username", user.UserAuth.Username,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User authenticate successfully",
		"token":   tokenStr,
	})
}
