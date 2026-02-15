package errors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound         = errors.New("recurso no encontrado")
	ErrUnauthorized     = errors.New("no autorizado")
	ErrForbidden        = errors.New("acceso denegado")
	ErrBadRequest       = errors.New("solicitud inválida")
	ErrConflict         = errors.New("el recurso ya existe")
	ErrInternal         = errors.New("error interno del servidor")
	ErrInvalidCredentials = errors.New("credenciales inválidas")
)

// AppError representa un error de aplicación con código HTTP
type AppError struct {
	Err     error
	Code    int
	Message string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, code int, message string) *AppError {
	return &AppError{Err: err, Code: code, Message: message}
}

// Is delega a la biblioteca estándar (para evitar conflictos de importación)
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// HTTPStatus devuelve el código HTTP apropiado
func HTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrInvalidCredentials) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden
	}
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrConflict) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}
