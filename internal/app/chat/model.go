package chat

import (
	"github.com/gorilla/websocket"
	"time"
)

type Message struct {
	ID          int       `json:"id"`
	DialogID    string    `json:"dialog_id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	Content     string    `json:"content"`
	Time        time.Time `json:"timestamp"`
}

type Clients struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Conn     *websocket.Conn
}
