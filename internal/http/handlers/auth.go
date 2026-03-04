package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gofermart_/internal/auth"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
)

// AuthHandler предоставляет методы для регистрации и аутентификации пользователей.
type AuthHandler struct {
	Repo storage.Repository
}

// Register обрабатывает запрос на регистрацию нового пользователя.
// JSON тело запроса должно содержать login и password.
// Если регистрация успешна, возвращает HTTP 200 и JWT в заголовке Authorization.
// Возможные ошибки:
// - 400 Bad Request — некорректный формат запроса или пустой логин/пароль
// - 409 Conflict — логин уже существует
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var c models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if c.Login == "" || c.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(c.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	u := &models.User{
		Login:        c.Login,
		PasswordHash: hash,
	}

	if err := h.Repo.CreateUser(r.Context(), u); err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			http.Error(w, "login already exists", http.StatusConflict)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

// Register обрабатывает запрос на регистрацию нового пользователя.
// JSON тело запроса должно содержать login и password.
// Если регистрация успешна, возвращает HTTP 200 и JWT в заголовке Authorization.
// Возможные ошибки:
// - 400 Bad Request — некорректный формат запроса или пустой логин/пароль
// - 409 Conflict — логин уже существует
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var c models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	if c.Login == "" || c.Password == "" {
		http.Error(w, "login and password required", http.StatusBadRequest)
		return
	}

	u, err := h.Repo.GetUserByLogin(r.Context(), c.Login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			http.Error(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := auth.CheckPassword(u.PasswordHash, c.Password); err != nil {
		http.Error(w, "invalid login or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
