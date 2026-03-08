package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"gofermart_/internal/helpers"
	"gofermart_/internal/models"
	"gofermart_/internal/service"
)

// AuthHandler предоставляет методы для регистрации и аутентификации пользователей.
type AuthHandler struct {
	Service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
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
		helpers.WriteJSONError(w, "invalid request format", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Register(r.Context(), c)
	if err != nil {

		if errors.Is(err, service.ErrInvalidCredentials) {
			helpers.WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrUserExists) {
			helpers.WriteJSONError(w, err.Error(), http.StatusConflict)
			return
		}

		helpers.WriteJSONError(w, "internal error", http.StatusInternalServerError)
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
		helpers.WriteJSONError(w, "invalid request format", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(r.Context(), c)
	if err != nil {

		if errors.Is(err, service.ErrInvalidCredentials) {
			helpers.WriteJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrInvalidCredentials) {
			helpers.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		helpers.WriteJSONError(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
