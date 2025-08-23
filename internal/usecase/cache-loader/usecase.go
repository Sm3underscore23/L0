package cacheloader

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"taskL0/internal/entity/order"
	"taskL0/internal/repository"
	"taskL0/internal/usecase"
)

// UseCase -.
type cacheLoader struct {
	orderRepo    repository.OrderRepo
	cache        *lru.Cache[order.OrderUID, order.OrderInfo]
	limitRecover int
}

// New -.
func New(orderRepo repository.OrderRepo, cache *lru.Cache[order.OrderUID, order.OrderInfo], limitRecover int) usecase.CacheLoader {
	return &cacheLoader{orderRepo: orderRepo, cache: cache, limitRecover: limitRecover}
}
