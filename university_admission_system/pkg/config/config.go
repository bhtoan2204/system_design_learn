package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents application runtime configuration.
type Config struct {
	HTTPPort      string
	MinimumScore  float64
	SeedDemoData  bool
	EnableSwagger bool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:      defaultString(os.Getenv("APP_HTTP_PORT"), "8080"),
		MinimumScore:  75,
		SeedDemoData:  defaultBool(os.Getenv("APP_SEED_DATA"), true),
		EnableSwagger: defaultBool(os.Getenv("APP_ENABLE_SWAGGER"), true),
	}

	if v := os.Getenv("APP_MINIMUM_SCORE"); v != "" {
		score, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid APP_MINIMUM_SCORE: %w", err)
		}
		cfg.MinimumScore = score
	}

	return cfg, nil
}

func defaultString(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func defaultBool(value string, fallback bool) bool {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
