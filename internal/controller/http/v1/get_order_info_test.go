package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	customerrors "taskL0/internal/entity/custom_errors"
	"taskL0/internal/entity/order"
	"taskL0/internal/usecase"
	mock_usecase "taskL0/internal/usecase/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderInfo(t *testing.T) {
	orderUId := order.OrderUID("123qwerty")

	validOrder := order.OrderInfo{
		OrderUID: orderUId,
	}

	validOrderJSON, err := json.Marshal(validOrder)
	if err != nil {
		t.Fatal(err)
	}

	// invalidOrderJSON := []byte("{invalid}")

	testTable := []struct {
		name           string
		mockOrderUC    func(ctrl *gomock.Controller) usecase.Order
		testData       order.OrderUID
		expectedData   []byte
		expectedStatus int
	}{
		{
			name: "OK",
			mockOrderUC: func(ctrl *gomock.Controller) usecase.Order {
				mock := mock_usecase.NewMockOrder(ctrl)
				mock.EXPECT().GetInfo(gomock.Any(), orderUId).Return(
					validOrder,
					nil,
				)
				return mock
			},
			testData:       orderUId,
			expectedData:   validOrderJSON,
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty order_uid",
			mockOrderUC: func(ctrl *gomock.Controller) usecase.Order {
				return mock_usecase.NewMockOrder(ctrl)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid order_uid",
			mockOrderUC: func(ctrl *gomock.Controller) usecase.Order {
				return mock_usecase.NewMockOrder(ctrl)
			},
			testData:       order.OrderUID("qwertyuiopasdfghjklzxcvbnm12345678910"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "order not exists",
			mockOrderUC: func(ctrl *gomock.Controller) usecase.Order {
				mock := mock_usecase.NewMockOrder(ctrl)
				mock.EXPECT().GetInfo(gomock.Any(), orderUId).Return(
					order.OrderInfo{},
					customerrors.ErrOrderNotExists,
				)
				return mock
			},
			testData:       orderUId,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "db error",
			mockOrderUC: func(ctrl *gomock.Controller) usecase.Order {
				mock := mock_usecase.NewMockOrder(ctrl)
				mock.EXPECT().GetInfo(gomock.Any(), orderUId).Return(
					order.OrderInfo{},
					errors.New("db error"),
				)
				return mock
			},
			testData:       orderUId,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			h := &V1{
				orderUsecase: tc.mockOrderUC(ctrl),
			}

			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/get_info/%s", orderUId), nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("order_id", string(tc.testData))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			h.orderInfo(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if resp.StatusCode != http.StatusOK {
				return
			}

			body, _ := io.ReadAll(resp.Body)

			assert.JSONEq(t, string(tc.expectedData), string(body))
		})
	}
}
