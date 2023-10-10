package wait

import (
	"context"
	"time"
)

func Wait(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	return WaitCh(ctx, timer.C)
}

func WaitTo(ctx context.Context, t time.Time) error {
	timer := time.NewTimer(time.Until(t))
	defer timer.Stop()

	return WaitCh(ctx, timer.C)
}

func WaitCh[T any](ctx context.Context, ch <-chan T) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}
