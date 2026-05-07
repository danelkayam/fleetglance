package console

import (
	"fleetglance/internal/console/config"
	"strings"
	"testing"
	"time"
)

func TestConsoleStartValidatesFleet(t *testing.T) {
	tests := []struct {
		name    string
		fleet   *config.Fleet
		wantErr string
	}{
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

func TestConsoleStartStops(t *testing.T) {
	c := NewConsole(&config.Fleet{
		Version: 1,
		Ships: map[string]config.Ship{
			"donnager": {URL: "http://donnager:9800"},
		},
	})

	errChan := make(chan error, 1)
	go func() {
		errChan <- c.Start()
	}()

	select {
	case err := <-errChan:
		t.Fatalf("start returned before stop: %v", err)
	case <-time.After(10 * time.Millisecond):
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
