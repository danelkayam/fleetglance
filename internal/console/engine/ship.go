package engine

import "time"

type Ship struct {
	Name         string
	URL          string
	PullInterval time.Duration
	Timeout      time.Duration
}
