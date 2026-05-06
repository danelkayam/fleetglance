package console

import (
	"fleetglance/internal/console/config"
	"sync"

	"github.com/rs/zerolog/log"
)

type Console struct {
	fleet *config.Fleet
	done  chan struct{}
	once  sync.Once
}

// NewConsole creates a single-use console runtime.
//
// A Console is expected to run through one Start/Stop lifecycle. Stop is safe
// to call more than once, but restarting requires creating a new Console.
func NewConsole(fleet *config.Fleet) *Console {
	return &Console{
		fleet: fleet,
		done:  make(chan struct{}),
	}
}

// Start validates the fleet config, logs the startup message, and blocks until
// Stop is called.
func (c *Console) Start() error {
	log.Info().Msg("Starting console...")

	if err := config.ValidateFleet(c.fleet); err != nil {
		return err
	}

	log.Info().Msg("hello console")
	log.Info().Msg("Starting console... DONE")

	<-c.done

	return nil
}

// Stop unblocks Start. It is safe to call more than once.
func (c *Console) Stop() error {
	log.Info().Msg("Stopping console...")

	c.once.Do(func() {
		close(c.done)
	})

	log.Info().Msg("Stopping console... DONE")

	return nil
}
