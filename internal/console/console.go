package console

import (
	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"sync"

	"github.com/rs/zerolog/log"
)

type Console struct {
	fleet  *config.Fleet
	engine engine.Engine

	done chan struct{}
	once sync.Once
}

// NewConsole creates a single-use console runtime.
//
// A Console is expected to run through one Start/Stop lifecycle. Stop is safe
// to call more than once, but restarting requires creating a new Console.
func NewConsole(fleet *config.Fleet) *Console {
	return &Console{
		fleet:  fleet,
		engine: engine.New(fleet),
		done:   make(chan struct{}),
	}
}

// Start validates the fleet config, logs the startup message, and blocks until
// Stop is called.
func (c *Console) Start() error {
	log.Info().Msg("Starting console...")

	if err := config.ValidateFleet(c.fleet); err != nil {
		return err
	}

	events, err := c.engine.Start()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start engine")
		return err
	}

	log.Info().Msg("Starting console... DONE")

	for event := range events {
		if event.Telemetry == nil {
			log.Warn().
				Str("ship", event.ShipName).
				Msg("Received event with no telemetry")

		} else {
			log.Info().
				Str("ship", event.ShipName).
				Float64("CPU", event.Telemetry.CPU.UsagePercent).
				Float64("RAM", event.Telemetry.Memory.UsagePercent).
				Float64("Disk", event.Telemetry.Storage.UsagePercent).
				Msg("Received telemetry event")
		}
	}

	<-c.done

	return nil
}

// Stop unblocks Start. It is safe to call more than once.
func (c *Console) Stop() error {
	log.Info().Msg("Stopping console...")

	err := c.engine.Stop()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to stop engine")
	}

	c.once.Do(func() {
		close(c.done)
	})

	log.Info().Msg("Stopping console... DONE")

	return err
}
