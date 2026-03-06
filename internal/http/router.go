package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"gofermart_/internal/http/handlers"
	"gofermart_/internal/http/middleware"
)

// NewRouter создаёт и возвращает новый HTTP-роутер.
//
// repo — интерфейс Repository для работы с базой данных.
// Возвращает http.Handler с настроенными маршрутами и middleware.
//
// Настроенные маршруты:
// - POST /api/user/register — регистрация пользователя
// - POST /api/user/login — вход пользователя
// - POST /api/user/orders — загрузка заказов (требует авторизации)
// - GET /api/user/orders — список заказов пользователя (требует авторизации)
// - GET /api/user/balance — получение баланса (требует авторизации)
// - POST /api/user/balance/withdraw — вывод средств (требует авторизации)
// - GET /api/user/withdrawals — список выводов средств (требует авторизации)
//
// Дополнительно:
// - GET /chin — тестовый защищённый маршрут, возвращает "chopa authorized".
func NewRouter(c *Container) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logging)

	authH := handlers.NewAuthHandler(c.AuthService)
	r.Post("/api/user/register", authH.Register)
	r.Post("/api/user/login", authH.Login)

	orderH := handlers.NewOrderHandler(c.OrderService)
	balanceH := handlers.NewBalanceHandler(c.BalanceService)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		r.Post("/api/user/orders", orderH.Upload)
		r.Get("/api/user/orders", orderH.List)

		r.Get("/api/user/balance", balanceH.Get)
		r.Post("/api/user/balance/withdraw", balanceH.Withdraw)
		r.Get("/api/user/withdrawals", balanceH.Withdrawals)

		r.Get("/chin", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("chopa authorized")) })
	})

	return r
}
