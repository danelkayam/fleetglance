package providers

import (
	"fleetglance/internal/protocol"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/mem"
)

func (p *TelemetryProvider) getMemory() *protocol.Memory {
	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting memory telemetry")
		return nil
	}

	return &protocol.Memory{
		UsedBytes:    int64(vm.Used),
		TotalBytes:   int64(vm.Total),
		UsagePercent: vm.UsedPercent,
	}
}
