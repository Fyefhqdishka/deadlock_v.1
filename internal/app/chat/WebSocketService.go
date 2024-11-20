package chat

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Проверка источника, здесь можно настроить проверку
		// на допустимые источники запросов, чтобы не было CORS ошибок
		return true
	},
}

type WebSocketManager struct {
	clients map[string]*websocket.Conn // user_id -> WebSocket connection
	mu      sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[string]*websocket.Conn),
	}
}

func (wm *WebSocketManager) Register(userID string, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.clients[userID] = conn
}

func (wm *WebSocketManager) SendToUser(userID string, msg Message) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	conn, ok := wm.clients[userID]
	if !ok {
		return errors.New("пользователь не подключен")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

func (wsManager *WebSocketManager) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем WebSocket-соединение
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установлении WebSocket соединения:", err)
		return
	}
	defer conn.Close()

	// Добавляем клиента в список
	wsManager.clients[conn] = true

	for {
		// Чтение сообщений от клиента
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Ошибка при чтении сообщения:", err)
			delete(wsManager.clients, conn)
			break
		}

		// Отправка сообщения всем подключенным клиентам
		for client := range wsManager.clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Ошибка при отправке сообщения:", err)
				client.Close()
				delete(wsManager.clients, client)
			}
		}
	}
}
