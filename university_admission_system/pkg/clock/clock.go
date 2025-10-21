package clock

import "time"

// Clock abstracts time retrieval for deterministic tests.
type Clock interface {
	Now() time.Time
}

// SystemClock implements Clock using time.Now.
type SystemClock struct{}

// Now returns the current UTC time.
func (SystemClock) Now() time.Time {
	return time.Now().UTC()
}
