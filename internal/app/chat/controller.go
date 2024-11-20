package chat

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type MessageRepo interface {
	CreateMessage(dialogID, senderID, recipientID, content string) (int, error)
	GetMessagesByDialogID(dialogID string) ([]Message, error)
}

type ControllerChat struct {
	repo      MessageRepo
	Logger    *slog.Logger
	WebSocket *WebSocketManager
}

func NewControllerChat(repo MessageRepo, logger *slog.Logger, websocket *WebSocketManager) *ControllerChat {
	return &ControllerChat{
		repo,
		logger,
		websocket,
	}
}

func (c *ControllerChat) SendMessage(w http.ResponseWriter, r *http.Request) {
	senderID, ok := r.Context().Value("user_id").(string)
	if !ok {
		c.Logger.Error("SendMessage", "Не удалось получить user_id")
		http.Error(w, "Не авторизован", http.StatusUnauthorized)
		return
	}

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		c.Logger.Error("SendMessage", "Ошибка декодирования JSON", "error", err)
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	msg.SenderID = senderID
	msg.Time = time.Now()

	messageID, err := c.repo.CreateMessage(msg.DialogID, msg.SenderID, msg.RecipientID, msg.Content)
	if err != nil {
		http.Error(w, "Ошибка сохранения сообщения", http.StatusInternalServerError)
		return
	}

	msg.ID = messageID

	if err = c.WebSocket.SendToUser(msg.RecipientID, msg); err != nil {
		c.Logger.Error("messageID:", messageID)
		http.Error(w, "Не удалось отправить сообщение", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (c *ControllerChat) GetMessages(w http.ResponseWriter, r *http.Request) {
	dialogID := r.URL.Query().Get("dialog_id")
	if dialogID == "" {
		http.Error(w, "dialog_id обязателен", http.StatusBadRequest)
		return
	}

	messages, err := c.repo.GetMessagesByDialogID(dialogID)
	if err != nil {
		http.Error(w, "Ошибка получения сообщений", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
