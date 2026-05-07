package engine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fleetglance/internal/protocol"
)

func TestClientGetTelemetrySuccess(t *testing.T) {
	wantTelemetry := protocol.Telemetry{
		AgentVersion:  "test-version",
		Timestamp:     time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC),
		UptimeSeconds: 123,
		CPU: &protocol.CPU{
			UsagePercent: 12.5,
		},
		Memory: &protocol.Memory{
			UsedBytes:    1024,
			TotalBytes:   2048,
			UsagePercent: 50,
		},
		Storage: &protocol.Storage{
			UsedBytes:    4096,
			TotalBytes:   8192,
			UsagePercent: 50,
		},
		Containers: &protocol.Containers{
			Running: 2,
			Total:   3,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/telemetry" {
			t.Fatalf("expected telemetry path, got %q", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Fatalf("expected method %s, got %s", http.MethodGet, r.Method)
		}
		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Fatalf("expected Accept application/json, got %q", accept)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(protocol.Response[protocol.Telemetry]{
			Data:  &wantTelemetry,
			Error: nil,
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	got, err := NewClient(time.Second).Get(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("get telemetry: %v", err)
	}

	if got.AgentVersion != wantTelemetry.AgentVersion {
		t.Fatalf("expected agent version %q, got %q", wantTelemetry.AgentVersion, got.AgentVersion)
	}
	if !got.Timestamp.Equal(wantTelemetry.Timestamp) {
		t.Fatalf("expected timestamp %s, got %s", wantTelemetry.Timestamp, got.Timestamp)
	}
	if got.CPU == nil || got.CPU.UsagePercent != wantTelemetry.CPU.UsagePercent {
		t.Fatalf("expected cpu telemetry %#v, got %#v", wantTelemetry.CPU, got.CPU)
	}
}

func TestClientGetTelemetryDoesNotDuplicateTelemetryPath(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		telemetry := protocol.Telemetry{AgentVersion: "test-version"}
		if err := json.NewEncoder(w).Encode(protocol.Response[protocol.Telemetry]{
			Data: &telemetry,
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	_, err := NewClient(time.Second).Get(context.Background(), server.URL+"/api/telemetry")
	if err != nil {
		t.Fatalf("get telemetry: %v", err)
	}

	if gotPath != "/api/telemetry" {
		t.Fatalf("expected existing telemetry path to be reused, got %q", gotPath)
	}
}

func TestClientGetTelemetryErrors(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr string
	}{
		{
			name: "non ok status",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
			},
			wantErr: "unexpected status 502",
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("{"))
			},
			wantErr: "decode telemetry response",
		},
		{
			name: "error envelope",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				err := json.NewEncoder(w).Encode(protocol.Response[protocol.Telemetry]{
					Error: &protocol.ResponseError{
						Message: "ship offline",
					},
				})
				if err != nil {
					t.Fatalf("encode response: %v", err)
				}
			},
			wantErr: "telemetry response error: ship offline",
		},
		{
			name: "missing data",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				err := json.NewEncoder(w).Encode(protocol.Response[protocol.Telemetry]{})
				if err != nil {
					t.Fatalf("encode response: %v", err)
				}
			},
			wantErr: "telemetry response missing data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			_, err := NewClient(time.Second).Get(context.Background(), server.URL)
			if err == nil {
				t.Fatal("expected error")
			}

			if got := err.Error(); !strings.Contains(got, tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, got)
			}
		})
	}
}
