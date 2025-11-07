package test

import (
	"context"
	"database_sharding/persistent"
	"fmt"
	"testing"
)

func TestUnit(t *testing.T) {
	ctx := context.Background()
	db, err := persistent.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if closeErr := persistent.Close(db); closeErr != nil {
			t.Fatalf("failed to close database: %v", closeErr)
		}
	}()

	tenant, err := persistent.GetTenantWithShard(ctx, db, 20)
	if err != nil {
		t.Fatalf("failed to get tenant: %v", err)
	}
	fmt.Println(tenant)
}
