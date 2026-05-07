package providers

import "math"

func roundPercent(value float64) float64 {
	return math.Round(value*10) / 10
}
