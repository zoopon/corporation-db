package domain

import (
	"crypto/rand"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
)

// NewUUIDv7 generates a new UUIDv7 with timestamp-based ordering.
// UUIDv7 provides better database performance due to sequential ordering.
func NewUUIDv7() uuid.UUID {
	// Get current timestamp in milliseconds since Unix epoch
	now := time.Now().UnixMilli()

	// Create 16-byte array for UUID
	var uuidBytes [16]byte

	// Set the 48-bit timestamp (6 bytes) in big-endian format
	binary.BigEndian.PutUint64(uuidBytes[0:8], uint64(now))
	// Shift to get only the lower 48 bits (6 bytes)
	copy(uuidBytes[0:6], uuidBytes[2:8])

	// Fill remaining 10 bytes with random data
	rand.Read(uuidBytes[6:16])

	// Set version (7) in the most significant 4 bits of byte 6
	uuidBytes[6] = (uuidBytes[6] & 0x0f) | 0x70

	// Set variant (10) in the most significant 2 bits of byte 8
	uuidBytes[8] = (uuidBytes[8] & 0x3f) | 0x80

	return uuid.UUID(uuidBytes)
}

// MustNewUUIDv7 generates a new UUIDv7 and panics on error.
// Use this when you're certain UUID generation will succeed.
func MustNewUUIDv7() uuid.UUID {
	return NewUUIDv7()
}
