package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var jwtSecret = []byte("jwt-secret")

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(id string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	return claims, nil
}

func SetJWTCookie(w http.ResponseWriter, userID string) error {
	token, err := GenerateToken(userID)
	if err != nil {
		return err
	}

	// Установка cookie с JWT токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id", // Имя cookie
		Value:    token,     // Значение (сам токен)
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteNoneMode, // Для кросс-доменных запросов
	})

	return nil
}
