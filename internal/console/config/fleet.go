package config

type Fleet struct {
	Version int             `yaml:"version"`
	Ships   map[string]Ship `yaml:"ships"`
}

type Ship struct {
	URL string `yaml:"url"`
}
