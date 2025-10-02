package repos

import (
	"context"
	"fmt"

	"event_sourcing_golang/pkg/eventsourcing"

	"gorm.io/gorm"
)

var (
	_ EventStore = (*eventStore)(nil)
)

type EventStore interface {
	Get(ctx context.Context, aggregateID string, fromVersion, toVersion int, agg eventsourcing.Aggregate) error
	Append(context.Context, eventsourcing.Event) error

	CreateIfNotExist(ctx context.Context, id, typ string) error
	CheckAndUpdateVersion(ctx context.Context, agg eventsourcing.Aggregate) bool
	CreateSnapshot(ctx context.Context, agg eventsourcing.Aggregate) error
	ReadSnapshot(ctx context.Context, aggregateID string, version int, agg eventsourcing.Aggregate) bool
	WithTransaction(ctx context.Context, fn func(EventStore) error) (err error)

	// List returns all events for an aggregate, fully deserialized using agg's registered types
	List(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) ([]eventsourcing.Event, error)
	// SnapshotVersion returns the latest snapshot version for the aggregate if exists
	SnapshotVersion(ctx context.Context, aggregateID string) (int, bool)
}

type eventStore struct {
	*aggregateRepo
	*eventRepo
	db         *gorm.DB
	serializer eventsourcing.Serializer
}

func newEventStore(db *gorm.DB, s eventsourcing.Serializer) EventStore {
	return &eventStore{
		db:            db,
		serializer:    s,
		aggregateRepo: newAggregateRepo(db, s),
		eventRepo:     newEventRepo(db, s),
	}
}

func (r *eventStore) WithTransaction(ctx context.Context, fn func(EventStore) error) (err error) {
	tx := r.db.Begin()
	tr := &eventStore{
		db:            tx,
		aggregateRepo: newAggregateRepo(tx, r.serializer),
		eventRepo:     newEventRepo(tx, r.serializer),
	}
	err = tx.Error
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil { // nolint
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = fn(tr)

	return err
}

// List returns all events for an aggregate, deserialized using the aggregate's registered event types
func (r *eventStore) List(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) ([]eventsourcing.Event, error) {
	root := agg.Root()
	rows, err := r.db.Raw(`
            SELECT id, aggregate_id, event_type, version, data, metadata, created_at
            FROM es_event
            WHERE aggregate_id = ?
            ORDER BY version ASC`, aggregateID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []eventsourcing.Event
	for rows.Next() {
		var evt eventsourcing.Event
		var dataStr string
		var metaStr string
		if err := rows.Scan(&evt.ID, &evt.AggregateID, &evt.EventType, &evt.Version, &dataStr, &metaStr, &evt.CreatedAt); err != nil {
			return nil, err
		}

		f, ok := r.serializer.Type(root.AggregateType(), evt.EventType)
		if !ok {
			return nil, fmt.Errorf("cant serialize event with type: %s_%s", root.AggregateType(), evt.EventType)
		}
		eventData := f()
		if err := r.serializer.Unmarshal([]byte(dataStr), &eventData); err != nil {
			return nil, err
		}
		var metadata interface{}
		if metaStr != "" {
			if err := r.serializer.Unmarshal([]byte(metaStr), &metadata); err != nil {
				return nil, err
			}
		}
		evt.Data = eventData
		evt.Metadata = metadata
		result = append(result, evt)
	}
	return result, nil
}

// SnapshotVersion returns the latest snapshot version if exists
func (r *eventStore) SnapshotVersion(ctx context.Context, aggregateID string) (int, bool) {
	var v int
	err := r.db.Raw(`
        SELECT version
        FROM es_aggregate_snapshot
        WHERE aggregate_id = ?
        ORDER BY version DESC
        LIMIT 1
    `, aggregateID).Scan(&v).Error
	if err != nil {
		return 0, false
	}
	if v == 0 {
		return 0, false
	}
	return v, true
}
