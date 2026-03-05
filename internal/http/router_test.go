package http

import (
	"context"
	"gofermart_/internal/auth"
	"gofermart_/internal/logger"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockRepo заглушка Repository для роутера
type mockRepo struct{}

func (m *mockRepo) CreateUser(ctx context.Context, u *models.User) error { return nil }
func (m *mockRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	return nil, storage.ErrUserNotFound
}
func (m *mockRepo) CreateOrder(ctx context.Context, o *models.Order) error { return nil }
func (m *mockRepo) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return nil, nil
}
func (m *mockRepo) ClaimOrders(ctx context.Context, limit int) ([]models.Order, error) {
	return nil, nil
}
func (m *mockRepo) UpdateOrderStatus(ctx context.Context, number string, status models.OrderStatus, accrual float64) error {
	return nil
}
func (m *mockRepo) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return &models.Balance{Current: 100, Withdrawn: 50}, nil
}
func (m *mockRepo) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	return nil
}
func (m *mockRepo) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return []models.Withdrawal{{Order: "12345678903", Sum: 10}}, nil
}

// Тест публичных маршрутов
func TestRouter_PublicRoutes(t *testing.T) {
	auth.Init("test-secret")
	logger.Init("/dev/null", logger.DEBUG)
	r := NewRouter(&mockRepo{})

	tests := []struct {
		method string
		path   string
		want   int
	}{
		{"POST", "/api/user/register", http.StatusOK},
		{"POST", "/api/user/login", http.StatusUnauthorized}, // mock возвращает ErrUserNotFound
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(`{"login":"user","password":"pass"}`))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != tt.want {
			t.Errorf("%s %s: expected %d, got %d", tt.method, tt.path, tt.want, w.Code)
		}
	}
}

func TestRouter_ProtectedRoutes_Authorized(t *testing.T) {
	auth.Init("test-secret")
	logger.Init("/dev/null", logger.DEBUG)
	r := NewRouter(&mockRepo{})

	// Генерируем JWT для userID=1
	token, err := auth.GenerateToken(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		method string
		path   string
		want   int
		body   string
	}{
		{"GET", "/api/user/balance", http.StatusOK, ""},
		{"POST", "/api/user/balance/withdraw", http.StatusOK, `{"order":"79927398713","sum":10}`},
		{"GET", "/api/user/withdrawals", http.StatusOK, ""},
		{"POST", "/api/user/orders", http.StatusAccepted, "79927398713"},
		{"GET", "/api/user/orders", http.StatusNoContent, ""},
		{"GET", "/chin", http.StatusOK, ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != tt.want {
			t.Errorf("%s %s: expected %d, got %d", tt.method, tt.path, tt.want, w.Code)
		}
	}
}

// Тест защищённых маршрутов без авторизации
func TestRouter_ProtectedRoutes_Unauthorized(t *testing.T) {
	logger.Init("/dev/null", logger.DEBUG)
	r := NewRouter(&mockRepo{})

	tests := []struct {
		method string
		path   string
		want   int
	}{
		{"GET", "/api/user/balance", http.StatusUnauthorized},
		{"POST", "/api/user/balance/withdraw", http.StatusUnauthorized},
		{"GET", "/api/user/withdrawals", http.StatusUnauthorized},
		{"POST", "/api/user/orders", http.StatusUnauthorized},
		{"GET", "/api/user/orders", http.StatusUnauthorized},
		{"GET", "/chin", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != tt.want {
			t.Errorf("%s %s: expected %d, got %d", tt.method, tt.path, tt.want, w.Code)
		}
	}
}
