package order

import (
	"taskL0/internal/entity/order"
)

func (uc *orderUC) GetCache() []order.OrderInfo {
	keys := uc.cache.Keys()

	orders := make([]order.OrderInfo, 0, len(keys))

	for _, k := range keys {
		if order, ok := uc.cache.Get(k); ok {
			orders = append(orders, order)
		}
	}

	return orders
}
