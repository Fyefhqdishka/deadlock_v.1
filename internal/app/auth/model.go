package auth

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Fyefhqdishka/deadlock_v.1/internal/app/user"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 8

type Auth struct {
	UserAuth user.User `json:"user"`
}

// логика проверки валидности пароля с почтой
func (a *Auth) Validate() error {
	var validationErrors []string

	if a.UserAuth.Username == "" {
		validationErrors = append(validationErrors, "username is empty")
	}

	if a.UserAuth.Email == "" {
		validationErrors = append(validationErrors, "email is empty")
	} else if !isValidEmail(a.UserAuth.Email) {
		validationErrors = append(validationErrors, "email format is invalid")
	}

	if a.UserAuth.Password == "" {
		validationErrors = append(validationErrors, "password is empty")
	} else if len(a.UserAuth.Password) < minPasswordLength {
		validationErrors = append(validationErrors, fmt.Sprintf("password is too short, minimum length is %d", minPasswordLength))
	}

	if len(validationErrors) > 0 {
		return errors.New(fmt.Sprintf("validation errors: %v", validationErrors))
	}

	return nil
}

// создает хэш для пароля.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %v", err)
	}
	return string(hashedBytes), nil
}

// проверяет, соответствует ли пароль хэшу
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}
