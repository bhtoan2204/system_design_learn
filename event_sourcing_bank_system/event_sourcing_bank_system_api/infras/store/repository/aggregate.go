package repository

import (
	"context"
	"event_sourcing_bank_system_api/package/eventsourcing"

	"event_sourcing_bank_system_api/package/logger"

	"gorm.io/gorm"
)

type aggregateRepo struct {
	db        *gorm.DB
	serialize eventsourcing.Serializer
}

func newAggregateRepo(db *gorm.DB, s eventsourcing.Serializer) *aggregateRepo {
	return &aggregateRepo{
		serialize: s,
		db:        db,
	}
}

func (r *aggregateRepo) CreateIfNotExist(ctx context.Context, id, typ string) error {
	log := logger.WithPrefix(ctx, "CreateIfNotExist")

	err := r.db.Exec(`
		INSERT IGNORE INTO es_aggregate(id, version, aggregate_type)	
		VALUES(?, 0, ?)
	`, id, typ).Error
	if err != nil {
		log.Warnf("ExecCreateIfNotExist failed with err=%v", err)
		return err
	}

	return nil
}

func (r *aggregateRepo) CheckAndUpdateVersion(ctx context.Context, agg eventsourcing.Aggregate) bool {
	log := logger.WithPrefix(ctx, "CheckAndUpdateVersion")

	root := agg.Root()

	aggregateId := root.AggregateID()
	expectedVersion := root.BaseVersion()
	newVersion := root.Version()

	query := r.db.Exec(`
		UPDATE es_aggregate
		SET version = ?
		WHERE id = ?
			AND version = ?
	`, newVersion, aggregateId, expectedVersion)
	err := query.Error
	if err != nil {
		log.Warnf("CheckAndUpdateVersion err=%v", err)
		return false
	}

	return query.RowsAffected > 0
}

func (r *aggregateRepo) ReadSnapshot(ctx context.Context, aggregateID string, version int, agg eventsourcing.Aggregate) bool {
	log := logger.WithPrefix(ctx, "ReadSnapshot")

	result := struct {
		AggregateType    string `json:"aggregate_type"`
		AggregateVersion int    `json:"aggregate_version"`
		SnapshotVersion  int    `json:"snapshot_version"`
		Data             string `json:"data"`
	}{}

	err := r.db.Raw(`
		SELECT a.aggregate_type, a.version as aggregate_version, eas.version as snapshot_version, eas.data
		FROM es_aggregate_snapshot eas
		JOIN es_aggregate a ON eas.aggregate_id = a.id
		WHERE eas.aggregate_id = ?
			AND eas.version >= ?
		ORDER BY eas.version DESC
		LIMIT 1
		`, aggregateID, version).Scan(&result).Error
	if err != nil {
		return false
	}

	if result.Data == "" {
		return false
	}

	root := agg.Root()

	err = r.serialize.Unmarshal([]byte(result.Data), agg)
	if err != nil {
		log.Warnf("Serialize aggData failed err=%v", err)
		return false
	}

	root.SetInternal(aggregateID, result.SnapshotVersion, result.AggregateVersion)

	return true
}

func (r *aggregateRepo) CreateSnapshot(ctx context.Context, agg eventsourcing.Aggregate) error {
	log := logger.WithPrefix(ctx, "CreateSnapshot")

	root := agg.Root()
	aggregateId := root.AggregateID()
	version := root.Version()
	data, err := r.serialize.Marshal(agg)
	if err != nil {
		return err
	}

	err = r.db.Exec(`
		INSERT INTO es_aggregate_snapshot (aggregate_id, version, data)
		VALUES(?, ?, ?)`, aggregateId, version, data).Error
	if err != nil {
		log.Warnf("Insert es_aggregate_snapshot err=%v", err)
		return err
	}

	return nil
}
