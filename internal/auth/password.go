package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword генерирует bcrypt-хеш для указанного пароля.
// Возвращает хеш в виде строки или ошибку.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword проверяет, соответствует ли пароль указанному хешу bcrypt.
// Возвращает nil, если пароль верный, иначе ошибку.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
