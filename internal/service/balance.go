package service

import (
	"context"
	"errors"

	"gofermart_/internal/helpers"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
)

type BalanceService interface {
	Withdraw(ctx context.Context, userID int, order string, sum float64) error
	GetBalance(ctx context.Context, userID int) (*models.Balance, error)
	GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)
}

type balanceService struct {
	repo storage.Repository
}

func NewBalanceService(repo storage.Repository) BalanceService {
	return &balanceService{repo: repo}
}

func (s *balanceService) Withdraw(ctx context.Context, userID int, order string, sum float64) error {

	if sum <= 0 {
		return ErrInvalidSum
	}

	if !helpers.ValidLuhn(order) {
		return ErrInvalidOrderNumber
	}

	err := s.repo.Withdraw(ctx, userID, order, sum)

	if errors.Is(err, storage.ErrNotEnoughFunds) {
		return ErrNotEnoughFunds
	}

	return err
}

func (s *balanceService) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return s.repo.GetBalance(ctx, userID)
}

func (s *balanceService) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return s.repo.GetWithdrawals(ctx, userID)
}
