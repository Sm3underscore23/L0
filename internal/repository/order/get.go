package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"taskL0/internal/entity/order"
)

func (r *orderRepo) GetInfo(ctx context.Context, orderUID order.OrderUID) (order.OrderInfo, error) {
	var orderInfo order.OrderInfo

	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return orderInfo, err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	builder := r.Builder.Select(
		"o."+OrdersOrderUID+" AS order_uid",
		OrdersTrackNumber,
		OrdersEntry,
		OrdersLocale,
		OrdersInternalSig,
		OrdersCustomerID,
		OrdersDeliveryService,
		OrdersShardKey,
		OrdersSmID,
		OrdersDateCreated,
		OrdersOOFShard,

		DeliveryName,
		DeliveryPhone,
		DeliveryZip,
		DeliveryCity,
		DeliveryAddress,
		DeliveryRegion,
		DeliveryEmail,

		PaymentsTransaction,
		PaymentsRequestID,
		PaymentsCurrency,
		PaymentsProvider,
		PaymentsAmount,
		PaymentsPaymentDt,
		PaymentsBank,
		PaymentsDeliveryCost,
		PaymentsGoodsTotal,
		PaymentsCustomFee,
	).
		From(TableOrders + " o").
		LeftJoin(TableDelivery + " d ON o." + OrdersOrderUID + " = d." + DeliveryOrderUID).
		LeftJoin(TablePayments + " p ON o." + OrdersOrderUID + " = p." + PaymentsTransaction).
		Where(sq.Eq{"o." + OrdersOrderUID: orderUID})

	query, args, err := builder.ToSql()
	if err != nil {
		return orderInfo, err
	}

	err = pgxscan.Get(ctx, tx, &orderInfo, query, args...)
	if err != nil {
		return orderInfo, err
	}

	builder = r.Builder.Select(
		ItemsTrackNumber,
		ItemsChrtID,
		ItemsPrice,
		ItemsRID,
		ItemsName,
		ItemsSale,
		ItemsSize,
		ItemsTotalPrice,
		ItemsNmID,
		ItemsBrand,
		ItemsStatus,
	).
		From(TableItems).
		Where(sq.Eq{ItemsTrackNumber: orderInfo.TrackNumber})

	query, args, err = builder.ToSql()
	if err != nil {
		return orderInfo, err
	}

	var items []order.Item

	err = pgxscan.Select(ctx, tx, &items, query, args...)
	if err != nil {
		return orderInfo, err
	}

	if err := tx.Commit(ctx); err != nil {
		return orderInfo, err
	}

	orderInfo.Items = items

	return orderInfo, nil
}
