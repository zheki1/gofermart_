package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"gofermart_/internal/helpers"
	"gofermart_/internal/logger"
	"gofermart_/internal/models"
	"gofermart_/internal/service"
)

// OrderHandler обрабатывает запросы, связанные с заказами пользователей.
type OrderHandler struct {
	Service service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{Service: s}
}

// Upload загружает новый заказ для авторизованного пользователя.
// Тело запроса должно содержать номер заказа (plain text).
// Возможные ответы:
// - 202 Accepted — заказ успешно принят
// - 200 OK — заказ уже загружен текущим пользователем
// - 409 Conflict — заказ загружен другим пользователем
// - 400 Bad Request — некорректное тело запроса или пустой номер заказа
// - 401 Unauthorized — пользователь не авторизован
// - 422 Unprocessable Entity — неверный номер заказа (Luhn не проходит)
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *OrderHandler) Upload(w http.ResponseWriter, r *http.Request) {
	uid, err := helpers.GetUserID(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		helpers.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	number := strings.TrimSpace(string(body))
	if number == "" {
		helpers.WriteJSONError(w, "empty order number", http.StatusBadRequest)
		return
	}

	err = h.Service.CreateOrder(r.Context(), uid, number)
	switch {
	case err == nil:
		w.WriteHeader(http.StatusAccepted)
	case errors.Is(err, service.ErrInvalidOrderNumber):
		helpers.WriteJSONError(w, err.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, service.ErrOrderExistsForUser):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, service.ErrOrderExistsForOther):
		helpers.WriteJSONError(w, err.Error(), http.StatusConflict)
	default:
		helpers.WriteJSONError(w, "internal server error", http.StatusInternalServerError)
	}
}

// List возвращает список заказов для авторизованного пользователя.
// Ответ формируется в формате JSON: массив объектов OrderResponse.
// Возможные ответы:
// - 200 OK — список заказов возвращён
// - 204 No Content — заказы отсутствуют
// - 401 Unauthorized — пользователь не авторизован
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	uid, err := helpers.GetUserID(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	orders, err := h.Service.GetUserOrders(r.Context(), uid)
	if err != nil {
		helpers.WriteJSONError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]models.OrderResponse, 0, len(orders))

	for _, o := range orders {
		var accrual *float64
		if o.Status == models.OrderProcessed {
			accrual = &o.Accrual
		}

		resp = append(resp, models.OrderResponse{
			Number:     o.Number,
			Status:     string(o.Status),
			Accrual:    accrual,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log.Error(err)
		return
	}
}
