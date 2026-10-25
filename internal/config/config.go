package config

import (
	"fmt"
	"os"
)

type Config struct {
	URL    string
	APIKey string
}

func Load() (Config, error) {
	cfg := Config{
		URL:    os.Getenv("MINIFLUX_URL"),
		APIKey: os.Getenv("MINIFLUX_API_KEY"),
	}

	if cfg.URL == "" {
		return Config{}, fmt.Errorf("MINIFLUX_URL is required")
	}
	if cfg.APIKey == "" {
		return Config{}, fmt.Errorf("MINIFLUX_API_KEY is required")
	}

	return cfg, nil
}
