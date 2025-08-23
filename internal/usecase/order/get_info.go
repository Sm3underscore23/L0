package order

import (
	"context"

	cu "taskL0/internal/entity/custom_errors"
	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/internal/entity/order"
	"taskL0/pkg/logger"
)

func (uc *orderUC) GetInfo(ctx context.Context, orderUID order.OrderUID) (order.OrderInfo, error) {
	var orderInfo order.OrderInfo

	orderInfo, ok := uc.cache.Get(orderUID)
	if ok {
		logger.Debug(ctx, loggertag.FromCacheEvent)
		return orderInfo, nil
	}

	isExists, err := uc.orderRepo.IsExists(ctx, orderUID)
	if err != nil {
		return orderInfo, err
	}

	if !isExists {
		return orderInfo, cu.ErrOrderNotExists
	}

	orderInfo, err = uc.orderRepo.GetInfo(ctx, orderUID)
	if err != nil {
		return orderInfo, err
	}

	logger.Debug(ctx, loggertag.FromDBEvent)

	uc.cache.ContainsOrAdd(orderUID, orderInfo)

	return orderInfo, nil
}
