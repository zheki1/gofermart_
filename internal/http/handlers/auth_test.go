package handlers

import (
	"gofermart_/internal/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthHandler_Register(t *testing.T) {
	auth.Init("test-secret")
	h := &AuthHandler{Repo: &mockRepo{}}

	req := httptest.NewRequest("POST", "/api/user/register",
		strings.NewReader(`{"login":"test","password":"123456"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.Register(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	authHeader := w.Header().Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Errorf("expected Authorization header with Bearer token, got %q", authHeader)
	}
}

func TestAuthHandler_Login(t *testing.T) {
	auth.Init("test-secret")
	h := &AuthHandler{Repo: &mockRepo{}}

	req := httptest.NewRequest("POST", "/api/user/login",
		strings.NewReader(`{"login":"test","password":"123456"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.Login(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}

	authHeader := w.Header().Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		t.Errorf("expected Authorization header with Bearer token, got %q", authHeader)
	}
}
