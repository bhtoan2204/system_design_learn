package eventsourcing

import "time"

type AggregateRoot struct {
	ID                string
	Version           int64
	Events            []Event
	Snapshot          interface{}
	SnapshotVersion   int64
	SnapshotTimestamp time.Time
}
