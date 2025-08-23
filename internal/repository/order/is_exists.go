package order

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"taskL0/internal/entity/order"
)

func (r *orderRepo) IsExists(ctx context.Context, orderUID order.OrderUID) (bool, error) {
	builder := r.Builder.Select("1").
		From(TableOrders).
		Where(squirrel.Eq{OrdersOrderUID: orderUID})

	query, args, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	var checker int
	err = r.Pool.QueryRow(ctx, query, args...).Scan(&checker)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, err
}
