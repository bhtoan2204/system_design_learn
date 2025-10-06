package repository

import (
	"event_sourcing_bank_system_api/package/eventsourcing"

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

func (r *repos) EventStore() EventStore {
	return r.ev
}
