package chat

import (
	"database/sql"
	"log/slog"
	"time"
)

type MessageRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewChatRepository(db *sql.DB, logger *slog.Logger) *MessageRepository {
	return &MessageRepository{
		db,
		logger,
	}
}

func (mr *MessageRepository) CreateMessage(dialogID, senderID, recipientID, content string) (int, error) {
	query := `
		INSERT INTO messages (dialog_id, sender_id, recipient_id, content, timestamp)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var messageID int
	err := mr.DB.QueryRow(query, dialogID, senderID, recipientID, content, time.Now()).Scan(&messageID)
	return messageID, err
}

func (mr *MessageRepository) GetMessagesByDialogID(dialogID string) ([]Message, error) {
	query := `
		SELECT id, dialog_id, sender_id, recipient_id, content, timestamp
		FROM messages
		WHERE dialog_id = $1
		ORDER BY timestamp ASC
	`
	rows, err := mr.DB.Query(query, dialogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.DialogID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.Time); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
