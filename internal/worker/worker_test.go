package worker

import (
	"context"
	"testing"
	"time"

	"gofermart_/internal/accrual"
	"gofermart_/internal/logger"
	"gofermart_/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ----------- Моки ------------

type mockRepo struct {
	mock.Mock
}

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
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *mockRepo) UpdateOrderStatus(ctx context.Context, number string, status models.OrderStatus, accrual float64) error {
	args := m.Called(ctx, number, status, accrual)
	return args.Error(0)
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

type MockAccrualClient struct {
	mock.Mock
}

func (m *MockAccrualClient) GetOrder(number string) (*accrual.Response, error) {
	args := m.Called(number)
	return args.Get(0).(*accrual.Response), args.Error(1)
}

// --- tests ---

func TestWorker_handleOrder(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	t.Run("PROCESSED order updates accrual", func(t *testing.T) {
		logger.Init("/dev/null", logger.DEBUG)
		repo := &mockRepo{}
		acClient := &MockAccrualClient{}
		w := New(repo, acClient)

		order := models.Order{Number: "123"}

		acClient.On("GetOrder", order.Number).Return(&accrual.Response{
			Order:   order.Number,
			Status:  "PROCESSED",
			Accrual: func() *float64 { v := 100.0; return &v }(),
		}, nil)

		repo.On("UpdateOrderStatus", ctx, order.Number, models.OrderProcessed, 100.0).Return(nil)

		w.handleOrder(ctx, order)

		acClient.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("INVALID order updates status to invalid", func(t *testing.T) {
		repo := &mockRepo{}
		acClient := &MockAccrualClient{}
		w := New(repo, acClient)

		order := models.Order{Number: "456"}

		acClient.On("GetOrder", order.Number).Return(&accrual.Response{
			Order:  order.Number,
			Status: "INVALID",
		}, nil)

		repo.On("UpdateOrderStatus", ctx, order.Number, models.OrderInvalid, 0.0).Return(nil)

		w.handleOrder(ctx, order)

		acClient.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("PROCESSING order does not update status", func(t *testing.T) {
		repo := &mockRepo{}
		acClient := &MockAccrualClient{}
		w := New(repo, acClient)

		order := models.Order{Number: "789"}

		acClient.On("GetOrder", order.Number).Return(&accrual.Response{
			Order:  order.Number,
			Status: "PROCESSING",
		}, nil)

		w.handleOrder(ctx, order)

		acClient.AssertExpectations(t)
	})

	t.Run("retry on error and respect context cancellation", func(t *testing.T) {
		repo := &mockRepo{}
		acClient := &MockAccrualClient{}
		w := New(repo, acClient)

		order := models.Order{Number: "999"}

		// Возвращаем валидный Response и ошибку, чтобы сработал retry
		acClient.On("GetOrder", order.Number).Return(&accrual.Response{
			Order:  order.Number,
			Status: "NEW",
		}, assert.AnError)

		// контекст с коротким таймаутом для остановки retry
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		w.handleOrder(ctx, order)

		acClient.AssertExpectations(t)
	})
}
