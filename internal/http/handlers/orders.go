package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"gofermart_/internal/http/middleware"
	"gofermart_/internal/logger"
	"gofermart_/internal/models"
	"gofermart_/internal/service"
	"gofermart_/internal/storage"
)

// OrderHandler обрабатывает запросы, связанные с заказами пользователей.
type OrderHandler struct {
	Repo storage.Repository
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
	uidValue := r.Context().Value(middleware.UserIDKey)
	if uidValue == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, ok := uidValue.(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	number := strings.TrimSpace(string(body))
	if number == "" {
		http.Error(w, "empty order number", http.StatusBadRequest)
		return
	}

	if !service.ValidLuhn(number) {
		http.Error(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	order := &models.Order{
		Number: number,
		UserID: uid,
		Status: models.OrderNew,
	}

	err = h.Repo.CreateOrder(r.Context(), order)

	switch err {
	case nil:
		w.WriteHeader(http.StatusAccepted)
		return

	case storage.ErrOrderExistsForUser:
		w.WriteHeader(http.StatusOK)
		return

	case storage.ErrOrderExistsForOther:
		http.Error(w, "order already uploaded by another user", http.StatusConflict)
		return

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
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
	uidValue := r.Context().Value(middleware.UserIDKey)
	if uidValue == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, ok := uidValue.(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.Repo.GetUserOrders(r.Context(), uid)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
