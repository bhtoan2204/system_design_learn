package id

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomGenerator creates random hexadecimal identifiers.
type RandomGenerator struct{}

// NewID generates a new identifier.
func (RandomGenerator) NewID() string {
	const size = 16
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
