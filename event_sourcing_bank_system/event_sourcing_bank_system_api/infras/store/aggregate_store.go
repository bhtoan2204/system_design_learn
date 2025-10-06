package store

import (
	"context"
	"errors"
	"event_sourcing_bank_system_api/infras/store/repository"
	"event_sourcing_bank_system_api/package/eventsourcing"
	"event_sourcing_bank_system_api/package/logger"
	"fmt"
	"reflect"
)

var _ AggregateStore = (*aggregateStore)(nil)

type AggregateStore interface {
	Get(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) error
	Save(ctx context.Context, agg eventsourcing.Aggregate) error
}

func NewAggregateStore(repository repository.Repos) AggregateStore {
	return &aggregateStore{
		repository: repository.EventStore(),
	}
}

type aggregateStore struct {
	repository repository.EventStore
}

// Get fetches the events and build up the aggregate
func (as *aggregateStore) Get(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) error {
	log := logger.WithPrefix(ctx, "Get")

	if reflect.ValueOf(agg).Kind() != reflect.Ptr {
		log.Warnf("aggregate must to be a pointer, aggregate_id=%s", aggregateID)
		return errors.New("aggregate must to be a pointer")
	}

	aggType := reflect.TypeOf(agg).Elem().Name()
	agg.Root().SetAggregateType(aggType)

	// FIXME: getFromSnapshot should be configurable base on aggregate
	has := as.getFromSnapshot(ctx, aggregateID, agg)
	if !has {
		return as.getFromEvents(ctx, aggregateID, agg)
	}

	return nil
}

func (as *aggregateStore) Save(ctx context.Context, agg eventsourcing.Aggregate) error {
	log := logger.WithPrefix(ctx, "Save")

	root := agg.Root()
	aggType := reflect.TypeOf(agg).Elem().Name()
	root.SetAggregateType(aggType)

	err := as.repository.CreateIfNotExist(ctx, root.AggregateID(), aggType)
	if err != nil {
		return err
	}

	err = as.repository.WithTransaction(ctx, func(txn repository.EventStore) error {
		if !txn.CheckAndUpdateVersion(ctx, agg) {
			return fmt.Errorf("optimistic concurrency control failed id=%s, expectedVersion=%d, newversion=%d",
				root.AggregateID(), root.BaseVersion(), root.Version())
		}

		for _, event := range root.Events() {
			err := txn.Append(ctx, event)
			if err != nil {
				return err
			}

			// FIXME: createSnapshot should be configurable
			nthEvent := 10
			version := root.Version()
			if version%nthEvent == 0 {
				log.Infof("Create snapshot of aggID=%s, aggType=%s", root.AggregateID(), root.AggregateType())
				err = txn.CreateSnapshot(ctx, agg)
				if err != nil {
					log.Warnf("Create snapshot aggID=%s got err=%v", root.AggregateID(), err)
					return err
				}
			}
			root.Update()
		}

		return nil
	})
	if err != nil {
		log.Warnf("Save EventStore with transaction got err=%v", err)
		return err
	}

	return nil
}

func (as *aggregateStore) getFromSnapshot(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) bool {
	log := logger.WithPrefix(ctx, "getFromSnapshot")

	has := as.repository.ReadSnapshot(ctx, aggregateID, agg.Root().Version(), agg)
	if !has {
		return false
	}

	root := agg.Root()
	if root.BaseVersion() < root.Version() {
		err := as.repository.Get(ctx, aggregateID, root.BaseVersion(), root.Version(), agg)
		if err != nil {
			log.Warnf("Get event fromVersion=%d to toVersion=%s err=%v", root.BaseVersion(), root.Version(), err)
			return false
		}
	}

	return true
}

func (as *aggregateStore) getFromEvents(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) error {
	err := as.repository.Get(ctx, aggregateID, 0, 0, agg)
	if err != nil {
		return err
	}

	return nil
}
