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
// TODO: replace static data with host, runtime, and container collectors.
func (p *TelemetryProvider) GetTelemetry() (*protocol.Telemetry, error) {
	// TODO - implement this
	return &protocol.Telemetry{
		AgentVersion:  version.Version,
		Hostname:      "mothership",
		Timestamp:     time.Now(),
		UptimeSeconds: 3600,
		Cpu: &protocol.Cpu{
			UsagePercent: 42.5,
		},
		Memory: &protocol.Memory{
			UsedBytes:    1024 * 1024 * 1024,     // 1 GB
			TotalBytes:   8 * 1024 * 1024 * 1024, // 8 GB
			UsagePercent: 12.5,
		},
		Storage: &protocol.Storage{
			UsedBytes:    512 * 1024 * 1024,      // 512 MB
			TotalBytes:   4 * 1024 * 1024 * 1024, // 4 GB
			UsagePercent: 12.8,
		},
		Temperature: &protocol.Temperature{
			Value: 75.5,
		},
		Containers: &protocol.Containers{
			Running: 5,
			Total:   10,
			Status:  "running",
		},
	}, nil
}
