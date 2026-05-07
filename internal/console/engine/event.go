package engine

import "fleetglance/internal/protocol"

type Event struct {
	ShipName  string
	Telemetry *protocol.Telemetry
	Error     error
}
