package http

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"gofermart_/internal/logger"
	"gofermart_/internal/models"
)

type mockRepo struct{}

func (m *mockRepo) CreateUser(ctx context.Context, u *models.User) error { return nil }
func (m *mockRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	return nil, nil
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
func (m *mockRepo) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	return nil
}
func (m *mockRepo) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return &models.Balance{}, nil
}
func (m *mockRepo) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return nil, nil
}

func TestRouter_RoutesExist(t *testing.T) {
	logger.Init("/dev/null", logger.DEBUG)
	repo := &mockRepo{}
	container := NewContainer(repo)
	router := NewRouter(container)

	routes := []struct {
		method string
		route  string
	}{
		{"POST", "/api/user/register"},
		{"POST", "/api/user/login"},
		{"POST", "/api/user/orders"},
		{"GET", "/api/user/orders"},
		{"GET", "/api/user/balance"},
		{"POST", "/api/user/balance/withdraw"},
		{"GET", "/api/user/withdrawals"},
		{"GET", "/chin"},
	}

	for _, rt := range routes {
		req := httptest.NewRequest(rt.method, rt.route, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.NotEqual(t, 404, rec.Code, "route %s %s should exist", rt.method, rt.route)
	}
}
