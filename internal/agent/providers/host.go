package providers

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/host"
)

func (p *TelemetryProvider) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting hostname")
		return ""
	}

	return hostname
}

func (p *TelemetryProvider) getUptimeSeconds() int64 {
	uptime, err := host.Uptime()
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting uptime telemetry")
		return 0
	}

	return int64(uptime)
}
