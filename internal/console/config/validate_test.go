package config

import (
	"strings"
	"testing"
)

func TestValidateFleetRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name    string
		fleet   *Fleet
		wantErr string
	}{
		{
			name:    "nil fleet",
			fleet:   nil,
			wantErr: "fleet config is required",
		},
		{
			name: "unsupported version",
			fleet: &Fleet{
				Version: 2,
				Ships: map[string]Ship{
					"donnager": {URL: "http://donnager:9800"},
				},
			},
			wantErr: "unsupported fleet config version: 2",
		},
		{
			name: "empty ships",
			fleet: &Fleet{
				Version: 1,
				Ships:   map[string]Ship{},
			},
			wantErr: "fleet must contain at least one ship",
		},
		{
			name: "missing ship url",
			fleet: &Fleet{
				Version: 1,
				Ships: map[string]Ship{
					"donnager": {},
				},
			},
			wantErr: `ship "donnager" url is required`,
		},
		{
			name: "invalid ship url",
			fleet: &Fleet{
				Version: 1,
				Ships: map[string]Ship{
					"donnager": {URL: "donnager:9800"},
				},
			},
			wantErr: `ship "donnager" url must be absolute http/https URL`,
		},
		{
			name:    "too many ships",
			fleet:   testFleetWithShips(MaxShips + 1),
			wantErr: "fleet supports at most 8 ships",
		},
		{
			name: "ship url with path",
			fleet: &Fleet{
				Version: 1,
				Ships: map[string]Ship{
					"donnager": {URL: "http://donnager:9800/api/telemetry"},
				},
			},
			wantErr: `ship "donnager" url must be agent base URL without path`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFleet(tt.fleet)
			if err == nil {
				t.Fatal("expected error")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidateFleetAcceptsValidConfig(t *testing.T) {
	err := ValidateFleet(&Fleet{
		Version: 1,
		Ships: map[string]Ship{
			"donnager": {URL: "http://donnager:9800"},
			"nostromo": {URL: "https://nostromo:9800/"},
		},
	})
	if err != nil {
		t.Fatalf("validate fleet: %v", err)
	}
}

func testFleetWithShips(count int) *Fleet {
	ships := make(map[string]Ship, count)
	for i := range count {
		name := "ship" + string(rune('a'+i))
		ships[name] = Ship{URL: "http://" + name + ":9800"}
	}

	return &Fleet{
		Version: 1,
		Ships:   ships,
	}
}
