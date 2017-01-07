package webfinger

import (
	"context"
	"net/http"
)

// ErrorKeyType is the type for the context error key
type ErrorKeyType int

// ErrorKey is the key for the context error
var ErrorKey ErrorKeyType

// ErrorFromContext gets the error from the context
func ErrorFromContext(ctx context.Context) error {
	v, ok := ctx.Value(ErrorKey).(error)
	if !ok {
		return nil
	}
	return v
}

func addError(r *http.Request, err error) *http.Request {
	if err == nil {
		return r
	}
	ctx := r.Context()
	ctx = context.WithValue(ctx, ErrorKey, err)
	r = r.WithContext(ctx)
	return r
}
