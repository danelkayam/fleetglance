package providers

import (
	"fleetglance/internal/protocol"
	"fleetglance/internal/version"
	"fmt"
	"time"
)

type TelemetryProvider interface {
	GetTelemetry() (*protocol.Telemetry, error)
}

type provider struct{}

func New() TelemetryProvider {
	return &provider{}
}

// GetTelemetry returns the current agent telemetry.
func (p *provider) GetTelemetry() (*protocol.Telemetry, error) {
	return &protocol.Telemetry{
		AgentVersion:  fmt.Sprintf("%s-%s", version.Version, version.Commit),
		Timestamp:     time.Now(),
		UptimeSeconds: p.getUptimeSeconds(),
		CPU:           p.getCPU(),
		Memory:        p.getMemory(),
		Storage:       p.getStorage(),
		Containers:    p.getContainers(),
	}, nil
}
