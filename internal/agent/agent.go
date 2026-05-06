package agent

import (
	"context"
	"errors"
	"fleetglance/internal/agent/providers"
	"fleetglance/internal/agent/routers"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Params struct {
	Port  int
	Debug bool
}

type Agent struct {
	params            Params
	telemetryProvider providers.TelemetryProvider
	server            *http.Server
}

func NewAgent(params Params) *Agent {
	return &Agent{
		params:            params,
		telemetryProvider: providers.New(),
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

	log.Info().Msgf("Agent listening on: %v", a.server.Addr)
	log.Info().Msg("Starting agent... DONE")

	err := a.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (a *Agent) Stop() error {
	log.Info().Msg("Stopping agent...")

	if a.server == nil {
		log.Info().Msg("Stopping agent... DONE")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Stopping agent... FAILED")
		return err
	}

	log.Info().Msg("Stopping agent... DONE")

	return nil
}
