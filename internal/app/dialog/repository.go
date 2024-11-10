package dialog

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Repository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewRepository(db *sql.DB, logger *slog.Logger) *Repository {
	return &Repository{db, logger}
}

func (m *Repository) CreateDialog(UserIDOne, UserIDTwo string) error {
	var dialog Dialog

	stmt := `INSERT INTO dialogs (user_id_1, user_id_2) VALUES ($1, $2) ON CONFLICT (user_id_1, user_id_2) DO NOTHING RETURNING dialog_id`
	err := m.DB.QueryRow(stmt, UserIDOne, UserIDTwo).Scan(&dialog.DialogID)
	if err != nil {
		m.Logger.Error(
			"CreateDialog",
			"Ошибка добавления диалога в БД",
			"err:", err,
		)

		return err
	}

	return nil
}

func (m *Repository) DeleteDialog(DialogID string) error {
	stmt := `DELETE FROM dialogs
				WHERE dialog_id = $1`
	_, err := m.DB.Exec(stmt, DialogID)
	if err != nil {
		m.Logger.Error(
			"DeleteDialog",
			"Ошибка обработки запроса в БД",
			"err:", err,
		)

		return err
	}

	return nil
}

func (m *Repository) GetDialogs(UserID string) ([]Dialog, error) {
	stmt := `
SELECT 
    d.dialog_id, 
    d.user_id_1, 
    d.user_id_2, 
    COALESCE(d.dialog_avatar, '') AS dialog_avatar,
    d.last_message,
    u1.username AS user_one_username, 
    u2.username AS user_two_username
FROM dialogs d
LEFT JOIN users u1 ON d.user_id_1 = u1.id
LEFT JOIN users u2 ON d.user_id_2 = u2.id
WHERE d.user_id_1 = $1 OR d.user_id_2 = $1`
	fmt.Println("Executing query:", stmt)
	rows, err := m.DB.Query(stmt, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dialogs []Dialog
	for rows.Next() {
		var dil Dialog
		if err = rows.Scan(&dil.DialogID, &dil.UserIDOne, &dil.UserIDTwo, &dil.Avatar, &dil.LastMessage, &dil.UserOneUsername, &dil.UserTwoUsername); err != nil {
			m.Logger.Error("GetDialogs", "Ошибка при извлечении данных из базы", "error", err)
			return nil, err
		}

		dialogs = append(dialogs, Dialog{
			DialogID:        dil.DialogID,
			UserIDOne:       dil.UserIDOne,
			UserIDTwo:       dil.UserIDTwo,
			Avatar:          dil.Avatar,
			LastMessage:     dil.LastMessage,
			UserOneUsername: dil.UserOneUsername,
			UserTwoUsername: dil.UserTwoUsername,
		})
	}
	return dialogs, nil
}

func (m *Repository) GetDialog(UserID string) ([]Dialog, error) {
	var dialog Dialog

	stmt := `SELECT * FROM dialogs WHERE dialog_id = $1`
	rows, err := m.DB.Query(stmt, dialog.UserIDOne)
	if err != nil {
		m.Logger.Error(
			"GetDialog",
			"Не удалось обработать запрос БД",
		)

		return nil, err
	}
	defer rows.Close()

	var dialogs []Dialog
	for rows.Next() {
		var dialog Dialog
		err = rows.Scan(&UserID)
		if err != nil {
			m.Logger.Error(
				"GetDialog",
				"Ошибка при сканировании строки",
				"error", err,
			)

			return nil, err
		}

		dialogs = append(dialogs, dialog)

	}

	if err = rows.Err(); err != nil {
		m.Logger.Error(
			"PostGet",
			"Ошибка при итерации по строкам",
			"error", err,
		)
	}

	return dialogs, err
}
