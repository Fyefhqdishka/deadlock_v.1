package post

type Post struct {
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
	Username string `json:"username"`
	ID       int    `json:"id"`
}
