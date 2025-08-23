package order

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/internal/entity/order"
)

func (uc *orderUC) HandleOrder(ctx context.Context, rowOrderInfo []byte) error {
	var orderInfo order.OrderInfo
	err := json.Unmarshal(rowOrderInfo, &orderInfo)
	if err != nil {
		return err
	}

	if err := uc.v.Struct(&orderInfo); err != nil {
		return err
	}

	if err := orderInfo.ValidateOrder(); err != nil {
		return fmt.Errorf("validation error order %s: %w", orderInfo.OrderUID, err)
	}

	if err := uc.orderRepo.Create(ctx, orderInfo); err != nil {
		return err
	}

	uc.cache.Add(orderInfo.OrderUID, orderInfo)

	slog.InfoContext(
		ctx,
		loggertag.OrderCrtdSccEvent,
		loggertag.OrderUID, orderInfo.OrderUID,
	)

	return nil
}
