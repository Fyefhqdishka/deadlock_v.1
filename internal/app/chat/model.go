package chat

type Message struct {
	UserID   string `json:"user_id"`
	DialogID string `json:"dialog_id"`
	Message  string `json:"message"`
	Time     string `json:"time"`
	Typing   bool   `json:"typing"`
}
