package order

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	cu "taskL0/internal/entity/custom_errors"
	"taskL0/internal/entity/order"
	"taskL0/internal/repository"
	mock_repository "taskL0/internal/repository/mocks"
)

func TestHandleOrder(t *testing.T) {
	err := os.Chdir("../../../")
	if err != nil {
		t.Fatal(err)
	}
	rowValidOrderPath := "testdata/order_test/valid_order.json"
	rowInvalidOrderPath := "testdata/order_test/invalid_order.json"

	rowValidOrderInfo, err := os.ReadFile(rowValidOrderPath)
	if err != nil {
		t.Fatal(err)
	}
	rowInvalidOrderInfo, err := os.ReadFile(rowInvalidOrderPath)
	if err != nil {
		t.Fatal(err)
	}

	var testOrder order.OrderInfo

	err = json.Unmarshal(rowValidOrderInfo, &testOrder)
	if err != nil {
		t.Fatal(err)
	}

	cache, err := lru.New[order.OrderUID, order.OrderInfo](1)
	if err != nil {
		t.Fatal(err)
	}

	validator := validator.New(validator.WithRequiredStructEnabled())

	testTable := []struct {
		name          string
		mockOrderRepo func(*gomock.Controller) repository.OrderRepo
		testData      []byte
		isError       bool
	}{
		{
			name: "OK",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				mock := mock_repository.NewMockOrderRepo(ctrl)
				mock.EXPECT().Create(gomock.Any(), testOrder).Return(nil)
				return mock
			},
			testData: rowValidOrderInfo,
		},
		{
			name: "validation error - by validator",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				return mock_repository.NewMockOrderRepo(ctrl)
			},
			isError:  true,
			testData: []byte("invalid_data"),
		},
		{
			name: "validation error - by custom func",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				return mock_repository.NewMockOrderRepo(ctrl)
			},
			isError:  true,
			testData: rowInvalidOrderInfo,
		},
		{
			name: "order already exists",
			mockOrderRepo: func(ctrl *gomock.Controller) repository.OrderRepo {
				mock := mock_repository.NewMockOrderRepo(ctrl)
				mock.EXPECT().Create(gomock.Any(), testOrder).Return(cu.ErrOrderAlreadyExists)
				return mock
			},
			isError:  true,
			testData: rowValidOrderInfo,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			uc := New(tc.mockOrderRepo(ctrl), cache, validator)

			err := uc.HandleOrder(ctx, tc.testData)

			if tc.isError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
