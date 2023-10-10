package wait

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

func Wait(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}

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

func WaitRand(ctx context.Context, minD, maxD time.Duration) error {
	if minD > maxD {
		return errors.New("failed to wait: min time duration greater than max time duration")
	}

	ms := minD.Milliseconds() + int64(rand.Intn(int(maxD.Milliseconds()-minD.Milliseconds())+1))

	return Wait(ctx, time.Duration(ms*int64(time.Millisecond)))
}
