package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func LoadConfig(ctx context.Context) (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	loader := &envconfig.Config{
		Target: cfg,
		Lookuper: envconfig.MultiLookuper(
			envconfig.OsLookuper(),
		),
	}
	if err := envconfig.ProcessWith(ctx, loader); err != nil {
		return nil, fmt.Errorf("envconfig.Process has err=%w", err)
	}
	fmt.Println(cfg)
	return cfg, nil
}
