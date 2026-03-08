package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gofermart_/internal/http/middleware"
	"gofermart_/internal/models"
	"gofermart_/internal/service"

	"github.com/stretchr/testify/require"
)

type mockOrderService struct {
	createFunc func(ctx context.Context, userID int, number string) error
	listFunc   func(ctx context.Context, userID int) ([]models.Order, error)
}

func (m *mockOrderService) CreateOrder(ctx context.Context, userID int, number string) error {
	return m.createFunc(ctx, userID, number)
}

func (m *mockOrderService) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return m.listFunc(ctx, userID)
}

func authCtxR(req *http.Request) *http.Request {
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	return req.WithContext(ctx)
}

func TestOrderHandler_Upload(t *testing.T) {

	tests := []struct {
		name       string
		body       string
		serviceErr error
		expected   int
	}{
		{
			name:     "success",
			body:     "79927398713",
			expected: http.StatusAccepted,
		},
		{
			name:       "invalid order number",
			body:       "123",
			serviceErr: service.ErrInvalidOrderNumber,
			expected:   http.StatusUnprocessableEntity,
		},
		{
			name:       "exists for user",
			body:       "79927398713",
			serviceErr: service.ErrOrderExistsForUser,
			expected:   http.StatusOK,
		},
		{
			name:       "exists for other",
			body:       "79927398713",
			serviceErr: service.ErrOrderExistsForOther,
			expected:   http.StatusConflict,
		},
		{
			name:       "internal error",
			body:       "79927398713",
			serviceErr: errors.New("db error"),
			expected:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mock := &mockOrderService{
				createFunc: func(ctx context.Context, userID int, number string) error {
					return tt.serviceErr
				},
			}

			handler := NewOrderHandler(mock)

			req := httptest.NewRequest(
				http.MethodPost,
				"/api/user/orders",
				bytes.NewBufferString(tt.body),
			)

			req = authCtxR(req)

			rec := httptest.NewRecorder()

			handler.Upload(rec, req)

			require.Equal(t, tt.expected, rec.Code)
		})
	}
}

func TestOrderHandler_Upload_EmptyBody(t *testing.T) {

	mock := &mockOrderService{}

	handler := NewOrderHandler(mock)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/user/orders",
		bytes.NewBufferString(" "),
	)

	req = authCtxR(req)

	rec := httptest.NewRecorder()

	handler.Upload(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOrderHandler_List(t *testing.T) {

	now := time.Now()

	tests := []struct {
		name     string
		orders   []models.Order
		err      error
		expected int
	}{
		{
			name:     "orders exist",
			expected: http.StatusOK,
			orders: []models.Order{
				{
					Number:     "79927398713",
					Status:     models.OrderProcessed,
					Accrual:    100,
					UploadedAt: now,
				},
			},
		},
		{
			name:     "no orders",
			expected: http.StatusNoContent,
			orders:   []models.Order{},
		},
		{
			name:     "internal error",
			expected: http.StatusInternalServerError,
			err:      errors.New("db error"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mock := &mockOrderService{
				listFunc: func(ctx context.Context, userID int) ([]models.Order, error) {
					return tt.orders, tt.err
				},
			}

			handler := NewOrderHandler(mock)

			req := httptest.NewRequest(
				http.MethodGet,
				"/api/user/orders",
				nil,
			)

			req = authCtxR(req)

			rec := httptest.NewRecorder()

			handler.List(rec, req)

			require.Equal(t, tt.expected, rec.Code)
		})
	}
}
