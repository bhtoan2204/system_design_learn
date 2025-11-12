package migration

import (
	"clean_architect/config"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

type MigrationRunner struct {
	db     *gorm.DB
	config *config.Config
}

func NewMigrationRunner(cfg *config.Config) (*MigrationRunner, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.ConnectionURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &MigrationRunner{
		db:     db,
		config: cfg,
	}, nil
}

func (m *MigrationRunner) ensureMigrationTable() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`
	return m.db.Exec(createTableSQL).Error
}

func (m *MigrationRunner) getAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)

	var results []struct {
		Version int
	}
	if err := m.db.Raw("SELECT version FROM schema_migrations ORDER BY version").Scan(&results).Error; err != nil {
		return nil, err
	}

	for _, r := range results {
		applied[r.Version] = true
	}
	return applied, nil
}

func (m *MigrationRunner) recordMigration(version int, name string) error {
	return m.db.Exec(
		"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
		version, name,
	).Error
}

func (m *MigrationRunner) RunMigrations(ctx context.Context, migrationsDir string) error {
	if err := m.ensureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations, err := m.loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			fmt.Printf("Migration %d_%s already applied, skipping\n", migration.Version, migration.Name)
			continue
		}

		fmt.Printf("Running migration %d_%s...\n", migration.Version, migration.Name)

		tx, err := sqlDB.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if err := m.runMigrationSQL(tx, migration.UpSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to run migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		if err := m.recordMigrationInTx(tx, migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d_%s: %w", migration.Version, migration.Name, err)
		}

		fmt.Printf("Migration %d_%s applied successfully\n", migration.Version, migration.Name)
	}

	return nil
}

func (m *MigrationRunner) loadMigrations(migrationsDir string) ([]Migration, error) {
	var migrations []Migration

	err := filepath.Walk(migrationsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".up.sql") {
			return nil
		}

		baseName := strings.TrimSuffix(filepath.Base(path), ".up.sql")
		parts := strings.SplitN(baseName, "_", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid migration file name: %s (expected: VERSION_NAME.up.sql)", baseName)
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return fmt.Errorf("invalid version in migration file: %s", baseName)
		}

		name := parts[1]

		upSQL, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read up migration file %s: %w", path, err)
		}

		downPath := strings.Replace(path, ".up.sql", ".down.sql", 1)
		downSQL := ""
		if _, err := os.Stat(downPath); err == nil {
			downSQLBytes, err := os.ReadFile(downPath)
			if err == nil {
				downSQL = string(downSQLBytes)
			}
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(upSQL),
			DownSQL: downSQL,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (m *MigrationRunner) runMigrationSQL(tx *sql.Tx, sql string) error {
	lines := strings.Split(sql, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		cleanedLines = append(cleanedLines, line)
	}

	cleanedSQL := strings.Join(cleanedLines, "\n")
	statements := strings.Split(cleanedSQL, ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute SQL: %w\nSQL: %s", err, stmt)
		}
	}
	return nil
}

func (m *MigrationRunner) recordMigrationInTx(tx *sql.Tx, version int, name string) error {
	_, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", version, name)
	return err
}
