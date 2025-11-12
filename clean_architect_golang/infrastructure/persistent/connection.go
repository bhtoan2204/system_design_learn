package persistent

import (
	"clean_architect/config"
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.ConnectionURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open got err=%w", err)
	}

	// if err := db.Use(otelgorm.NewPlugin()); err != nil {
	// 	return nil, fmt.Errorf("new otel-gorm plugin got err=%w", err)
	// }

	dbConfig, err := db.DB()
	if err != nil {
		return nil, err
	}
	dbConfig.SetMaxOpenConns(cfg.Database.MaxOpenConnNumber)
	dbConfig.SetMaxIdleConns(cfg.Database.MaxIdleConnNumber)
	dbConfig.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifeTimeSeconds) * time.Second)
	dbConfig.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxIdleTimeSeconds) * time.Second)

	if err = dbConfig.Ping(); err != nil {
		return nil, fmt.Errorf("ping db got err=%w", err)
	}

	return db, nil
}
