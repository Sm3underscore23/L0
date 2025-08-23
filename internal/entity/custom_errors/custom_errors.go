package customerrors

import (
	"errors"
	"net/http"
)

var (
	ErrHtmlPageNotFound = errors.New("html page not found")

	ErrOrderUID     = errors.New("empty or invalid order_uid")
	ErrInvalidInput = errors.New("invalid input data")

	ErrOrderNotExists     = errors.New("order not exists")
	ErrOrderAlreadyExists = errors.New("order already existed")

	ErrDbPasswordEpt = errors.New("db password is empty")

	ErrDifferentUID      = errors.New("order_uid and trx is different")
	ErrDifferentTrackNum = errors.New("item track_namber and order track_namber is different")

	errWithStatus = map[error]int{
		ErrOrderUID:       http.StatusBadRequest,
		ErrInvalidInput:   http.StatusBadRequest,
		ErrOrderNotExists: http.StatusBadRequest,
	}
)

func GetStatusCode(err error) int {
	for mapError, code := range errWithStatus {
		if errors.Is(err, mapError) {
			return code
		}
	}
	return http.StatusInternalServerError
}
