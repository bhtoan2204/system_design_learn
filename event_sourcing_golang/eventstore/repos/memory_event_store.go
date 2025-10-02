package repos

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"event_sourcing_golang/pkg/eventsourcing"
)

// memEventStore is an in-memory implementation of EventStore for testing/demo purposes.
type memEventStore struct {
	mu         sync.RWMutex
	serializer eventsourcing.Serializer

	// aggregates keeps latest persisted version per aggregate
	aggregates map[string]*memAggregate
	// events keeps ordered events per aggregate id
	events map[string][]eventsourcing.Event
	// snapshots keeps the latest snapshot per aggregate id (by version)
	snapshots map[string]memSnapshot
}

type memAggregate struct {
	Version       int
	AggregateType string
}

type memSnapshot struct {
	Version int
	Data    []byte
	// AggregateVersion at snapshot time
	AggregateVersion int
}

func newMemoryEventStore(s eventsourcing.Serializer) EventStore {
	return &memEventStore{
		serializer: s,
		aggregates: make(map[string]*memAggregate),
		events:     make(map[string][]eventsourcing.Event),
		snapshots:  make(map[string]memSnapshot),
	}
}

func (m *memEventStore) WithTransaction(ctx context.Context, fn func(EventStore) error) (err error) {
	// In-memory variant runs the fn under the same mutex; emulate transactional behavior.
	m.mu.Lock()
	defer m.mu.Unlock()
	return fn(m)
}

func (m *memEventStore) CreateIfNotExist(ctx context.Context, id, typ string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.aggregates[id]; !ok {
		m.aggregates[id] = &memAggregate{Version: 0, AggregateType: typ}
	}
	return nil
}

func (m *memEventStore) CheckAndUpdateVersion(ctx context.Context, agg eventsourcing.Aggregate) bool {
	root := agg.Root()
	id := root.AggregateID()
	expected := root.BaseVersion()
	newVersion := root.Version()

	a, ok := m.aggregates[id]
	if !ok {
		return false
	}
	if a.Version != expected {
		return false
	}
	a.Version = newVersion
	return true
}

func (m *memEventStore) Append(ctx context.Context, e eventsourcing.Event) error {
	// Ensure sequence
	list := m.events[e.AggregateID]
	if len(list) > 0 {
		last := list[len(list)-1]
		if e.Version != last.Version+1 {
			return fmt.Errorf("version gap for aggregate %s, got=%d want=%d", e.AggregateID, e.Version, last.Version+1)
		}
	} else {
		if e.Version != 1 {
			return fmt.Errorf("first event must have version=1, got=%d", e.Version)
		}
	}
	m.events[e.AggregateID] = append(list, e)
	return nil
}

func (m *memEventStore) Get(ctx context.Context, aggregateID string, fromVersion, toVersion int, agg eventsourcing.Aggregate) error {
	root := agg.Root()
	events := m.events[aggregateID]
	for _, evt := range events {
		if (fromVersion == 0 || evt.Version > fromVersion) && (toVersion == 0 || evt.Version <= toVersion) {
			f, ok := m.serializer.Type(root.AggregateType(), evt.EventType)
			if !ok {
				return fmt.Errorf("cant serialize event with type: %s_%s", root.AggregateType(), evt.EventType)
			}
			eventData := f()
			// Convert evt.Data to correct concrete type via marshal/unmarshal to be safe.
			b, err := m.serializer.Marshal(evt.Data)
			if err != nil {
				return err
			}
			if err := m.serializer.Unmarshal(b, &eventData); err != nil {
				return err
			}
			evt.Data = eventData
			root.LoadFromHistory(agg, []eventsourcing.Event{evt})
		}
	}
	return nil
}

func (m *memEventStore) CreateSnapshot(ctx context.Context, agg eventsourcing.Aggregate) error {
	root := agg.Root()
	data, err := m.serializer.Marshal(agg)
	if err != nil {
		return err
	}
	m.snapshots[root.AggregateID()] = memSnapshot{Version: root.Version(), Data: data, AggregateVersion: root.Version()}
	return nil
}

func (m *memEventStore) ReadSnapshot(ctx context.Context, aggregateID string, version int, agg eventsourcing.Aggregate) bool {
	snap, ok := m.snapshots[aggregateID]
	if !ok {
		return false
	}
	if snap.Version < version {
		return false
	}
	if err := m.serializer.Unmarshal(snap.Data, agg); err != nil {
		return false
	}
	agg.Root().SetInternal(aggregateID, snap.Version, snap.AggregateVersion)
	return true
}

func (m *memEventStore) List(ctx context.Context, aggregateID string, agg eventsourcing.Aggregate) ([]eventsourcing.Event, error) {
	// derive aggregate type via reflection to ensure it's set
	aggType := reflect.TypeOf(agg).Elem().Name()
	list := m.events[aggregateID]
	res := make([]eventsourcing.Event, 0, len(list))
	for _, evt := range list {
		f, ok := m.serializer.Type(aggType, evt.EventType)
		if !ok {
			return nil, fmt.Errorf("cant serialize event with type: %s_%s", aggType, evt.EventType)
		}
		eventData := f()
		b, err := m.serializer.Marshal(evt.Data)
		if err != nil {
			return nil, err
		}
		if err := m.serializer.Unmarshal(b, &eventData); err != nil {
			return nil, err
		}
		evt.Data = eventData
		res = append(res, evt)
	}
	return res, nil
}

func (m *memEventStore) SnapshotVersion(ctx context.Context, aggregateID string) (int, bool) {
	snap, ok := m.snapshots[aggregateID]
	if !ok {
		return 0, false
	}
	return snap.Version, true
}
