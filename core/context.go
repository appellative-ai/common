package core

import (
	"context"
	"time"
)

var (
	cancelFn = func() {}
)

func NewContext(ctx context.Context, timeout time.Duration) (context.Context, func()) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		return ctx, cancelFn
	}
	return context.WithTimeout(ctx, timeout)
}
