package providers

import (
	"fleetglance/internal/protocol"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/v4/disk"
)

func (p *provider) getStorage() *protocol.Storage {
	usage, err := disk.Usage("/")
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting storage telemetry")
		return nil
	}

	return &protocol.Storage{
		UsedBytes:    int64(usage.Used),
		TotalBytes:   int64(usage.Total),
		UsagePercent: usage.UsedPercent,
	}
}
