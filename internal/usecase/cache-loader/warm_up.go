package cacheloader

import "context"

func (uc *cacheLoader) WarmUp(ctx context.Context) error {
	orders, err := uc.orderRepo.GetLastsOrders(ctx, uc.limitRecover)
	if err != nil {
		return err
	}
	for _, order := range orders {
		uc.cache.ContainsOrAdd(order.OrderUID, order)
	}
	return nil
}
