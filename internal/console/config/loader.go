package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFleet(path string) (*Fleet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fleet config %q: %w", path, err)
	}

	var fleet Fleet
	if err := yaml.Unmarshal(data, &fleet); err != nil {
		return nil, fmt.Errorf("parse fleet config %q: %w", path, err)
	}

	return &fleet, nil
}
