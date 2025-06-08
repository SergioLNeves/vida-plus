// Package pkg provides common utilities.
package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID generates a random unique identifier.
func GenerateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
