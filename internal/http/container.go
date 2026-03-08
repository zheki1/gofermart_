package http

import (
	"gofermart_/internal/service"
	"gofermart_/internal/storage"
)

type Container struct {
	AuthService    service.AuthService
	OrderService   service.OrderService
	BalanceService service.BalanceService
}

func NewContainer(repo storage.Repository) *Container {
	return &Container{
		AuthService:    service.NewAuthService(repo),
		OrderService:   service.NewOrderService(repo),
		BalanceService: service.NewBalanceService(repo),
	}
}
