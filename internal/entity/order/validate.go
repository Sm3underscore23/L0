package order

import (
	"errors"
	"fmt"
	"log"

	cu "taskL0/internal/entity/custom_errors"
)

func (o *OrderInfo) ValidateOrder() error {
	if err := o.validateUID(); err != nil {
		return err
	}

	var validGoodsTotal int
	for _, item := range o.Items {
		if err := item.validateTrackNumber(o.TrackNumber); err != nil {
			return err
		}
		if err := item.validateItemPrice(); err != nil {
			return fmt.Errorf("invalid item price: %w", err)
		}

		validGoodsTotal += item.TotalPrice
	}

	if err := o.validateOrederPrice(validGoodsTotal); err != nil {
		return fmt.Errorf("invalid order price: %w", err)
	}

	return nil
}

func (i *Item) validateItemPrice() error {
	if i.TotalPrice != i.Price-i.Sale {
		return errors.New("total_price != price-sale")
	}

	if 0 > i.TotalPrice {
		return errors.New("0 > total_price")
	}

	if i.TotalPrice == 0 {
		log.Printf("WARNING: item %v with price 0, track_number %s", i.ChrtID, i.TrackNumber)
	}

	return nil
}

func (o *OrderInfo) validateOrederPrice(validGoodsTotal int) error {
	if validGoodsTotal != o.Payment.GoodsTotal {
		return errors.New("valid_goods_total != goods_total")
	}

	if o.Payment.GoodsTotal == 0 {
		log.Printf("WARNING: order %s total_price is 0", o.OrderUID)
	}

	if o.Payment.Amount != o.Payment.CustomFee+o.Payment.DeliveryCost+o.Payment.GoodsTotal {
		return errors.New("amount != custom_fee + delivery_cost + goods_total")
	}

	return nil
}

func (o *OrderInfo) validateUID() error {
	if o.OrderUID != OrderUID(o.Payment.Transaction) {
		return cu.ErrDifferentUID
	}
	return nil
}

func (i *Item) validateTrackNumber(trackNumber string) error {
	if i.TrackNumber != trackNumber {
		return cu.ErrDifferentTrackNum
	}
	return nil
}
