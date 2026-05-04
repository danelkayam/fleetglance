package routers

import (
	"fleetglance/internal/agent/providers"
	"fleetglance/internal/agent/routers/handlers"

	"github.com/gin-gonic/gin"
)

type Params struct {
	Debug             bool
	TelemetryProvider *providers.TelemetryProvider
}

func NewRouter(params Params) *gin.Engine {
	statusHandler := handlers.NewStatusHandler()
	telemetryHandler := handlers.NewTelemetryHandler(params.TelemetryProvider)

	if params.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default(
		withStrictRouting(),
		withMethodNotAllowed(),
		withRequestContextFallback(),
	)

	statusHandler.BindRoutes(router)
	telemetryHandler.BindRoutes(router)

	return router
}
