package middleware

import (
	"context"
	"net/http"
	"strings"

	"gofermart_/internal/auth"
)

// ctxKey используется для хранения значений в контексте запроса
type ctxKey string

// UserIDKey ключ для userID в контексте
const UserIDKey ctxKey = "uid"

// Auth проверяет JWT токен из заголовка Authorization.
// Если токен валиден, добавляет userID в контекст запроса.
// В противном случае возвращает 401 Unauthorized.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		uid, err := auth.ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
