package order

import (
	"github.com/go-playground/validator/v10"
	lru "github.com/hashicorp/golang-lru/v2"
	"taskL0/internal/entity/order"
	"taskL0/internal/repository"
	"taskL0/internal/usecase"
)

// UseCase -.
type orderUC struct {
	orderRepo repository.OrderRepo
	cache     *lru.Cache[order.OrderUID, order.OrderInfo]
	v         *validator.Validate
}

// New -.
func New(orderRepo repository.OrderRepo, cache *lru.Cache[order.OrderUID, order.OrderInfo], validator *validator.Validate) usecase.Order {
	return &orderUC{
		orderRepo: orderRepo,
		cache:     cache,
		v:         validator,
	}
}
