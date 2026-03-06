package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gofermart_/internal/http/middleware"
	"gofermart_/internal/models"
	"gofermart_/internal/service"
)

type mockBalanceService struct {
	withdrawFn       func(ctx context.Context, userID int, order string, sum float64) error
	getBalanceFn     func(ctx context.Context, userID int) (*models.Balance, error)
	getWithdrawalsFn func(ctx context.Context, userID int) ([]models.Withdrawal, error)
}

func (m *mockBalanceService) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	return m.withdrawFn(ctx, userID, order, sum)
}

func (m *mockBalanceService) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	return m.getBalanceFn(ctx, userID)
}

func (m *mockBalanceService) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	return m.getWithdrawalsFn(ctx, userID)
}

func authCtx() context.Context {
	return context.WithValue(context.Background(), middleware.UserIDKey, 1)
}

func TestWithdraw_Unauthorized(t *testing.T) {

	h := NewBalanceHandler(&mockBalanceService{})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", nil)
	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestWithdraw_InvalidJSON(t *testing.T) {

	h := NewBalanceHandler(&mockBalanceService{})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBufferString("bad"))
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWithdraw_InvalidSum(t *testing.T) {

	s := &mockBalanceService{
		withdrawFn: func(ctx context.Context, userID int, order string, sum float64) error {
			return service.ErrInvalidSum
		},
	}

	h := NewBalanceHandler(s)

	body, _ := json.Marshal(map[string]interface{}{
		"order": "79927398713",
		"sum":   -10,
	})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWithdraw_InvalidOrder(t *testing.T) {

	s := &mockBalanceService{
		withdrawFn: func(ctx context.Context, userID int, order string, sum float64) error {
			return service.ErrInvalidOrderNumber
		},
	}

	h := NewBalanceHandler(s)

	body, _ := json.Marshal(map[string]interface{}{
		"order": "123",
		"sum":   10,
	})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestWithdraw_NotEnoughFunds(t *testing.T) {

	s := &mockBalanceService{
		withdrawFn: func(ctx context.Context, userID int, order string, sum float64) error {
			return service.ErrNotEnoughFunds
		},
	}

	h := NewBalanceHandler(s)

	body, _ := json.Marshal(map[string]interface{}{
		"order": "79927398713",
		"sum":   10,
	})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	require.Equal(t, http.StatusPaymentRequired, w.Code)
}

func TestWithdraw_Success(t *testing.T) {

	s := &mockBalanceService{
		withdrawFn: func(ctx context.Context, userID int, order string, sum float64) error {
			return nil
		},
	}

	h := NewBalanceHandler(s)

	body, _ := json.Marshal(map[string]interface{}{
		"order": "79927398713",
		"sum":   10,
	})

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdraw(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetBalance_Success(t *testing.T) {

	s := &mockBalanceService{
		getBalanceFn: func(ctx context.Context, userID int) (*models.Balance, error) {
			return &models.Balance{Current: 100, Withdrawn: 50}, nil
		},
	}

	h := NewBalanceHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Get(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestWithdrawals_NoContent(t *testing.T) {

	s := &mockBalanceService{
		getWithdrawalsFn: func(ctx context.Context, userID int) ([]models.Withdrawal, error) {
			return []models.Withdrawal{}, nil
		},
	}

	h := NewBalanceHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdrawals(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestWithdrawals_Success(t *testing.T) {

	s := &mockBalanceService{
		getWithdrawalsFn: func(ctx context.Context, userID int) ([]models.Withdrawal, error) {
			return []models.Withdrawal{
				{
					Order:       "12345678903",
					Sum:         10,
					ProcessedAt: time.Now(),
				},
			}, nil
		},
	}

	h := NewBalanceHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/withdrawals", nil)
	req = req.WithContext(authCtx())

	w := httptest.NewRecorder()

	h.Withdrawals(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
