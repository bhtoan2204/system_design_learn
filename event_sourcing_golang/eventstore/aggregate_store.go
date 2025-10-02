package eventstore

import (
	"context"
	"errors"
	"event_sourcing_golang/eventstore/repos"
	"fmt"
	"reflect"

	"event_sourcing_golang/pkg/eventsourcing"
)

var _ AggregateStore = (*aggregateStore)(nil)

type AggregateStore interface {
	Get(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) error
	Save(ctx context.Context, agg eventsourcing.Aggregate) error
}

type aggregateStore struct {
	repo repos.EventStore
}

func NewAggregateStore(repos repos.Repos) AggregateStore {
	return &aggregateStore{
		repo: repos.EventStore(),
	}
}

// Get fetches the events and build up the aggregate
func (as *aggregateStore) Get(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) error {
	if reflect.ValueOf(agg).Kind() != reflect.Ptr {
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
	root := agg.Root()
	aggType := reflect.TypeOf(agg).Elem().Name()
	root.SetAggregateType(aggType)

	err := as.repo.CreateIfNotExist(ctx, root.AggregateID(), aggType)
	if err != nil {
		return err
	}

	err = as.repo.WithTransaction(ctx, func(txn repos.EventStore) error {
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
			nthEvent := 20
			version := root.Version()
			if version%nthEvent == 0 {

				err = txn.CreateSnapshot(ctx, agg)
				if err != nil {

					return err
				}
			}
			root.Update()
		}

		return nil
	})
	if err != nil {

		return err
	}

	return nil
}

func (as *aggregateStore) getFromSnapshot(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) bool {
	has := as.repo.ReadSnapshot(ctx, aggregateID, agg.Root().Version(), agg)
	if !has {
		return false
	}

	root := agg.Root()
	if root.BaseVersion() < root.Version() {
		err := as.repo.Get(ctx, aggregateID, root.BaseVersion(), root.Version(), agg)
		if err != nil {

			return false
		}
	}

	return true
}

func (as *aggregateStore) getFromEvents(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) error {
	err := as.repo.Get(ctx, aggregateID, 0, 0, agg)
	if err != nil {
		return err
	}

	return nil
}
