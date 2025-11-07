package persistent

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultDatabaseURL = "postgres://postgres:secret@localhost:5432/appdb?sslmode=disable"

func Connect(ctx context.Context) (*gorm.DB, error) {
	dsn := databaseURL()
	var (
		db  *gorm.DB
		err error
	)

	backoff := time.Second
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, derr := db.DB()
			if derr == nil {
				if derr = sqlDB.Ping(); derr == nil {
					return db, nil
				}
			}
			err = derr
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if attempt == 10 {
			break
		}
		time.Sleep(backoff)
		if backoff < 5*time.Second {
			backoff *= 2
			if backoff > 5*time.Second {
				backoff = 5 * time.Second
			}
		}
	}
	return nil, fmt.Errorf("connect to database: %w", err)
}

func Close(db *gorm.DB) error {
	dbConfig, err := db.DB()
	if err != nil {
		return err
	}
	return dbConfig.Close()
}

func databaseURL() string {
	if env := os.Getenv("DATABASE_URL"); env != "" {
		return env
	}
	return defaultDatabaseURL
}
