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

// GetTelemetry is a placeholder for the actual telemetry data retrieval logic.
// In a real implementation, this would gather data from various sources and return it.
func (p *TelemetryProvider) GetTelemetry() (*protocol.Telemetry, error) {
	// TODO - implement this
	return &protocol.Telemetry{
		Agent_version:  version.Version,
		Hostname:       "example-hostname",
		Timestamp:      time.Now(),
		Uptime_seconds: 3600,
		Cpu: protocol.Cpu{
			UsagePercent: 42.5,
		},
		Memory: protocol.Memory{
			UsedBytes:    1024 * 1024 * 1024,     // 1 GB
			TotalBytes:   8 * 1024 * 1024 * 1024, // 8 GB
			UsagePercent: 12.5,
		},
		Storage: protocol.Storage{
			UsedBytes:    512 * 1024 * 1024,      // 512 MB
			TotalBytes:   4 * 1024 * 1024 * 1024, // 4 GB
			UsagePercent: 12.8,
		},
		Temperature: protocol.Temperature{
			Cpu: 75.5,
		},
		Containers: protocol.Containers{
			Running: 5,
			Total:   10,
			Status:  "running",
		},
	}, nil
}
