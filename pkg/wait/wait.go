package wait

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

var errWrongDuration = errors.New(
	"failed to wait: min time duration greater than max time duration",
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
		return errWrongDuration
	}

	ms := minD.Milliseconds() + int64(rand.Intn(int(maxD.Milliseconds()-minD.Milliseconds())+1))

	return Wait(ctx, time.Duration(ms*int64(time.Millisecond)))
}

func WaitRandMs(ctx context.Context, minMs, maxMs int) error {
	if minMs > maxMs {
		return errWrongDuration
	}

	ms := minMs + rand.Intn(maxMs-minMs+1)

	return Wait(ctx, time.Millisecond*time.Duration(ms))
}
