package engine

import (
	"context"
	"strings"
	"testing"

	"fleetglance/internal/console/config"
)

func TestToShipsSortsFleetShips(t *testing.T) {
	ships := toShips(&config.Fleet{
		Version: 1,
		Ships: map[string]config.Ship{
			"ship-c": {URL: "http://ship-c:9800"},
			"ship-a": {URL: "http://ship-a:9800"},
			"ship-b": {URL: "http://ship-b:9800"},
		},
	})

	got := []string{}
	for _, ship := range ships {
		got = append(got, ship.Name)
	}

	want := []string{"ship-a", "ship-b", "ship-c"}
	if len(got) != len(want) {
		t.Fatalf("ship order = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ship order = %v, want %v", got, want)
		}
	}
}

func TestEngineStartIsSingleUse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	e := &engine{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	events, err := e.Start()
	if err != nil {
		t.Fatalf("start engine: %v", err)
	}
	for range events {
	}

	_, err = e.Start()
	if err == nil {
		t.Fatal("expected second start to fail")
	}
	if !strings.Contains(err.Error(), "engine already started") {
		t.Fatalf("second start error = %q", err.Error())
	}

	if err := e.Stop(); err != nil {
		t.Fatalf("stop engine: %v", err)
	}
}
