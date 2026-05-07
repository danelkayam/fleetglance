package console

import (
	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestConsoleStartValidatesFleet(t *testing.T) {
	tests := []struct {
		name    string
		fleet   *config.Fleet
		wantErr string
	}{
		{
			name:    "nil fleet",
			fleet:   nil,
			wantErr: "fleet config is required",
		},
		{
			name: "unsupported version",
			fleet: &config.Fleet{
				Version: 2,
				Ships: map[string]config.Ship{
					"donnager": {URL: "http://donnager:9800"},
				},
			},
			wantErr: "unsupported fleet config version: 2",
		},
		{
			name: "empty ships",
			fleet: &config.Fleet{
				Version: 1,
				Ships:   map[string]config.Ship{},
			},
			wantErr: "fleet must contain at least one ship",
		},
		{
			name: "missing ship url",
			fleet: &config.Fleet{
				Version: 1,
				Ships: map[string]config.Ship{
					"donnager": {},
				},
			},
			wantErr: `ship "donnager" url is required`,
		},
		{
			name: "invalid ship url",
			fleet: &config.Fleet{
				Version: 1,
				Ships: map[string]config.Ship{
					"donnager": {URL: "donnager:9800"},
				},
			},
			wantErr: `ship "donnager" url must be absolute http/https URL`,
		},
		{
			name:    "too many ships",
			fleet:   testFleetWithShips(9),
			wantErr: "fleet supports at most 8 ships",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConsole(tt.fleet)

			err := c.Start()
			if err == nil {
				t.Fatal("expected error")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func testFleetWithShips(count int) *config.Fleet {
	ships := make(map[string]config.Ship, count)
	for i := range count {
		name := "ship" + string(rune('a'+i))
		ships[name] = config.Ship{URL: "http://" + name + ":9800"}
	}

	return &config.Fleet{
		Version: 1,
		Ships:   ships,
	}
}

func TestNewConsoleNilFleetDoesNotPanic(t *testing.T) {
	c := NewConsole(nil)
	if c == nil {
		t.Fatal("console should be created")
	}
	if c.engine != nil {
		t.Fatal("engine should not be constructed before Start")
	}
}

func TestConsoleStopBeforeEngineStartIsSafe(t *testing.T) {
	c := NewConsole(nil)

	if err := c.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
}

func TestConsoleStartStops(t *testing.T) {
	c := NewConsole(&config.Fleet{
		Version: 1,
		Ships: map[string]config.Ship{
			"donnager": {URL: "http://donnager:9800"},
		},
	})
	fakeEngine := &fakeEngine{events: make(chan engine.Event)}
	c.engine = fakeEngine
	c.programOptions = quietProgramOptions()

	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Start()
	}()

	select {
	case err := <-errChan:
		t.Fatalf("start returned before stop: %v", err)
	case <-time.After(10 * time.Millisecond):
	}

	if !fakeEngine.started {
		t.Fatal("engine was not started")
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	if err := c.Stop(); err != nil {
		t.Fatalf("second stop failed: %v", err)
	}

	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("start failed: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("start did not stop")
	}
}

type fakeEngine struct {
	events  chan engine.Event
	started bool
	once    sync.Once
}

func (f *fakeEngine) Start() (<-chan engine.Event, error) {
	f.started = true
	return f.events, nil
}

func (f *fakeEngine) Stop() error {
	f.once.Do(func() {
		close(f.events)
	})
	return nil
}

func quietProgramOptions() []tea.ProgramOption {
	return []tea.ProgramOption{
		tea.WithInput(nil),
		tea.WithOutput(io.Discard),
		tea.WithoutRenderer(),
		tea.WithoutSignalHandler(),
		tea.WithoutSignals(),
	}
}
