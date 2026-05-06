package providers

import (
	"context"
	"time"

	"fleetglance/internal/protocol"

	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

const dockerTelemetryTimeout = 2 * time.Second

func (p *provider) getContainers() *protocol.Containers {
	apiClient, err := client.New(
		client.WithTimeout(dockerTelemetryTimeout),
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting containers telemetry")
		return nil
	}
	defer func() {
		_ = apiClient.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), dockerTelemetryTimeout)
	defer cancel()

	result, err := apiClient.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed collecting containers telemetry")
		return nil
	}

	containers := &protocol.Containers{
		Total: len(result.Items),
	}
	for _, item := range result.Items {
		if item.State == "running" {
			containers.Running++
		}
	}

	return containers
}
