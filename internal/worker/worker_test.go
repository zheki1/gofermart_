package worker

import (
	"context"
	"testing"
	"time"

	"gofermart_/internal/accrual"
	"gofermart_/internal/logger"
	"gofermart_/internal/models"
)

// ----------- Моки ------------

type mockRepo struct{}

func (m *mockRepo) CreateUser(ctx context.Context, u *models.User) error {
	return nil
}
func (m *mockRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	return &models.User{ID: 1}, nil
}
func (m *mockRepo) CreateOrder(ctx context.Context, o *models.Order) error { return nil }
func (m *mockRepo) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return []models.Order{}, nil
}
func (m *mockRepo) ClaimOrders(ctx context.Context, limit int) ([]models.Order, error) {
	return []models.Order{
		{Number: "12345678903", Status: models.OrderNew, UserID: 1},
	}, nil
}
func (m *mockRepo) UpdateOrderStatus(ctx context.Context, number string, status models.OrderStatus, accrual float64) error {
	return nil
}
func (m *mockRepo) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return &models.Balance{Current: 100, Withdrawn: 0}, nil
}
func (m *mockRepo) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	return nil
}
func (m *mockRepo) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return nil, nil
}

// ----------- Мок accrual клиента -----------

type mockAccrual struct{}

func (c *mockAccrual) GetOrder(number string) (*accrual.Response, int, time.Duration, error) {
	accr := 100.0
	return &accrual.Response{
		Order:   number,
		Status:  "PROCESSED",
		Accrual: &accr,
	}, 200, 0, nil
}

// ----------- Тесты Worker -----------

func TestWorker_ProcessOrder(t *testing.T) {
	logger.Init("/dev/null", logger.DEBUG)
	repo := &mockRepo{}
	client := &mockAccrual{}

	worker := New(repo, client)
	worker.concurrency = 1
	worker.interval = 50 * time.Millisecond // чтобы быстро тикнул

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
}

func TestWorker_HandleOrder_InvalidStatus(t *testing.T) {
	repo := &mockRepo{}
	client := &mockAccrualInvalid{}

	w := New(repo, client)

	order := models.Order{Number: "11111111111", Status: models.OrderNew, UserID: 1}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	w.handleOrder(ctx, order)
}

// мок для INVALID заказа
type mockAccrualInvalid struct{}

func (c *mockAccrualInvalid) GetOrder(number string) (*accrual.Response, int, time.Duration, error) {
	return &accrual.Response{
		Order:  number,
		Status: "INVALID",
	}, 200, 0, nil
}
