package ctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lgtmpub/smartchat/pkg/utils"
)

func TestRequestID(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "", GetRequestID(ctx))

	reqID := utils.UUIDv7().String()
	ctx = context.WithValue(context.Background(), RequestIDKey, reqID)
	assert.Equal(t, reqID, GetRequestID(ctx))
	assert.Equal(t, GetRequestID(ctx), GetRequestID(ctx))
	ctx = WithRequestID(ctx)
	assert.Equal(t, reqID, GetRequestID(ctx))

	ctx = context.WithValue(context.Background(), RequestIDKey, "123")
	reqID = GetRequestID(ctx)
	assert.Equal(t, "123", reqID)
	ctx = WithRequestID(ctx)
	assert.Equal(t, "123", GetRequestID(ctx))
}
