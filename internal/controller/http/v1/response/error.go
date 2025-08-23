package response

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	ce "taskL0/internal/entity/custom_errors"
	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/pkg/logger"
)

type errorResponse struct {
	Message string `json:"error" example:"message"`
}

func WriteJSONError(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	response := errorResponse{
		Message: err.Error(),
	}
	statusCode := ce.GetStatusCode(err)
	w.WriteHeader(statusCode)

	logger.Error(ctx,
		loggertag.HandlerErrorEvent,
		loggertag.StatusCode, statusCode,
	)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to write JSONE: %s", err)
	}
}
