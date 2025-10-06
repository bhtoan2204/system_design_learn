package repository

import (
	"context"
	"event_sourcing_bank_system_api/package/eventsourcing"

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
