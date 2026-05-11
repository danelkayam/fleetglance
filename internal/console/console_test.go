package console

import (
	"fleetglance/internal/console/config"
	"fleetglance/internal/console/engine"
	"io"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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

func TestNewConsoleCreatesEngine(t *testing.T) {
	c := NewConsole(testFleetWithShips(1))
	if c == nil {
		t.Fatal("console should be created")
	}
	if c.engine == nil {
		t.Fatal("engine should be constructed")
	}
}

func TestConsoleStartStops(t *testing.T) {
	c := NewConsole(&config.Fleet{
		Version: 1,
		Ships: map[string]config.Ship{
			"donnager": {URL: "http://donnager:9800"},
		},
	})
	fakeEngine := &fakeEngine{
		events:  make(chan engine.Event),
		started: make(chan struct{}),
	}
	c.engine = fakeEngine
	c.programOptions = quietProgramOptions()

	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Start()
	}()

	select {
	case err := <-errChan:
		t.Fatalf("start returned before stop: %v", err)
	case <-fakeEngine.started:
	case <-time.After(time.Second):
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
	started chan struct{}
	once    sync.Once
}

func (f *fakeEngine) Start() (<-chan engine.Event, error) {
	close(f.started)
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
