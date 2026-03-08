package storage

import (
	"context"
	"errors"

	"gofermart_/internal/models"
)

// Стандартные ошибки пакета storage
var (
	ErrUserExists            = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrNotEnoughFunds        = errors.New("not enough funds")
	ErrOrderAlreadyWithdrawn = errors.New("order already withdrawn")
	ErrOrderExistsForUser    = errors.New("order already uploaded by this user")
	ErrOrderExistsForOther   = errors.New("order uploaded by another user")
)

// Repository описывает интерфейс хранения данных приложения.
type Repository interface {
	// auth
	CreateUser(ctx context.Context, u *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)

	// orders
	CreateOrder(ctx context.Context, o *models.Order) error
	GetUserOrders(ctx context.Context, userID int) ([]models.Order, error)
	ClaimOrders(ctx context.Context, limit int) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status models.OrderStatus, accrual float64) error

	// balance / withdrawals
	GetBalance(ctx context.Context, userID int) (*models.Balance, error)
	Withdraw(ctx context.Context, userID int, order string, sum float64) error
	GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)
}
