package chat

import (
	"database/sql"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type ControllerChat struct {
	DB     *sql.DB
	Logger *slog.Logger
}

var (
	upgrader  = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients   = make(map[*websocket.Conn]string)          // Map client connection to userID
	dialogs   = make(map[string]map[*websocket.Conn]bool) // Map dialog_id to WebSocket connections
	broadcast = make(chan Message)
	mu        sync.Mutex
)

func (c *ControllerChat) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.Logger.Error(
			"HandleConnections",
			"Ошибка при обновлении до WebSocket",
			"error", err,
		)
		return
	}
	defer func() {
		mu.Lock()
		delete(clients, ws)
		for dialogID := range dialogs {
			delete(dialogs[dialogID], ws)
		}
		mu.Unlock()
		ws.Close()
	}()

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		c.Logger.Error(
			"HandleConnections",
			"Ошибка: userID не найден в контексте",
		)
		return
	}

	c.Logger.Info(
		"HandleConnections",
		"Пользователь подключился UserID:", userID,
	)

	mu.Lock()
	clients[ws] = userID
	mu.Unlock()

	// Listen for incoming messages
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			c.Logger.Warn(
				"HandleConnections",
				"Ошибка чтения JSON от клиента",
				"error", err,
			)
			break
		}
		msg.UserID = userID

		// Ensure dialog_id exists, create it if not
		mu.Lock()
		if _, exists := dialogs[msg.DialogID]; !exists {
			dialogs[msg.DialogID] = make(map[*websocket.Conn]bool)
		}
		dialogs[msg.DialogID][ws] = true
		mu.Unlock()

		broadcast <- msg
	}
}

func (c *ControllerChat) HandleMessages() {
	for {
		msg := <-broadcast

		// Send the message only to clients in the same dialog
		mu.Lock()
		clientsInDialog, exists := dialogs[msg.DialogID]
		if exists {
			for client := range clientsInDialog {
				err := client.WriteJSON(msg)
				if err != nil {
					c.Logger.Warn(
						"HandleMessages",
						"Ошибка отправки сообщения клиенту",
						"error", err,
					)

					client.Close()
					delete(clientsInDialog, client)
				}
			}
		}
		mu.Unlock()
	}
}
