package eventsourcing

import "time"

type Event struct {
	ID            int64
	AggregateID   string
	AggregateType string
	Data          interface{}
	Metadata      interface{}
	Version       int64
	Timestamp     time.Time
}

type EventBuilder struct {
	Event *Event
}

func (e *EventBuilder) SetID(id int64) *EventBuilder {
	e.Event.ID = id
	return e
}

func (e *EventBuilder) SetAggregateID(aggregateID string) *EventBuilder {
	e.Event.AggregateID = aggregateID
	return e
}

func (e *EventBuilder) SetAggregateType(aggregateType string) *EventBuilder {
	e.Event.AggregateType = aggregateType
	return e
}

func (e *EventBuilder) SetData(data interface{}) *EventBuilder {
	e.Event.Data = data
	return e
}

func (e *EventBuilder) SetMetadata(metadata interface{}) *EventBuilder {
	e.Event.Metadata = metadata
	return e
}

func (e *EventBuilder) SetVersion(version int64) *EventBuilder {
	e.Event.Version = version
	return e
}

func (e *EventBuilder) SetTimestamp(timestamp time.Time) *EventBuilder {
	e.Event.Timestamp = timestamp
	return e
}

func (e *EventBuilder) Build() *Event {
	return e.Event
}
