package chat

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Time     string `json:"time"`
	Typing   bool   `json:"typing"`
}
