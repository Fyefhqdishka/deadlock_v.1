package user

import "time"

type User struct {
	UserID            string    `json:"user_id"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Gender            string    `json:"gender"`
	Dob               time.Time `json:"dob"`
	Avatar            string    `json:"avatar"`
	Time_registration time.Time `json:"time_registration"`
}

type Search struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
