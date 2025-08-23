package v1

import (
	"taskL0/internal/controller/http/middleware"
	"taskL0/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewOrderRoutes(
	r chi.Router,
	orderUsecase usecase.Order,
	validator *validator.Validate,
) {

	apiV1Group := &V1{
		orderUsecase: orderUsecase, validator: validator,
	}

	r.Use(middleware.Logging)

	r.Route("/order", func(r chi.Router) {
		r.Get("/get_info/{order_id}", apiV1Group.orderInfo)
	})

	r.Get("/test/get_cache", apiV1Group.getCache)
}
