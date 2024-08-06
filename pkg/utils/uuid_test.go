package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUID(t *testing.T) {
	u := NewUUID()
	s := u.String()
	assert.Equal(t, 36, len(s))
	b58 := u.B58()
	assert.Equal(t, 21, len(b58))

	u2 := UUID{}
	err := u2.FromString(s)
	assert.NoError(t, err)
	assert.Equal(t, s, u2.String())
	assert.Equal(t, b58, u2.B58())
	assert.Equal(t, u, u2)

	u3 := UUID{}
	err = u3.FromB58(b58)
	assert.NoError(t, err)
	assert.Equal(t, s, u3.String())
	assert.Equal(t, b58, u3.B58())
	assert.Equal(t, u, u3)

	us := "01912531-36cf-71ff-bf13-5a7f87bf243f"
	ub58 := "CDvWW3fmbeRpALuhCjr86"
	u4 := UUID{}
	err = u4.FromString(us)
	assert.NoError(t, err)
	assert.Equal(t, us, u4.String())
	assert.Equal(t, ub58, u4.B58())
	u5 := UUID{}
	err = u5.FromB58(ub58)
	assert.NoError(t, err)
	assert.Equal(t, us, u5.String())
	assert.Equal(t, ub58, u5.B58())

	id, err := uuid.Parse(us)
	assert.NoError(t, err)
	assert.Equal(t, time.Unix(id.Time().UnixTime()), u4.Time())
}

func BenchmarkUUIDString(b *testing.B) {
	u := NewUUID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = u.String()
	}
}

func BenchmarkFromString(b *testing.B) {
	s := NewUUID().String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := UUID{}
		_ = u.FromString(s)
	}
}

func BenchmarkB58(b *testing.B) {
	u := NewUUID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u.B58()
	}
}

func BenchmarkFromB58(b *testing.B) {
	b58 := NewUUID().B58()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := UUID{}
		_ = u.FromB58(b58)
	}
}

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
