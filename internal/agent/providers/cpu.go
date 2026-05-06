package providers

import (
	"fleetglance/internal/protocol"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/cpu"
)

func (p *provider) getCPU() *protocol.CPU {
	values, err := cpu.Percent(0, false)
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting CPU telemetry")
		return nil
	}

	if len(values) == 0 {
		log.Warn().Msg("Failed collecting CPU telemetry: no usage values returned")
		return nil
	}

	return &protocol.CPU{
		UsagePercent: values[0],
	}
}
