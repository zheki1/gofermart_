package service

import (
	"context"
	"errors"
	"gofermart_/internal/helpers"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID int, number string) error
	GetUserOrders(ctx context.Context, userID int) ([]models.Order, error)
}

type orderService struct {
	repo storage.Repository
}

func NewOrderService(repo storage.Repository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, userID int, number string) error {

	if !helpers.ValidLuhn(number) {
		return ErrInvalidOrderNumber
	}

	order := &models.Order{
		Number: number,
		UserID: userID,
		Status: models.OrderNew,
	}

	err := s.repo.CreateOrder(ctx, order)

	if errors.Is(err, storage.ErrOrderExistsForUser) {
		return ErrOrderExistsForUser
	}

	if errors.Is(err, storage.ErrOrderExistsForOther) {
		return ErrOrderExistsForOther
	}

	return err
}

func (s *orderService) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return s.repo.GetUserOrders(ctx, userID)
}
