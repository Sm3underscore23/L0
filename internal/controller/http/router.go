// Package v1 implements routing paths. Each services in own file.
package http

import (
	"net/http"

	v1 "taskL0/internal/controller/http/v1"
	"taskL0/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// NewRouter -.
func NewRouter(
	orderUsecase usecase.Order,
	validator *validator.Validate,
) http.Handler {
	r := chi.NewRouter()

	// Routers
	r.Route("/v1", func(r chi.Router) {
		v1.NewOrderRoutes(r, orderUsecase, validator)
	})

	return r
}
