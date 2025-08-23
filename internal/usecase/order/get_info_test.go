package order

import (
	"context"
	"errors"
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"taskL0/internal/entity/order"
	"taskL0/internal/repository"
	mock_repository "taskL0/internal/repository/mocks"
)

func TestGetInfo(t *testing.T) {
	testOrderUID := order.OrderUID("test_order_uid")

	testOrderInfo := order.OrderInfo{
		OrderUID: testOrderUID,
	}

	cache, err := lru.New[order.OrderUID, order.OrderInfo](1)
	if err != nil {
		t.Fatal(err)
	}

	testTable := []struct {
		name          string
		mockOrderRepo func(ctrl *gomock.Controller) repository.OrderRepo
		withCache     bool
		isError       bool
		expectedData  order.OrderInfo
	}{
		{
			name: "OK. Order from chache",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				return mock_repository.NewMockOrderRepo(ctrl)
			},
			withCache:    true,
			isError:      false,
			expectedData: testOrderInfo,
		},
		{
			name: "OK. Order from db",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				mock := mock_repository.NewMockOrderRepo(ctrl)
				mock.EXPECT().IsExists(gomock.Any(), testOrderUID).Return(true, nil)
				mock.EXPECT().GetInfo(gomock.Any(), testOrderUID).Return(testOrderInfo, nil)
				return mock
			},
			isError:      false,
			expectedData: testOrderInfo,
		},
		{
			name: "order not exists",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				mock := mock_repository.NewMockOrderRepo(ctrl)
				mock.EXPECT().IsExists(gomock.Any(), testOrderUID).Return(false, nil)
				return mock
			},
			isError: true,
		},
		{
			name: "db error isExists func",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				mock := mock_repository.NewMockOrderRepo(ctrl)
				mock.EXPECT().IsExists(gomock.Any(), testOrderUID).Return(false, errors.New("db err"))
				return mock
			},
			isError: true,
		},
		{
			name: "db error GetInfo func",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				mock := mock_repository.NewMockOrderRepo(ctrl)
				mock.EXPECT().IsExists(gomock.Any(), testOrderUID).Return(true, nil)
				mock.EXPECT().GetInfo(gomock.Any(), testOrderUID).Return(order.OrderInfo{}, errors.New("db err"))
				return mock
			},
			isError: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			if tc.withCache {
				cache.Add(testOrderUID, testOrderInfo)
			}

			uc := New(tc.mockOrderRepo(ctrl), cache, nil)

			orderInfo, err := uc.GetInfo(ctx, testOrderUID)

			if tc.isError {
				assert.Error(t, err)
			} else {
				orderInfoFromCache, ok := cache.Get(testOrderUID)

				assert.True(t, ok)
				assert.Equal(t, tc.expectedData, orderInfoFromCache)
			}

			assert.Equal(t, tc.expectedData, orderInfo)

			cache.Purge()
		})
	}
}
