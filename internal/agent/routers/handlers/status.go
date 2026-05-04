package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct{}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

func (h *StatusHandler) BindRoutes(router *gin.Engine) {
	router.GET("/health", h.Health)
}

func (h *StatusHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
