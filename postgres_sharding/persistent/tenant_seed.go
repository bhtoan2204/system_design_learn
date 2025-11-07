package persistent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/gorm"
)

const tenantSeedFile = "mock/tenants.json"

type tenantSeed struct {
	Name string `json:"name"`
}

var (
	tenantSeedOnce sync.Once
	tenantSeedData []tenantSeed
	tenantSeedErr  error
)

func loadTenantSeeds() ([]tenantSeed, error) {
	tenantSeedOnce.Do(func() {
		path := tenantSeedFile
		if envPath := os.Getenv("TENANT_SEED_FILE"); envPath != "" {
			path = envPath
		}
		data, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			tenantSeedErr = fmt.Errorf("read tenant seed file: %w", err)
			return
		}
		if err := json.Unmarshal(data, &tenantSeedData); err != nil {
			tenantSeedErr = fmt.Errorf("unmarshal tenant seed data: %w", err)
			return
		}
		if len(tenantSeedData) == 0 {
			tenantSeedErr = errors.New("tenant seed file is empty")
		}
	})
	return tenantSeedData, tenantSeedErr
}

// SeedTenants inserts mock tenants from a JSON file so sharding distribution can be inspected.
func SeedTenants(ctx context.Context, db *gorm.DB) error {
	seeds, err := loadTenantSeeds()
	if err != nil {
		return err
	}

	const batchSize = 1000
	for start := 0; start < len(seeds); start += batchSize {
		end := start + batchSize
		if end > len(seeds) {
			end = len(seeds)
		}

		candidates := make([]Tenant, 0, end-start)
		names := make([]string, 0, end-start)
		for _, seed := range seeds[start:end] {
			if seed.Name == "" {
				continue
			}
			names = append(names, seed.Name)
			candidates = append(candidates, Tenant{Name: seed.Name})
		}
		if len(candidates) == 0 {
			continue
		}

		var existing []string
		if err := db.WithContext(ctx).
			Model(&Tenant{}).
			Where("name IN ?", names).
			Pluck("name", &existing).Error; err != nil {
			return fmt.Errorf("load existing tenants for batch starting at %d: %w", start, err)
		}
		if len(existing) > 0 {
			existingSet := make(map[string]struct{}, len(existing))
			for _, name := range existing {
				existingSet[name] = struct{}{}
			}

			filtered := candidates[:0]
			for _, tenant := range candidates {
				if _, found := existingSet[tenant.Name]; !found {
					filtered = append(filtered, tenant)
				}
			}
			candidates = filtered
		}

		if len(candidates) == 0 {
			continue
		}

		if err := db.WithContext(ctx).Create(&candidates).Error; err != nil {
			return fmt.Errorf("seed tenants batch starting at %d: %w", start, err)
		}
	}

	return nil
}
