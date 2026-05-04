package handlers

import (
	"fleetglance/internal/agent/providers"
	"fleetglance/internal/protocol"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TelemetryHandler struct {
	provider providers.TelemetryProvider
}

func NewTelemetryHandler(provider *providers.TelemetryProvider) *TelemetryHandler {
	return &TelemetryHandler{
		provider: *provider,
	}
}

func (h *TelemetryHandler) BindRoutes(router *gin.Engine) {
	router.GET("/api/telemetry", h.GetTelemetry)
}

func (h *TelemetryHandler) GetTelemetry(c *gin.Context) {
	telemetry, err := h.provider.GetTelemetry()
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Data:  telemetry,
		Error: nil,
	})
}
