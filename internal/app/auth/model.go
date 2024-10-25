package auth

import (
	"errors"
	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/user"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	UserAuth *user.User `json:"user"`
}

func (a *Auth) Validate() error {
	if a.UserAuth.Username == "" {
		return errors.New("username is empty")
	}
	if a.UserAuth.Email == "" {
		return errors.New("email is empty")
	}
	if a.UserAuth.Password == "" {
		return errors.New("password is empty")
	}
	if len(a.UserAuth.Password) < 8 {
		return errors.New("password is too short")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
