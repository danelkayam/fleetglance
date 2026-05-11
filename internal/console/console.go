package console

import (
	"sync"

	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"fleetglance/internal/console/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
)

type Console struct {
	fleet          *config.Fleet
	engine         engine.Engine
	program        *tea.Program
	programOptions []tea.ProgramOption

	stopOnce sync.Once
}

// NewConsole creates a single-use console runtime.
//
// A Console is expected to run through one Start/Stop lifecycle. Stop is safe
// to call more than once, but restarting requires creating a new Console.
func NewConsole(fleet *config.Fleet) *Console {
	return &Console{
		fleet:  fleet,
		engine: engine.New(fleet),
		programOptions: []tea.ProgramOption{
			tea.WithAltScreen(),
			tea.WithoutSignalHandler(),
		},
	}
}

// Start starts the engine and UI, and blocks until the UI exits or Stop is
// called.
func (c *Console) Start() error {
	log.Info().Msg("Starting console...")

	events, err := c.engine.Start()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start engine")
		return err
	}

	model := ui.NewModel(c.fleet)
	c.program = tea.NewProgram(model, c.programOptions...)
	runErr := make(chan error, 1)

	go func() {
		_, err := c.program.Run()
		runErr <- err
	}()

	go forwardEvents(c.program, events)

	log.Info().Msg("Starting console... DONE")

	err = <-runErr

	return err
}

// Stop stops the engine and unblocks Start. It is safe to call more than once.
func (c *Console) Stop() error {
	var err error

	c.stopOnce.Do(func() {
		log.Info().Msg("Stopping console...")

		err = c.engine.Stop()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to stop engine")
		}

		c.program.Quit()

		log.Info().Msg("Stopping console... DONE")
	})

	return err
}

func forwardEvents(program *tea.Program, events <-chan engine.Event) {
	for event := range events {
		program.Send(ui.TelemetryEventMsg(event))
	}
}
