package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithdraw_InvalidSum(t *testing.T) {
	s := NewBalanceService(&mockRepo{})

	err := s.Withdraw(context.Background(), 1, "79927398713", -10)

	require.ErrorIs(t, err, ErrInvalidSum)
}

func TestWithdraw_InvalidOrderNumber(t *testing.T) {
	s := NewBalanceService(&mockRepo{})

	err := s.Withdraw(context.Background(), 1, "123", 10)

	require.ErrorIs(t, err, ErrInvalidOrderNumber)
}

func TestWithdraw_NotEnoughFunds(t *testing.T) {
	s := NewBalanceService(&mockRepo{})

	// этот номер в mockRepo возвращает ErrNotEnoughFunds
	err := s.Withdraw(context.Background(), 1, "79927398713", 10)

	require.ErrorIs(t, err, ErrNotEnoughFunds)
}

func TestWithdraw_Success(t *testing.T) {
	s := NewBalanceService(&mockRepo{})

	err := s.Withdraw(context.Background(), 1, "12345678903", 10)

	require.NoError(t, err)
}

func TestGetBalance(t *testing.T) {
	s := NewBalanceService(&mockRepo{})

	b, err := s.GetBalance(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, 100.0, b.Current)
	require.Equal(t, 50.0, b.Withdrawn)
}

func TestGetWithdrawals(t *testing.T) {
	s := NewBalanceService(&mockRepo{})

	w, err := s.GetWithdrawals(context.Background(), 1)

	require.NoError(t, err)
	require.Len(t, w, 1)
	require.Equal(t, "12345678903", w[0].Order)
}
