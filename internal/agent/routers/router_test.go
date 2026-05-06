package routers

import (
	"encoding/json"
	"errors"
	"fleetglance/internal/protocol"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fakeTelemetryProvider struct {
	telemetry *protocol.Telemetry
	err       error
}

func (f fakeTelemetryProvider) GetTelemetry() (*protocol.Telemetry, error) {
	return f.telemetry, f.err
}

func TestHealthz(t *testing.T) {
	router := NewRouter(Params{
		Debug: false,
		TelemetryProvider: fakeTelemetryProvider{
			telemetry: testTelemetry(),
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body["status"] != "healthy" {
		t.Fatalf("expected healthy status body, got %#v", body)
	}
}

func TestTelemetrySuccess(t *testing.T) {
	router := NewRouter(Params{
		Debug: false,
		TelemetryProvider: fakeTelemetryProvider{
			telemetry: testTelemetry(),
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/telemetry", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode telemetry response: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected top-level data object, got %#v", body["data"])
	}

	assertKeys(t, data, "agent_version", "timestamp", "uptime_seconds", "cpu", "memory", "storage", "containers")

	if _, ok := data["agentVersion"]; ok {
		t.Fatal("expected agent_version to remain snake_case")
	}
	if _, ok := data["uptimeSeconds"]; ok {
		t.Fatal("expected uptime_seconds to remain snake_case")
	}

	if data["agent_version"] != "test-version" {
		t.Fatalf("expected agent_version %q, got %#v", "test-version", data["agent_version"])
	}
	if data["uptime_seconds"] != float64(123) {
		t.Fatalf("expected uptime_seconds 123, got %#v", data["uptime_seconds"])
	}
	if _, err := time.Parse(time.RFC3339Nano, data["timestamp"].(string)); err != nil {
		t.Fatalf("expected parseable timestamp, got %q: %v", data["timestamp"], err)
	}

	cpu := objectField(t, data, "cpu")
	assertKeys(t, cpu, "usage_percent")
	if cpu["usage_percent"] != 12.5 {
		t.Fatalf("expected cpu usage_percent 12.5, got %#v", cpu["usage_percent"])
	}

	memory := objectField(t, data, "memory")
	assertKeys(t, memory, "used_bytes", "total_bytes", "usage_percent")
	if memory["used_bytes"] != float64(1024) {
		t.Fatalf("expected memory used_bytes 1024, got %#v", memory["used_bytes"])
	}
	if memory["total_bytes"] != float64(2048) {
		t.Fatalf("expected memory total_bytes 2048, got %#v", memory["total_bytes"])
	}
	if memory["usage_percent"] != 50.0 {
		t.Fatalf("expected memory usage_percent 50, got %#v", memory["usage_percent"])
	}

	storage := objectField(t, data, "storage")
	assertKeys(t, storage, "used_bytes", "total_bytes", "usage_percent")
	if storage["used_bytes"] != float64(4096) {
		t.Fatalf("expected storage used_bytes 4096, got %#v", storage["used_bytes"])
	}
	if storage["total_bytes"] != float64(8192) {
		t.Fatalf("expected storage total_bytes 8192, got %#v", storage["total_bytes"])
	}
	if storage["usage_percent"] != 50.0 {
		t.Fatalf("expected storage usage_percent 50, got %#v", storage["usage_percent"])
	}

	containers := objectField(t, data, "containers")
	assertKeys(t, containers, "running", "total")
	if containers["running"] != float64(2) {
		t.Fatalf("expected containers running 2, got %#v", containers["running"])
	}
	if containers["total"] != float64(3) {
		t.Fatalf("expected containers total 3, got %#v", containers["total"])
	}
}

func TestTelemetryPartialSections(t *testing.T) {
	telemetry := testTelemetry()
	telemetry.Memory = nil
	telemetry.Containers = nil

	router := NewRouter(Params{
		Debug: false,
		TelemetryProvider: fakeTelemetryProvider{
			telemetry: telemetry,
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/telemetry", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode telemetry response: %v", err)
	}

	data := objectField(t, body, "data")
	if _, ok := data["memory"]; !ok {
		t.Fatal("expected memory field to be present")
	}
	if data["memory"] != nil {
		t.Fatalf("expected memory to serialize as null, got %#v", data["memory"])
	}
	if _, ok := data["containers"]; !ok {
		t.Fatal("expected containers field to be present")
	}
	if data["containers"] != nil {
		t.Fatalf("expected containers to serialize as null, got %#v", data["containers"])
	}
}

func TestTelemetryProviderError(t *testing.T) {
	router := NewRouter(Params{
		Debug: false,
		TelemetryProvider: fakeTelemetryProvider{
			err: errors.Join(protocol.ErrTelemetryFailed, errors.New("collector failed")),
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/telemetry", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if body["data"] != nil {
		t.Fatalf("expected data to be null, got %#v", body["data"])
	}

	responseError := objectField(t, body, "error")
	message, ok := responseError["message"].(string)
	if !ok {
		t.Fatalf("expected error message string, got %#v", responseError["message"])
	}
	if message == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestUnknownRoute(t *testing.T) {
	router := NewRouter(Params{
		Debug: false,
		TelemetryProvider: fakeTelemetryProvider{
			telemetry: testTelemetry(),
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func testTelemetry() *protocol.Telemetry {
	return &protocol.Telemetry{
		AgentVersion:  "test-version",
		Timestamp:     time.Date(2026, 5, 6, 12, 0, 0, 123456789, time.UTC),
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
}

func objectField(t *testing.T, object map[string]any, name string) map[string]any {
	t.Helper()

	field, ok := object[name].(map[string]any)
	if !ok {
		t.Fatalf("expected %s object, got %#v", name, object[name])
	}

	return field
}

func assertKeys(t *testing.T, object map[string]any, keys ...string) {
	t.Helper()

	for _, key := range keys {
		if _, ok := object[key]; !ok {
			t.Fatalf("expected key %q in %#v", key, object)
		}
	}
}
