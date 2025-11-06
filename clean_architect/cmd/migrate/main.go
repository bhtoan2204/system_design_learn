package main

import (
	"clean_architect/config"
	"clean_architect/infrastructure/persistent/migration"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var migrationsDir string
	flag.StringVar(&migrationsDir, "dir", "infrastructure/persistent/migration/migrations", "Directory containing migration files")
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	runner, err := migration.NewMigrationRunner(cfg)
	if err != nil {
		log.Fatalf("Failed to create migration runner: %v", err)
	}

	// Get absolute path
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", absPath)
	}

	fmt.Printf("Running migrations from: %s\n", absPath)

	if err := runner.RunMigrations(ctx, absPath); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("All migrations completed successfully!")
}
