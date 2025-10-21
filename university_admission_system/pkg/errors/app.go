package errors

import "errors"

var (
	// ErrNotFound indicates an entity was not located.
	ErrNotFound = errors.New("not found")
	// ErrConflict indicates a state conflict (e.g., invalid transition).
	ErrConflict = errors.New("conflict")
	// ErrInvalidInput indicates the request failed validation.
	ErrInvalidInput = errors.New("invalid input")
	// ErrInternal indicates an unexpected failure.
	ErrInternal = errors.New("internal error")
)
