package domain

// IDGenerator abstracts ID creation for aggregates.
type IDGenerator interface {
	NewID() string
}
