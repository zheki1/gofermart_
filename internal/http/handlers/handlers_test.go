package handlers

import (
	"context"
	"gofermart_/internal/auth"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
	"time"
)

type mockRepo struct{}

func (m *mockRepo) CreateUser(ctx context.Context, u *models.User) error {
	u.ID = 1
	return nil
}

func (m *mockRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	hashed, _ := auth.HashPassword("123456")
	return &models.User{
		ID:           1,
		Login:        login,
		PasswordHash: hashed,
	}, nil
}

func (m *mockRepo) CreateOrder(ctx context.Context, o *models.Order) error {
	switch o.Number {
	case "4532015112830366":
		return storage.ErrOrderExistsForUser
	case "4485275742308327":
		return storage.ErrOrderExistsForOther
	}
	return nil
}

func (m *mockRepo) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	if userID == 1 {
		return []models.Order{
			{
				Number:     "12345678903",
				UserID:     1,
				Status:     models.OrderNew,
				Accrual:    0,
				UploadedAt: time.Now(),
			},
		}, nil
	}
	return []models.Order{}, nil
}

func (m *mockRepo) ClaimOrders(ctx context.Context, limit int) ([]models.Order, error) {
	return nil, nil
}
func (m *mockRepo) UpdateOrderStatus(ctx context.Context, number string, status models.OrderStatus, accrual float64) error {
	return nil
}

func (m *mockRepo) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	switch order {
	case "79927398713": // валидный номер для теста "not enough"
		return storage.ErrNotEnoughFunds
	}
	return nil
}

func (m *mockRepo) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return &models.Balance{Current: 100, Withdrawn: 50}, nil
}

func (m *mockRepo) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	if userID == 1 {
		return []models.Withdrawal{
			{Order: "12345678903", Sum: 10, ProcessedAt: time.Now()},
		}, nil
	}
	return []models.Withdrawal{}, nil
}
