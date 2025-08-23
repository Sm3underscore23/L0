// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"taskL0/internal/entity/order"
)

//go:generate mockgen -source=contracts.go -destination=mocks/mock.go
type (
	Order interface {
		HandleOrder(ctx context.Context, value []byte) error
		GetInfo(ctx context.Context, orderUID order.OrderUID) (order.OrderInfo, error)
		GetCache() []order.OrderInfo
	}

	CacheLoader interface {
		WarmUp(ctx context.Context) error
	}
)
