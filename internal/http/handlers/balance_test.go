package handlers

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"gofermart_/internal/http/middleware"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// // Withdraw тесты
// func TestBalanceHandler_Withdraw(t *testing.T) {
// 	h := &BalanceHandler{Repo: &mockRepo{}}

// 	tests := []struct {
// 		name       string
// 		order      string
// 		sum        float64
// 		ctxUserID  interface{}
// 		wantStatus int
// 	}{
// 		{"ok", "4532015112830366", 10, 1, http.StatusOK},                 // валидный Luhn
// 		{"not enough", "79927398713", 10, 1, http.StatusPaymentRequired}, // валидный Luhn
// 		{"invalid sum", "4532015112830366", 0, 1, http.StatusBadRequest},
// 		{"invalid order", "123", 10, 1, http.StatusUnprocessableEntity},
// 		{"unauthorized", "4532015112830366", 10, nil, http.StatusUnauthorized},
// 	}

// 	for _, tt := range tests {
// 		body := map[string]interface{}{"order": tt.order, "sum": tt.sum}
// 		b, _ := json.Marshal(body)
// 		req := httptest.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewReader(b))
// 		if tt.ctxUserID != nil {
// 			req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, tt.ctxUserID))
// 		}
// 		w := httptest.NewRecorder()

// 		h.Withdraw(w, req)

// 		if w.Code != tt.wantStatus {
// 			t.Errorf("%s: expected %d, got %d", tt.name, tt.wantStatus, w.Code)
// 		}
// 	}
// }

// // Get тесты
// func TestBalanceHandler_Get(t *testing.T) {
// 	h := &BalanceHandler{Repo: &mockRepo{}}

// 	req := httptest.NewRequest("GET", "/api/user/balance", nil)
// 	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 1))
// 	w := httptest.NewRecorder()

// 	h.Get(w, req)

// 	if w.Code != http.StatusOK {
// 		t.Fatalf("expected 200, got %d", w.Code)
// 	}
// }

// // Withdrawals тесты
// func TestBalanceHandler_Withdrawals(t *testing.T) {
// 	h := &BalanceHandler{Repo: &mockRepo{}}

// 	tests := []struct {
// 		name       string
// 		ctxUserID  interface{}
// 		wantStatus int
// 	}{
// 		{"with withdrawals", 1, http.StatusOK},
// 		{"no withdrawals", 2, http.StatusNoContent},
// 		{"unauthorized", nil, http.StatusUnauthorized},
// 	}

// 	for _, tt := range tests {
// 		req := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
// 		if tt.ctxUserID != nil {
// 			req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, tt.ctxUserID))
// 		}
// 		w := httptest.NewRecorder()

// 		h.Withdrawals(w, req)

// 		if w.Code != tt.wantStatus {
// 			t.Errorf("%s: expected %d, got %d", tt.name, tt.wantStatus, w.Code)
// 		}
// 	}
// }
