package service

import (
	"fmt"

	"github.com/pos-saas/restaurant-pos/internal/errors"
)

type AppError = errors.AppError

func NewAppError(err error, code int, message string) *AppError {
	return errors.NewAppError(err, code, message)
}

func NewValidationError(field, msg string) *AppError {
	return errors.NewAppError(
		fmt.Errorf("validation: %s", field),
		400,
		fmt.Sprintf("%s: %s", field, msg),
	)
}
