package dialog

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type RepositoryDialog interface {
	CreateDialog(UserIDOne, UserIDTwo string) error
	GetDialogs(UserID string) ([]Dialog, error)
}

type DialogCreator interface {
	CreateDialog(UserIDOne, UserIDTwo string) error
}

type DialogGetter interface {
	GetDialogs(UserID string) ([]Dialog, error)
}

type ControllerDialog struct {
	dialogCreator DialogCreator
	dialogGetter  DialogGetter
	logger        *slog.Logger
}

func NewControllerDialog(dialogCreator DialogCreator, dialogGetter DialogGetter, logger *slog.Logger) *ControllerDialog {
	return &ControllerDialog{
		dialogCreator,
		dialogGetter,
		logger,
	}
}

func (c *ControllerDialog) Create(w http.ResponseWriter, r *http.Request) {
	c.logger.Debug(
		"Create",
		"Starting create dialog",
	)

	var dialog Dialog

	err := json.NewDecoder(r.Body).Decode(&dialog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		c.logger.Error(
			"Create",
			"ошибка декодирования json",
			"error:", err,
		)

		return
	}

	userID := r.Context().Value("user_id")
	if userID == nil {
		c.logger.Warn(
			"Create",
			"userID не найден в контексте",
		)

		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.logger.Warn(
			"Create",
			"userID не является строкой",
		)

		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	err = c.dialogCreator.CreateDialog(userIDStr, dialog.UserIDTwo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.logger.Debug(
		"Create",
		"finished create dialog",
		"User_1:", userIDStr,
		"User_2:", dialog.UserIDTwo,
	)

	w.WriteHeader(http.StatusCreated)
	dialogResponse := struct {
		UserIDTwo string `json:"user_id_two"`
	}{
		UserIDTwo: dialog.UserIDTwo,
	}

	if err := json.NewEncoder(w).Encode(dialogResponse); err != nil {
		c.logger.Error(
			"Create",
			"ошибка при отправке ответа",
			"error", err,
		)

		return
	}
}

func (c *ControllerDialog) Get(w http.ResponseWriter, r *http.Request) {
	UserID := r.Context().Value("user_id")
	if UserID == nil {
		c.logger.Error(
			"Get",
			"userID не найден в контексте",
		)
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	UserIDStr, ok := UserID.(string)
	if !ok {
		c.logger.Warn(
			"Get",
			"userID не является строкой",
		)
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	c.logger.Debug(
		"Get",
		"finished fetching dialogs",
		"user_id:", UserIDStr,
	)

	dialogs, err := c.dialogGetter.GetDialogs(UserIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(dialogs); err != nil {
		c.logger.Error(
			"Get",
			"ошибка при отправке ответа",
			"error", err,
		)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
	}
}
