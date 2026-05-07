package ui

import (
	"fleetglance/internal/console/engine"
	"time"
)

type TelemetryEventMsg engine.Event

type tickMsg struct {
	now time.Time
}
