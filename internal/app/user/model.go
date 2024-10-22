package user

import "time"

type User struct {
	UUID              string    `json:"uuid"`
	Name              string    `json:"name"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Gender            string    `json:"gender"`
	Dob               time.Time `json:"dob"`
	Time_registration time.Time `json:"time_registration"`
}
