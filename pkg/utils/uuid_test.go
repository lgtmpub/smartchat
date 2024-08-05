package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUIDv4(t *testing.T) {
	assert.Equal(t, 36, len(UUIDv4().String()))
}

func BenchmarkNewUUIDv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUIDv4()
	}
}

func TestNewUUIDv7(t *testing.T) {
	assert.Equal(t, 36, len(UUIDv7().String()))
}

func BenchmarkNewUUIDv7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UUIDv7()
	}
}

func TestRandBytes(t *testing.T) {
	bytes := RandBytes(32)
	assert.Equal(t, 32, len(bytes))
}

func BenchmarkRandByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandBytes(32)
	}
}
