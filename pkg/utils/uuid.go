package utils

import (
	"crypto/rand"
	"io"

	"github.com/google/uuid"
)

const MaxTry = 8

// UUIDv4 returns a new UUIDv4.
func UUIDv4() uuid.UUID {
	for i := 0; i < MaxTry; i++ {
		id, err := uuid.NewRandom()
		if err == nil {
			return id
		}
	}
	panic("failed to generate UUIDv4")
}

// UUIDv7 returns a new UUIDv7.
func UUIDv7() uuid.UUID {
	for i := 0; i < MaxTry; i++ {
		id, err := uuid.NewV7()
		if err == nil {
			return id
		}
	}
	panic("failed to generate UUIDv7")
}

// RandBytes returns a random byte slice.
func RandBytes(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < MaxTry; i++ {
		if _, err := io.ReadFull(rand.Reader, b); err == nil {
			return b
		}
	}
	panic("failed to generate random string")
}
