package providers

import (
	"fleetglance/internal/protocol"
	"fleetglance/internal/version"
	"time"
)

type TelemetryProvider struct{}

func NewTelemetryProvider() *TelemetryProvider {
	return &TelemetryProvider{}
}

// GetTelemetry returns the current agent telemetry.
func (p *TelemetryProvider) GetTelemetry() (*protocol.Telemetry, error) {
	return &protocol.Telemetry{
		AgentVersion:  version.Version,
		Hostname:      p.getHostname(),
		Timestamp:     time.Now(),
		UptimeSeconds: p.getUptimeSeconds(),
		CPU:           p.getCPU(),
		Memory:        p.getMemory(),
		Storage:       p.getStorage(),
		Containers:    p.getContainers(),
	}, nil
}
