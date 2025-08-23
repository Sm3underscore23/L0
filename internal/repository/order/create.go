package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	customerrors "taskL0/internal/entity/custom_errors"
	"taskL0/internal/entity/order"
)

func (r *orderRepo) Create(ctx context.Context, orderInfo order.OrderInfo) error {
	tx, err := r.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	builder := r.Builder.Insert(TableOrders).
		Columns(
			OrdersOrderUID,
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
		).
		Values(
			orderInfo.OrderUID,
			orderInfo.TrackNumber,
			orderInfo.Entry,
			orderInfo.Locale,
			orderInfo.InternalSignature,
			orderInfo.CustomerID,
			orderInfo.DeliveryService,
			orderInfo.ShardKey,
			orderInfo.SmID,
			orderInfo.DateCreated,
			orderInfo.OofShard,
		).
		Suffix(fmt.Sprintf("ON CONFLICT (%s) DO NOTHING RETURNING 1", OrdersOrderUID))

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	var checker int
	err = tx.QueryRow(ctx, query, args...).Scan(&checker)
	if errors.Is(err, pgx.ErrNoRows) {
		return customerrors.ErrOrderAlreadyExists
	}

	if err != nil {
		return err
	}

	builder = r.Builder.Insert(TableDelivery).
		Columns(
			DeliveryOrderUID,
			DeliveryName,
			DeliveryPhone,
			DeliveryZip,
			DeliveryCity,
			DeliveryAddress,
			DeliveryRegion,
			DeliveryEmail,
		).
		Values(
			orderInfo.OrderUID,
			orderInfo.Delivery.Name,
			orderInfo.Delivery.Phone,
			orderInfo.Delivery.Zip,
			orderInfo.Delivery.City,
			orderInfo.Delivery.Address,
			orderInfo.Delivery.Region,
			orderInfo.Delivery.Email,
		)

	query, args, err = builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	builder = r.Builder.Insert(TablePayments).
		Columns(
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
		Values(
			orderInfo.Payment.Transaction,
			orderInfo.Payment.RequestID,
			orderInfo.Payment.Currency,
			orderInfo.Payment.Provider,
			orderInfo.Payment.Amount,
			orderInfo.Payment.PaymentDT,
			orderInfo.Payment.Bank,
			orderInfo.Payment.DeliveryCost,
			orderInfo.Payment.GoodsTotal,
			orderInfo.Payment.CustomFee,
		)

	query, args, err = builder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	for _, item := range orderInfo.Items {
		builder = r.Builder.Insert(TableItems).
			Columns(
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
			Values(
				item.TrackNumber,
				item.ChrtID,
				item.Price,
				item.Rid,
				item.Name,
				item.Sale,
				item.Size,
				item.TotalPrice,
				item.NmID,
				item.Brand,
				item.Status,
			)

		query, args, err = builder.ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
