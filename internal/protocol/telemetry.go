package protocol

import (
	"time"
)

type Telemetry struct {
	AgentVersion  string      `json:"agent_version"`
	Timestamp     time.Time   `json:"timestamp"`
	UptimeSeconds int64       `json:"uptime_seconds"`
	CPU           *CPU        `json:"cpu"`
	Memory        *Memory     `json:"memory"`
	Storage       *Storage    `json:"storage"`
	Containers    *Containers `json:"containers"`
}

type CPU struct {
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

type Containers struct {
	Running int `json:"running"`
	Total   int `json:"total"`
}
