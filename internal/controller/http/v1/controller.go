package v1

import (
	"taskL0/internal/usecase"

	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	orderUsecase usecase.Order
	validator    *validator.Validate
}
