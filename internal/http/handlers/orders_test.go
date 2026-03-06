package handlers

// import (
// 	"context"
// 	"gofermart_/internal/http/middleware"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// // тест Upload
// func TestOrderHandler_Upload(t *testing.T) {
// 	h := &OrderHandler{Repo: &mockRepo{}}

// 	tests := []struct {
// 		name       string
// 		orderNum   string
// 		ctxUserID  interface{}
// 		wantStatus int
// 	}{
// 		{"new order", "79927398713", 1, http.StatusAccepted},
// 		{"exists for user", "4532015112830366", 1, http.StatusOK},        // валидный Luhn
// 		{"exists for other", "4485275742308327", 1, http.StatusConflict}, // валидный Luhn
// 		{"invalid Luhn", "123", 1, http.StatusUnprocessableEntity},
// 		{"empty order", "", 1, http.StatusBadRequest},
// 		{"unauthorized", "79927398713", nil, http.StatusUnauthorized},
// 	}

// 	for _, tt := range tests {
// 		req := httptest.NewRequest("POST", "/api/user/orders", strings.NewReader(tt.orderNum))
// 		if tt.ctxUserID != nil {
// 			req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, tt.ctxUserID))
// 		}
// 		w := httptest.NewRecorder()

// 		h.Upload(w, req)

// 		if w.Code != tt.wantStatus {
// 			t.Errorf("%s: expected %d, got %d", tt.name, tt.wantStatus, w.Code)
// 		}
// 	}
// }

// // тест List
// func TestOrderHandler_List(t *testing.T) {
// 	h := &OrderHandler{Repo: &mockRepo{}}

// 	tests := []struct {
// 		name       string
// 		ctxUserID  interface{}
// 		wantStatus int
// 	}{
// 		{"with orders", 1, http.StatusOK},
// 		{"no orders", 2, http.StatusNoContent},
// 		{"unauthorized", nil, http.StatusUnauthorized},
// 	}

// 	for _, tt := range tests {
// 		req := httptest.NewRequest("GET", "/api/user/orders", nil)
// 		if tt.ctxUserID != nil {
// 			req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, tt.ctxUserID))
// 		}
// 		w := httptest.NewRecorder()

// 		h.List(w, req)

// 		if w.Code != tt.wantStatus {
// 			t.Errorf("%s: expected %d, got %d", tt.name, tt.wantStatus, w.Code)
// 		}
// 	}
// }
