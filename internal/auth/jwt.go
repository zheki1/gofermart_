package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// secret — ключ для подписи JWT токенов.
// TODO вынести в конфиг или переменные окружения.
var secret = []byte("super-secret-key")

// GenerateToken создает JWT токен для пользователя с указанным userID.
// Токен действителен 24 часа.
// Возвращает строку токена или ошибку при создании.
func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"uid": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken проверяет JWT токен и извлекает userID.
// Возвращает userID и ошибку, если токен недействителен.
func ParseToken(tokenStr string) (int, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !t.Valid {
		return 0, err
	}

	claims := t.Claims.(jwt.MapClaims)
	return int(claims["uid"].(float64)), nil
}
