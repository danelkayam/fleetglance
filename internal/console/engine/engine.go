package engine

import (
	"context"
	"errors"
	"fleetglance/internal/console/config"
	"sort"
	"sync"
	"time"
)

type Engine interface {
	Start() (<-chan Event, error)
	Stop() error
}

type engine struct {
	ships []*Ship

	mu      sync.Mutex
	started bool
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan struct{}
}

func New(fleet *config.Fleet) Engine {
	ctx, cancel := context.WithCancel(context.Background())

	return &engine{
		ships:   toShips(fleet),
		started: false,
		ctx:     ctx,
		cancel:  cancel,
		done:    make(chan struct{}),
	}
}

func (e *engine) Start() (<-chan Event, error) {
	e.mu.Lock()
	if e.started {
		e.mu.Unlock()
		return nil, errors.New("engine already started")
	}
	e.started = true
	ctx := e.ctx
	e.mu.Unlock()

	out := make(chan Event, len(e.ships))
	var wg sync.WaitGroup

	for _, ship := range e.ships {
		wg.Go(func() {
			streamTelemetry(ctx, ship, out)
		})
	}

	go func() {
		wg.Wait()
		close(out)
		close(e.done)
	}()

	return out, nil
}

func (e *engine) Stop() error {
	e.mu.Lock()
	if !e.started {
		e.mu.Unlock()
		return nil
	}
	cancel := e.cancel
	done := e.done
	e.mu.Unlock()

	cancel()
	<-done

	return nil
}

func toShips(fleet *config.Fleet) []*Ship {
	ships := make([]*Ship, 0, len(fleet.Ships))
	names := make([]string, 0, len(fleet.Ships))
	for name := range fleet.Ships {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		shipConfig := fleet.Ships[name]
		ships = append(ships, &Ship{
			Name:         name,
			URL:          shipConfig.URL,
			PullInterval: fleet.PullInterval,
			Timeout:      fleet.Timeout,
		})
	}

	return ships
}

func streamTelemetry(ctx context.Context, ship *Ship, out chan<- Event) {
	client := NewClient(ship.Timeout)

	get := func() bool {
		telemetry, err := client.Get(ctx, ship.URL)
		event := Event{
			ShipName:  ship.Name,
			Telemetry: telemetry,
			Error:     err,
		}

		select {
		case <-ctx.Done():
			return false
		case out <- event:
			return true
		}
	}

	if !get() {
		return
	}

	ticker := time.NewTicker(ship.PullInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !get() {
				return
			}
		}
	}
}
