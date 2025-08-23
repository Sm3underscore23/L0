package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
	cu "taskL0/internal/entity/custom_errors"
)

const (
	testOrderUID    = "test_order_uid123"
	testTrackNumber = "track_number123"
)

var validOrder = OrderInfo{
	OrderUID:    testOrderUID,
	TrackNumber: testTrackNumber,
	Payment: Payment{
		Transaction:  testOrderUID,
		Amount:       5,
		DeliveryCost: 1,
		GoodsTotal:   3,
		CustomFee:    1,
	},
	Items: []Item{
		{
			TrackNumber: testTrackNumber,
			Price:       2,
			Sale:        1,
			TotalPrice:  1,
		},
		{
			TrackNumber: testTrackNumber,
			Price:       2,
			Sale:        1,
			TotalPrice:  1,
		},
		{
			TrackNumber: testTrackNumber,
			Price:       2,
			Sale:        1,
			TotalPrice:  1,
		},
	},
}

func TestValidateOrederPrice(t *testing.T) {
	validGoodsTotal := 3

	testTable := []struct {
		name     string
		testData OrderInfo
		noError  bool
	}{
		{
			name:     "OK",
			testData: validOrder,
			noError:  true,
		},
		{
			name: "amount != custom_fee + delivery_cost + goods_total",
			testData: OrderInfo{
				Payment: Payment{
					Amount:       3,
					CustomFee:    1,
					DeliveryCost: 1,
					GoodsTotal:   validGoodsTotal,
				},
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testData.validateOrederPrice(validGoodsTotal)

			if tc.noError {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
		})
	}
}

func TestValidateItemPrice(t *testing.T) {
	testTable := []struct {
		name     string
		testData Item
		noError  bool
	}{
		{
			name:     "OK",
			testData: validOrder.Items[0],
			noError:  true,
		},
		{
			name: "invalid total_price: total_price != price-sale",
			testData: Item{
				Price:      2,
				Sale:       1,
				TotalPrice: 0,
			},
			noError: false,
		},
		{
			name: "invalid total_price: 0 > total_price",
			testData: Item{
				Price:      2,
				Sale:       5,
				TotalPrice: -3,
			},
			noError: false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testData.validateItemPrice()

			if tc.noError {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
		})
	}
}

func TestValidateUID(t *testing.T) {
	invalidOrder := OrderInfo{
		OrderUID: testOrderUID,
		Payment: Payment{
			Transaction: "different",
		},
	}

	testTable := []struct {
		name          string
		testData      OrderInfo
		expectedError error
	}{
		{
			name:          "OK",
			testData:      validOrder,
			expectedError: nil,
		},
		{
			name:          "different",
			testData:      invalidOrder,
			expectedError: cu.ErrDifferentUID,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testData.validateUID()

			assert.ErrorIs(t, tc.expectedError, err)
		})
	}
}

func TestValidateTrackNumber(t *testing.T) {
	testItem := Item{
		TrackNumber: "different",
	}

	testTable := []struct {
		name          string
		testData      Item
		expectedError error
	}{
		{
			name:          "OK",
			testData:      validOrder.Items[0],
			expectedError: nil,
		},
		{
			name:          "different",
			testData:      testItem,
			expectedError: cu.ErrDifferentTrackNum,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.testData.validateTrackNumber(validOrder.TrackNumber)

			assert.ErrorIs(t, tc.expectedError, err)
		})
	}
}
