package middleware

import (
	"gofermart_/internal/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест простого прохождения запроса
func TestLoggingMiddleware_OK(t *testing.T) {
	logger.Init("/dev/null", logger.DEBUG)
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("GET", "/ok", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", w.Code)
	}
}

// Тест паники
func TestLoggingMiddleware_Panic(t *testing.T) {
	logger.Init("/dev/null", logger.DEBUG)
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 Internal Server Error, got %d", w.Code)
	}
}

// Тест 4xx ошибки
func TestLoggingMiddleware_ClientError(t *testing.T) {
	logger.Init("/dev/null", logger.DEBUG)
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))

	req := httptest.NewRequest("GET", "/bad", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
