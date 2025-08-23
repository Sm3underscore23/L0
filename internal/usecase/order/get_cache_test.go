package order

import (
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stretchr/testify/assert"
	"taskL0/internal/entity/order"
)

func TestGetCache(t *testing.T) {
	const (
		testOrderUID1 = "test_order_uid_1"
		testOrderUID2 = "test_order_uid_2"
		testOrderUID3 = "test_order_uid_3"
	)

	orders := []order.OrderInfo{
		{
			OrderUID: testOrderUID1,
		},
		{
			OrderUID: testOrderUID2,
		},
		{
			OrderUID: testOrderUID3,
		},
	}

	cacheCap := 3

	generateTestCache := func(orders []order.OrderInfo) *lru.Cache[order.OrderUID, order.OrderInfo] {
		cache, err := lru.New[order.OrderUID, order.OrderInfo](cacheCap)
		if err != nil {
			t.Fatal(err)
		}
		for _, order := range orders {
			cache.Add(order.OrderUID, order)
		}
		return cache
	}

	testTable := []struct {
		name         string
		isError      bool
		expectedData []order.OrderInfo
	}{
		{
			name:         "not empty cache",
			expectedData: orders,
		},
		{
			name:         "empty cache",
			expectedData: []order.OrderInfo{},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			cache := generateTestCache(tc.expectedData)

			uc := New(nil, cache, nil)

			orders := uc.GetCache()

			assert.Equal(t, tc.expectedData, orders)
			assert.Equal(t, len(tc.expectedData), cache.Len())
		})
	}
}
