package persistent

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func ListTenants(ctx context.Context, db *gorm.DB, limit, offset int) ([]Tenant, error) {
	var tenants []Tenant
	if err := db.WithContext(ctx).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func CountTenants(ctx context.Context, db *gorm.DB) (int64, error) {
	var count int64
	if err := db.WithContext(ctx).Model(&Tenant{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func RebalanceTenants(ctx context.Context, db *gorm.DB) error {
	if err := db.WithContext(ctx).
		Exec("SELECT rebalance_table_shards('tenants')").Error; err != nil {
		return fmt.Errorf("rebalance tenants shards: %w", err)
	}
	return nil
}

type TenantShardPlacement struct {
	ShardID       int64
	NodeName      string
	NodePort      int
	ShardMinValue string
	ShardMaxValue string
}

func ListTenantShardPlacements(ctx context.Context, db *gorm.DB) ([]TenantShardPlacement, error) {
	const shardSQL = `
SELECT
	s.shardid,
	s.shardminvalue,
	s.shardmaxvalue,
	p.nodename,
	p.nodeport
FROM pg_dist_shard s
JOIN pg_dist_shard_placement p ON p.shardid = s.shardid
WHERE s.logicalrelid = 'tenants'::regclass
ORDER BY s.shardid, p.nodename, p.nodeport`
	var placements []TenantShardPlacement
	if err := db.WithContext(ctx).Raw(shardSQL).Scan(&placements).Error; err != nil {
		return nil, fmt.Errorf("list tenant shard placements: %w", err)
	}
	return placements, nil
}

type TenantWithShard struct {
	Tenant
	ShardID  int64  `gorm:"column:shardid"`
	NodeName string `gorm:"column:nodename"`
	NodePort int    `gorm:"column:nodeport"`
}

func GetTenantWithShard(ctx context.Context, db *gorm.DB, id int64) (*TenantWithShard, error) {
	const sql = `
WITH target AS (
	SELECT
		t.*,
    	get_shard_id_for_distribution_column('tenants'::regclass, t.id) AS shardid
	FROM tenants t
	WHERE t.id = ?
	LIMIT 1
)
SELECT
	t.id,
	t.name,
	t.created_at,
	t.updated_at,
	t.shardid,
	p.nodename,
	p.nodeport
FROM target t
JOIN pg_dist_shard_placement p ON p.shardid = t.shardid
LIMIT 1`
	var result TenantWithShard
	if err := db.WithContext(ctx).Raw(sql, id).Scan(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("get tenant shard info: %w", err)
	}
	if result.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}
