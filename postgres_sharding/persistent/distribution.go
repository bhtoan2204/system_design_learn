package persistent

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// EnsureTenantDistribution makes sure the tenants table is distributed across Citus workers.
func EnsureTenantDistribution(ctx context.Context, db *gorm.DB) error {
	const partitionCheck = `
		SELECT EXISTS (
			SELECT 1
			FROM pg_dist_partition
			WHERE logicalrelid = 'tenants'::regclass
		)`
	var exists bool
	if err := db.WithContext(ctx).Raw(partitionCheck).Scan(&exists).Error; err != nil {
		return fmt.Errorf("check tenants distribution: %w", err)
	}
	if exists {
		return nil
	}
	if err := db.WithContext(ctx).
		Exec(`SELECT create_distributed_table('tenants', 'id');`).Error; err != nil {
		return fmt.Errorf("distribute tenants table: %w", err)
	}
	return nil
}
