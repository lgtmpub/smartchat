package ctx

import (
	"context"

	"github.com/spf13/cast"

	"github.com/lgtmpub/smartchat/pkg/utils"
)

// Key is the type used to store the request ID in the context.
type Key string

// RequestIDKey is the key used to store the request ID in the context.
var RequestIDKey Key = "X-Request-ID"

// GetRequestID returns the request ID from the context.
func GetRequestID(ctx context.Context) string {
	reqID := cast.ToString(ctx.Value(RequestIDKey))
	if reqID == "" {
		return ""
	}
	return reqID
}

// WithRequestID sets the request ID in the context.
func WithRequestID(ctx context.Context) context.Context {
	val := ctx.Value(RequestIDKey)
	if val != nil {
		return ctx
	}
	return context.WithValue(ctx, RequestIDKey, utils.UUIDv7())
}
