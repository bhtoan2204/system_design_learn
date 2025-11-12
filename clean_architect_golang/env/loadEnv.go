package env

import (
	"clean_architect/config"
	"clean_architect/infrastructure/persistent"
	"context"
)

func LoadEnv(ctx context.Context, cfg *config.Config) (*Env, error) {
	var envOption []Option

	// init database
	db, err := persistent.InitDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}
	envOption = append(envOption, WithDatabase(db))

	// init env
	env := NewEnv(envOption...)
	return env, nil
}
