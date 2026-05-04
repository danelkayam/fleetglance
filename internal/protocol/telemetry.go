package protocol

import (
	"time"
)

type Telemetry struct {
	Agent_version  string      `json:"agent_version"`
	Hostname       string      `json:"hostname"`
	Timestamp      time.Time   `json:"timestamp"`
	Uptime_seconds int64       `json:"uptime_seconds"`
	Cpu            Cpu         `json:"cpu"`
	Memory         Memory      `json:"memory"`
	Storage        Storage     `json:"storage"`
	Temperature    Temperature `json:"temperature"`
	Containers     Containers  `json:"containers"`
}

type Cpu struct {
	UsagePercent float64 `json:"usage_percent"`
}

type Memory struct {
	UsedBytes    int64   `json:"used_bytes"`
	TotalBytes   int64   `json:"total_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type Storage struct {
	UsedBytes    int64   `json:"used_bytes"`
	TotalBytes   int64   `json:"total_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type Temperature struct {
	Cpu float64 `json:"cpu"`
}

type Containers struct {
	Running int    `json:"running"`
	Total   int    `json:"total"`
	Status  string `json:"status"`
}
