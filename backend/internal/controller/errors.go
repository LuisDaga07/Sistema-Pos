package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pos-saas/restaurant-pos/internal/errors"
	"github.com/pos-saas/restaurant-pos/internal/service"
)

func handleError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	msg := "error interno del servidor"

	if appErr, ok := err.(*service.AppError); ok {
		code = appErr.Code
		if appErr.Message != "" {
			msg = appErr.Message
		}
	} else {
		code = errors.HTTPStatus(err)
		if code != http.StatusInternalServerError {
			msg = err.Error()
		}
	}

	c.JSON(code, gin.H{"error": msg})
}
