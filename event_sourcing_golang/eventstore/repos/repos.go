package repos

import (
	"event_sourcing_golang/pkg/eventsourcing"

	"gorm.io/gorm"
)

var _ Repos = (*repos)(nil)

type Repos interface {
	EventStore() EventStore
}

type repos struct {
	ev EventStore
}

func New(db *gorm.DB, s eventsourcing.Serializer) Repos {
	ev := newEventStore(db, s)

	return &repos{
		ev: ev,
	}
}

// NewInMemory creates a Repos backed by in-memory event store
func NewInMemory(s eventsourcing.Serializer) Repos {
	ev := newMemoryEventStore(s)

	return &repos{
		ev: ev,
	}
}

func (r *repos) EventStore() EventStore {
	return r.ev
}
