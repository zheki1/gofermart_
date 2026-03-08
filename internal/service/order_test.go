package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateOrder_InvalidNumber(t *testing.T) {
	s := NewOrderService(&mockRepo{})

	err := s.CreateOrder(context.Background(), 1, "123")

	require.ErrorIs(t, err, ErrInvalidOrderNumber)
}

func TestCreateOrder_OrderExistsForUser(t *testing.T) {
	s := NewOrderService(&mockRepo{})

	// номер, который mockRepo обрабатывает как ErrOrderExistsForUser
	err := s.CreateOrder(context.Background(), 1, "4532015112830366")

	require.ErrorIs(t, err, ErrOrderExistsForUser)
}

func TestCreateOrder_OrderExistsForOther(t *testing.T) {
	s := NewOrderService(&mockRepo{})

	// номер, который mockRepo обрабатывает как ErrOrderExistsForOther
	err := s.CreateOrder(context.Background(), 1, "4485275742308327")

	require.ErrorIs(t, err, ErrOrderExistsForOther)
}

func TestCreateOrder_Success(t *testing.T) {
	s := NewOrderService(&mockRepo{})

	err := s.CreateOrder(context.Background(), 1, "79927398713")

	require.NoError(t, err)
}

func TestGetUserOrders_WithOrders(t *testing.T) {
	s := NewOrderService(&mockRepo{})

	orders, err := s.GetUserOrders(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Equal(t, "12345678903", orders[0].Number)
}

func TestGetUserOrders_Empty(t *testing.T) {
	s := NewOrderService(&mockRepo{})

	orders, err := s.GetUserOrders(context.Background(), 999)

	require.NoError(t, err)
	require.Len(t, orders, 0)
}
