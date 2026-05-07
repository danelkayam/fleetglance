package config

import "time"

const (
	defaultPullInterval = 5 * time.Second
	defaultTimeout      = 2 * time.Second
)

type Fleet struct {
	Version      int             `yaml:"version"`
	PullInterval time.Duration   `yaml:"pull_interval,omitempty"`
	Timeout      time.Duration   `yaml:"timeout,omitempty"`
	Ships        map[string]Ship `yaml:"ships"`
}

type Ship struct {
	URL string `yaml:"url"`
}
