package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"gofermart_/internal/models"
	"gofermart_/internal/service"
)

type mockAuthService struct {
	registerFn func(models.Credentials) (string, error)
	loginFn    func(models.Credentials) (string, error)
}

func (m *mockAuthService) Register(ctx context.Context, c models.Credentials) (string, error) {
	return m.registerFn(c)
}

func (m *mockAuthService) Login(ctx context.Context, c models.Credentials) (string, error) {
	return m.loginFn(c)
}

func TestRegister_InvalidJSON(t *testing.T) {

	h := NewAuthHandler(&mockAuthService{})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_InvalidCredentials(t *testing.T) {

	s := &mockAuthService{
		registerFn: func(c models.Credentials) (string, error) {
			return "", service.ErrInvalidCredentials
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "", Password: ""})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_UserExists(t *testing.T) {

	s := &mockAuthService{
		registerFn: func(c models.Credentials) (string, error) {
			return "", service.ErrUserExists
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "user", Password: "pass"})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_InternalError(t *testing.T) {

	s := &mockAuthService{
		registerFn: func(c models.Credentials) (string, error) {
			return "", errors.New("db error")
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "user", Password: "pass"})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRegister_Success(t *testing.T) {

	s := &mockAuthService{
		registerFn: func(c models.Credentials) (string, error) {
			return "token123", nil
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "user", Password: "pass"})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Register(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "Bearer token123", w.Header().Get("Authorization"))
}

func TestLogin_InvalidJSON(t *testing.T) {

	h := NewAuthHandler(&mockAuthService{})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	h.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {

	s := &mockAuthService{
		loginFn: func(c models.Credentials) (string, error) {
			return "", service.ErrInvalidCredentials
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "user", Password: "wrong"})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_InternalError(t *testing.T) {

	s := &mockAuthService{
		loginFn: func(c models.Credentials) (string, error) {
			return "", errors.New("db error")
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "user", Password: "pass"})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestLogin_Success(t *testing.T) {

	s := &mockAuthService{
		loginFn: func(c models.Credentials) (string, error) {
			return "token123", nil
		},
	}

	h := NewAuthHandler(s)

	body, _ := json.Marshal(models.Credentials{Login: "user", Password: "pass"})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	h.Login(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "Bearer token123", w.Header().Get("Authorization"))
}
