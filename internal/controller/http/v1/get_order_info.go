package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"taskL0/internal/controller/http/v1/response"
	ce "taskL0/internal/entity/custom_errors"
	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/internal/entity/order"
	"taskL0/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func (h *V1) orderInfo(w http.ResponseWriter, r *http.Request) {
	ctx := logger.InfoAddValues(
		r.Context(),
		loggertag.HandlerStartedEvent,
		loggertag.APIMethod, "orderInfo",
	)

	reqOrderUID := chi.URLParam(r, "order_id")
	if reqOrderUID == "" || len(reqOrderUID) > 36 {
		ctx = logger.AddValuesToContext(ctx, loggertag.Error, ce.ErrOrderUID)
		response.WriteJSONError(ctx, w, ce.ErrOrderUID)
		return
	}

	orderUID := (order.OrderUID)(reqOrderUID)

	ctx = logger.AddValuesToContext(ctx, loggertag.OrderUID, orderUID)

	start := time.Now()
	orderInfo, err := h.orderUsecase.GetInfo(ctx, orderUID)
	if err != nil {
		ctx = logger.AddValuesToContext(ctx, loggertag.Error, err)
		response.WriteJSONError(ctx, w, err)
		return
	}
	elapsed := time.Since(start)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orderInfo); err != nil {
		logger.Error(
			ctx,
			loggertag.HandlerCompletedEvent,
			loggertag.Error, err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	logger.Info(
		ctx,
		loggertag.HandlerCompletedEvent,
		slog.Duration(loggertag.Time, elapsed),
	)

}
