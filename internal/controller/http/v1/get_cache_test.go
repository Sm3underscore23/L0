package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"taskL0/internal/controller/http/v1/response"
	"taskL0/internal/entity/order"
	"taskL0/internal/usecase"
	mock_usecase "taskL0/internal/usecase/mocks"
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

	testCache := response.Cache{
		Orders: orders,
	}

	testTable := []struct {
		name           string
		mockOrderUC    func(ctrl *gomock.Controller) usecase.Order
		expectedData   response.Cache
		expectedStatus int
	}{
		{
			name: "OK",
			mockOrderUC: func(ctrl *gomock.Controller) usecase.Order {
				mock := mock_usecase.NewMockOrder(ctrl)
				mock.EXPECT().GetCache().Return(orders)
				return mock
			},
			expectedData:   testCache,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := V1{
				orderUsecase: tc.mockOrderUC(ctrl),
			}

			req := httptest.NewRequest("GET", "/v1/test/get_cache", nil)

			w := httptest.NewRecorder()

			h.getCache(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
