// Package repository implements application outer layer logic. Each logic group in own file.
package repository

import (
	"context"

	"taskL0/internal/entity/order"
)

//go:generate mockgen -source=contracts.go -destination=mocks/mock.go
type (
	OrderRepo interface {
		Create(ctx context.Context, orderInfo order.OrderInfo) error
		IsExists(ctx context.Context, orderUID order.OrderUID) (bool, error)
		GetInfo(ctx context.Context, orderUID order.OrderUID) (order.OrderInfo, error)
		GetLastsOrders(ctx context.Context, limitRecover int) ([]order.OrderInfo, error)
	}
)
