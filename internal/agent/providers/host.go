package providers

import (
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/host"
)

func (p *provider) getUptimeSeconds() int64 {
	uptime, err := host.Uptime()
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting uptime telemetry")
		return 0
	}

	return int64(uptime)
}
