package request

type OrderUID struct {
	OrderUID int `json:"order_id" validate:"required"`
}
