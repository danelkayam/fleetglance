package components

import (
	"fmt"
	"math"
	"time"
)

func FormatPercent(value float64) string {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		value = 0
	}

	if value < 0 {
		value = 0
	}

	if value > 100 {
		value = 100
	}

	return fmt.Sprintf("%.1f%%", value)
}

func FormatUptime(seconds int64) string {
	if seconds <= 0 {
		return "0m"
	}

	duration := time.Duration(seconds) * time.Second
	days := duration / (24 * time.Hour)
	duration -= days * 24 * time.Hour
	hours := duration / time.Hour
	duration -= hours * time.Hour
	minutes := duration / time.Minute

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh", days, hours)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}

func FormatContainers(running int, total int) string {
	return fmt.Sprintf("%d/%d", running, total)
}
