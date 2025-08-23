package response

import "taskL0/internal/entity/order"

type Cache struct {
	Orders []order.OrderInfo `json:"orders"`
}
