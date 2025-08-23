package v1

import (
	"encoding/json"
	"net/http"

	"taskL0/internal/controller/http/v1/response"
	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/pkg/logger"
)

func (h *V1) getCache(w http.ResponseWriter, r *http.Request) {
	ctx := logger.InfoAddValues(
		r.Context(),
		loggertag.HandlerStartedEvent,
		loggertag.APIMethod, "getCache",
	)

	orders := h.orderUsecase.GetCache()

	resp := response.Cache{
		Orders: orders,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error(
			ctx,
			loggertag.HandlerErrorEvent,
			loggertag.Error, err,
		)
		http.Error(w, "can not encode json", http.StatusInternalServerError)
		return
	}

	logger.Info(
		ctx,
		loggertag.HandlerCompletedEvent,
	)
}
