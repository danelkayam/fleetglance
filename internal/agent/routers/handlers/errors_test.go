package handlers

import (
	"errors"
	"fleetglance/internal/protocol"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestResolveStatusAndMessage(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "invalid args",
			err:         protocol.ErrInvalidArgs,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Invalid request: invalid arguments",
		},
		{
			name:        "telemetry unavailable",
			err:         protocol.ErrTelemetryUnavailable,
			wantStatus:  http.StatusServiceUnavailable,
			wantMessage: "Telemetry unavailable: telemetry unavailable",
		},
		{
			name:        "telemetry failed",
			err:         protocol.ErrTelemetryFailed,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Telemetry collection failed: telemetry collection failed",
		},
		{
			name:        "unauthenticated",
			err:         protocol.ErrUnauthenticated,
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Unauthorized: unauthenticated",
		},
		{
			name:        "unauthorized",
			err:         protocol.ErrUnauthorized,
			wantStatus:  http.StatusForbidden,
			wantMessage: "Forbidden: unauthorized",
		},
		{
			name:        "unknown",
			err:         errors.New("boom"),
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Internal server error: boom",
		},
		{
			name:        "wrapped protocol error",
			err:         fmt.Errorf("collector wrapped: %w", protocol.ErrTelemetryUnavailable),
			wantStatus:  http.StatusServiceUnavailable,
			wantMessage: "Telemetry unavailable: collector wrapped: telemetry unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, message := resolveStatusAndMessage(tt.err)

			if status != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, status)
			}
			if message != tt.wantMessage {
				t.Fatalf("expected message %q, got %q", tt.wantMessage, message)
			}
		})
	}
}

func TestResolveStatusAndMessagePrefersFirstMatchingProtocolError(t *testing.T) {
	err := errors.Join(protocol.ErrInvalidArgs, protocol.ErrUnauthorized)

	status, message := resolveStatusAndMessage(err)

	if status != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, status)
	}
	if !strings.HasPrefix(message, "Invalid request:") {
		t.Fatalf("expected invalid request message, got %q", message)
	}
}
