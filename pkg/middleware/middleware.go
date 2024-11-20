package middleware

import (
	"context"
	"github.com/Fyefhqdishka/deadlock_v.1/pkg/jwt"
	"log/slog"
	"net/http"
)

func JWTMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Cookie")

			cookie, err := r.Cookie("token")
			if err != nil {
				logger.Warn("Cookie не найден", "error", err, "cookies", r.Cookies())
				http.Error(w, "Требуется аутентификация", http.StatusUnauthorized)
				return
			}

			logger.Info("Cookie найден", "token", cookie.Value)

			tokenStr := cookie.Value
			claims, err := jwt.VerifyJWT(tokenStr)
			if err != nil {
				logger.Warn("Недействительный токен", "error", err)
				http.Error(w, "Недействительный токен", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.ID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
