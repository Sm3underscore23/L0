package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"taskL0/internal/entity/order"
)

func (r *orderRepo) GetLastsOrders(ctx context.Context, limitRecover int) ([]order.OrderInfo, error) {
	var orders []order.OrderInfo

	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
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
		OrderBy(OrdersDateCreated).
		Limit(uint64(limitRecover))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	if err := pgxscan.Select(ctx, tx, &orders, query, args...); err != nil {
		return nil, err
	}

	trackNumbers := make([]string, 0, len(orders))
	for _, o := range orders {
		trackNumbers = append(trackNumbers, o.TrackNumber)
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
		Where(sq.Eq{ItemsTrackNumber: trackNumbers})

	query, args, err = builder.ToSql()
	if err != nil {
		return nil, err
	}

	var items []order.Item
	if err := pgxscan.Select(ctx, tx, &items, query, args...); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	itemMap := make(map[string][]order.Item)
	for _, item := range items {
		itemMap[item.TrackNumber] = append(itemMap[item.TrackNumber], item)
	}

	for i, o := range orders {
		orders[i].Items = itemMap[o.TrackNumber]
	}

	return orders, nil
}
