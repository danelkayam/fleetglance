package components

import (
	"fleetglance/internal/protocol"
	"time"
)

type Status int

const (
	StatusPending Status = iota
	StatusOnline
	StatusFailed
)

type Model struct {
	Width   int
	Height  int
	Now     time.Time
	Summary Summary
	Ships   []Ship
}

type Summary struct {
	OnlineShips       int
	TotalShips        int
	RunningContainers int
	TotalContainers   int
}

type Ship struct {
	Name      string
	Telemetry *protocol.Telemetry
	Status    Status
}
