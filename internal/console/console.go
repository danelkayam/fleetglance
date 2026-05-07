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

	mu       sync.Mutex
	stopOnce sync.Once
	stopErr  error
}

// NewConsole creates a single-use console runtime.
//
// A Console is expected to run through one Start/Stop lifecycle. Stop is safe
// to call more than once, but restarting requires creating a new Console.
func NewConsole(fleet *config.Fleet) *Console {
	return &Console{
		fleet: fleet,
		programOptions: []tea.ProgramOption{
			tea.WithAltScreen(),
			tea.WithoutSignalHandler(),
		},
	}
}

// Start validates the fleet config, starts the engine and UI, and blocks until
// the UI exits or Stop is called.
func (c *Console) Start() error {
	log.Info().Msg("Starting console...")

	if err := config.ValidateFleet(c.fleet); err != nil {
		return err
	}

	c.mu.Lock()
	if c.engine == nil {
		c.engine = engine.New(c.fleet)
	}
	currentEngine := c.engine
	c.mu.Unlock()

	events, err := currentEngine.Start()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start engine")
		return err
	}

	model := ui.NewModel(c.fleet)
	program := tea.NewProgram(model, c.programOptions...)
	runErr := make(chan error, 1)

	c.mu.Lock()
	c.program = program
	c.mu.Unlock()

	go func() {
		_, err := program.Run()
		runErr <- err
	}()

	go forwardEvents(program, events)

	log.Info().Msg("Starting console... DONE")

	err = <-runErr

	c.mu.Lock()
	if c.program == program {
		c.program = nil
	}
	c.mu.Unlock()

	return err
}

// Stop stops the engine and unblocks Start. It is safe to call more than once.
func (c *Console) Stop() error {
	c.stopOnce.Do(func() {
		log.Info().Msg("Stopping console...")

		c.mu.Lock()
		currentEngine := c.engine
		program := c.program
		c.mu.Unlock()

		if currentEngine != nil {
			if err := currentEngine.Stop(); err != nil {
				log.Warn().Err(err).Msg("Failed to stop engine")
				c.stopErr = err
			}
		}

		if program != nil {
			program.Quit()
		}

		log.Info().Msg("Stopping console... DONE")
	})

	return c.stopErr
}

func forwardEvents(program *tea.Program, events <-chan engine.Event) {
	for event := range events {
		program.Send(ui.TelemetryEventMsg(event))
	}
}
