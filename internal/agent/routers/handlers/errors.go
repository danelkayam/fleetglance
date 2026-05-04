package handlers

import (
	"errors"
	"fleetglance/internal/protocol"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AbortWithError(c *gin.Context, err error) {
	status, message := resolveStatusAndMessage(err)
	c.AbortWithStatusJSON(status, protocol.Response{
		Data:  nil,
		Error: &protocol.ResponseError{Message: message},
	})
}

func resolveStatusAndMessage(err error) (int, string) {
	switch {
	case errors.Is(err, protocol.ErrItemNotFound):
		return http.StatusNotFound,
			fmt.Sprintf("Item not found: %s", err.Error())

	case errors.Is(err, protocol.ErrItemConflict):
		return http.StatusConflict,
			fmt.Sprintf("Item conflict: %s", err.Error())

	case errors.Is(err, protocol.ErrInvalidArgs):
		return http.StatusBadRequest,
			fmt.Sprintf("Invalid request: %s", err.Error())

	case errors.Is(err, protocol.ErrUnauthenticated):
		return http.StatusUnauthorized,
			fmt.Sprintf("Unauthorized: %s", err.Error())

	case errors.Is(err, protocol.ErrUnauthorized):
		return http.StatusForbidden,
			fmt.Sprintf("Forbidden: %s", err.Error())

	default:
		return http.StatusInternalServerError,
			fmt.Sprintf("Internal server error: %s", err.Error())
	}
}
