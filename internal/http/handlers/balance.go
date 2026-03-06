package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gofermart_/internal/helpers"
	"gofermart_/internal/http/middleware"
	"gofermart_/internal/logger"
	"gofermart_/internal/storage"
)

// BalanceHandler обрабатывает запросы, связанные с балансом пользователя.
type BalanceHandler struct {
	Repo storage.Repository
}

// Withdraw обрабатывает запрос на снятие средств пользователем.
// JSON тело запроса должно содержать поля:
// - order (номер заказа, plain text)
// - sum (сумма для списания, >0)
// Возможные ответы:
// - 200 OK — снятие прошло успешно
// - 402 Payment Required — недостаточно средств
// - 400 Bad Request — некорректный JSON или сумма <=0
// - 401 Unauthorized — пользователь не авторизован
// - 422 Unprocessable Entity — неверный номер заказа
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	uidValue := r.Context().Value(middleware.UserIDKey)
	if uidValue == nil {
		helpers.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, ok := uidValue.(int)
	if !ok {
		helpers.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	type request struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteJSONError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Sum <= 0 {
		helpers.WriteJSONError(w, "invalid sum", http.StatusBadRequest)
		return
	}

	if !helpers.ValidLuhn(req.Order) {
		helpers.WriteJSONError(w, "invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err := h.Repo.Withdraw(r.Context(), uid, req.Order, req.Sum)

	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case storage.ErrNotEnoughFunds:
		helpers.WriteJSONError(w, "not enough funds", http.StatusPaymentRequired)
	default:
		helpers.WriteJSONError(w, "internal server error", http.StatusInternalServerError)
	}
}

// Get возвращает текущий баланс пользователя.
// Ответ JSON содержит:
// - current — доступные средства
// - withdrawn — общая сумма снятых средств
// Возможные ответы:
// - 200 OK — баланс успешно получен
// - 401 Unauthorized — пользователь не авторизован
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *BalanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	uidValue := r.Context().Value(middleware.UserIDKey)
	if uidValue == nil {
		helpers.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	uid, ok := uidValue.(int)
	if !ok {
		helpers.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := h.Repo.GetBalance(r.Context(), uid)
	if err != nil {
		helpers.WriteJSONError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Log.Error(err)
		return
	}
}

// Withdrawals возвращает список снятий пользователя.
// JSON массив объектов содержит:
// - order — номер заказа
// - sum — сумма снятия
// - processed_at — время снятия в UTC (RFC3339)
// Возможные ответы:
// - 200 OK — список снятий возвращён
// - 204 No Content — снятий нет
// - 401 Unauthorized — пользователь не авторизован
// - 500 Internal Server Error — внутренняя ошибка сервиса
func (h *BalanceHandler) Withdrawals(w http.ResponseWriter, r *http.Request) {
	uidVal := r.Context().Value(middleware.UserIDKey)
	if uidVal == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	uid, ok := uidVal.(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	list, err := h.Repo.GetWithdrawals(r.Context(), uid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(list) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type response struct {
		Order       string  `json:"order"`
		Sum         float64 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}

	out := make([]response, 0, len(list))
	for _, wth := range list {
		out = append(out, response{
			Order:       wth.Order,
			Sum:         wth.Sum,
			ProcessedAt: wth.ProcessedAt.UTC().Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(out); err != nil {
		logger.Log.Error(err)
		return
	}
}
