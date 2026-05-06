package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTelemetryJSONShape(t *testing.T) {
	telemetry := Telemetry{
		AgentVersion:  "test-version",
		Timestamp:     time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC),
		UptimeSeconds: 123,
		CPU: &CPU{
			UsagePercent: 12.5,
		},
		Memory: &Memory{
			UsedBytes:    1024,
			TotalBytes:   2048,
			UsagePercent: 50,
		},
		Storage: &Storage{
			UsedBytes:    4096,
			TotalBytes:   8192,
			UsagePercent: 50,
		},
		Containers: &Containers{
			Running: 2,
			Total:   3,
		},
	}

	var body map[string]any
	marshalToMap(t, telemetry, &body)

	assertJSONKeys(t, body, "agent_version", "timestamp", "uptime_seconds", "cpu", "memory", "storage", "containers")
	assertNoJSONKeys(t, body, "agentVersion", "uptimeSeconds")

	cpu := jsonObject(t, body, "cpu")
	assertJSONKeys(t, cpu, "usage_percent")

	memory := jsonObject(t, body, "memory")
	assertJSONKeys(t, memory, "used_bytes", "total_bytes", "usage_percent")
	assertNoJSONKeys(t, memory, "usedBytes", "totalBytes", "usagePercent")

	storage := jsonObject(t, body, "storage")
	assertJSONKeys(t, storage, "used_bytes", "total_bytes", "usage_percent")
	assertNoJSONKeys(t, storage, "usedBytes", "totalBytes", "usagePercent")

	containers := jsonObject(t, body, "containers")
	assertJSONKeys(t, containers, "running", "total")
}

func TestTelemetryNilSectionsSerializeAsNull(t *testing.T) {
	telemetry := Telemetry{
		AgentVersion:  "test-version",
		Timestamp:     time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC),
		UptimeSeconds: 123,
		CPU:           nil,
		Memory:        nil,
		Storage:       nil,
		Containers:    nil,
	}

	var body map[string]any
	marshalToMap(t, telemetry, &body)

	for _, key := range []string{"cpu", "memory", "storage", "containers"} {
		if _, ok := body[key]; !ok {
			t.Fatalf("expected key %q in %#v", key, body)
		}
		if body[key] != nil {
			t.Fatalf("expected %s to serialize as null, got %#v", key, body[key])
		}
	}
}

func TestResponseEnvelopeJSONShape(t *testing.T) {
	response := Response{
		Data: map[string]string{
			"status": "ok",
		},
		Error: nil,
	}

	var body map[string]any
	marshalToMap(t, response, &body)

	assertJSONKeys(t, body, "data")
	assertNoJSONKeys(t, body, "error")
	if _, ok := body["data"].(map[string]any); !ok {
		t.Fatalf("expected data object, got %#v", body["data"])
	}
}

func TestErrorResponseEnvelopeJSONShape(t *testing.T) {
	response := Response{
		Data: nil,
		Error: &ResponseError{
			Message: "failed",
		},
	}

	var body map[string]any
	marshalToMap(t, response, &body)

	assertJSONKeys(t, body, "data", "error")
	if body["data"] != nil {
		t.Fatalf("expected data to serialize as null, got %#v", body["data"])
	}

	responseError := jsonObject(t, body, "error")
	assertJSONKeys(t, responseError, "message")
	if responseError["message"] != "failed" {
		t.Fatalf("expected error message %q, got %#v", "failed", responseError["message"])
	}
}

func marshalToMap(t *testing.T, value any, target any) {
	t.Helper()

	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal value: %v", err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("unmarshal value: %v", err)
	}
}

func jsonObject(t *testing.T, object map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := object[key].(map[string]any)
	if !ok {
		t.Fatalf("expected %s object, got %#v", key, object[key])
	}

	return value
}

func assertJSONKeys(t *testing.T, object map[string]any, keys ...string) {
	t.Helper()

	for _, key := range keys {
		if _, ok := object[key]; !ok {
			t.Fatalf("expected key %q in %#v", key, object)
		}
	}
}

func assertNoJSONKeys(t *testing.T, object map[string]any, keys ...string) {
	t.Helper()

	for _, key := range keys {
		if _, ok := object[key]; ok {
			t.Fatalf("did not expect key %q in %#v", key, object)
		}
	}
}
