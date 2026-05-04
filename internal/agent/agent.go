package agent

import (
	"fleetglance/internal/agent/providers"
	"fleetglance/internal/agent/routers"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Params struct {
	Port  int
	Debug bool
}

type Agent struct {
	params            Params
	telemetryProvider *providers.TelemetryProvider
	server            *http.Server
}

func NewAgent(params Params) *Agent {
	return &Agent{
		params:            params,
		telemetryProvider: providers.NewTelemetryProvider(),
	}
}

func (a *Agent) Start() error {
	log.Info().Msg("Starting agent...")

	router := routers.NewRouter(routers.Params{
		Debug:             a.params.Debug,
		TelemetryProvider: a.telemetryProvider,
	})

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.params.Port),
		Handler: router,
	}

	log.Info().Msg("Starting agent... DONE")

	return a.server.ListenAndServe()
}

func (a *Agent) Stop() error {
	log.Info().Msg("Stopping agent...")
	err := a.server.Close()
	if err != nil {
		log.Error().Err(err).Msg("Stopping agent... FAILED")
		return err
	}

	log.Info().Msg("Stopping agent... DONE")

	return nil
}
