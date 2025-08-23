package order

import (
	"taskL0/pkg/postgres"
)

// Tables
const (
	TableOrders   = "orders"
	TableDelivery = "delivery"
	TablePayments = "payments"
	TableItems    = "items"
)

// Atributs orders
const (
	OrdersOrderUID        = "order_uid"
	OrdersTrackNumber     = "track_number"
	OrdersEntry           = "entry"
	OrdersLocale          = "locale"
	OrdersInternalSig     = "internal_signature"
	OrdersCustomerID      = "customer_id"
	OrdersDeliveryService = "delivery_service"
	OrdersShardKey        = "shardkey"
	OrdersSmID            = "sm_id"
	OrdersDateCreated     = "date_created"
	OrdersOOFShard        = "oof_shard"
)

// Atributs delivery
const (
	DeliveryOrderUID = "order_uid"
	DeliveryName     = "name"
	DeliveryPhone    = "phone"
	DeliveryZip      = "zip"
	DeliveryCity     = "city"
	DeliveryAddress  = "address"
	DeliveryRegion   = "region"
	DeliveryEmail    = "email"
)

// Atributs payments
const (
	PaymentsTransaction  = "transaction"
	PaymentsRequestID    = "request_id"
	PaymentsCurrency     = "currency"
	PaymentsProvider     = "provider"
	PaymentsAmount       = "amount"
	PaymentsPaymentDt    = "payment_dt"
	PaymentsBank         = "bank"
	PaymentsDeliveryCost = "delivery_cost"
	PaymentsGoodsTotal   = "goods_total"
	PaymentsCustomFee    = "custom_fee"
)

// Atributs items
const (
	ItemsID          = "id"
	ItemsTrackNumber = "track_number"
	ItemsChrtID      = "chrt_id"
	ItemsPrice       = "price"
	ItemsRID         = "rid"
	ItemsName        = "name"
	ItemsSale        = "sale"
	ItemsSize        = "size"
	ItemsTotalPrice  = "total_price"
	ItemsNmID        = "nm_id"
	ItemsBrand       = "brand"
	ItemsStatus      = "status"
)

type orderRepo struct {
	*postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *orderRepo {
	return &orderRepo{pg}
}
