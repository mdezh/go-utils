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
	defer func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}()

	return WaitCh(ctx, timer.C)
}

func WaitTo(ctx context.Context, t time.Time) error {
	return Wait(ctx, time.Until(t))
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

	ns := minD.Nanoseconds() + int64(rand.Intn(int(maxD.Nanoseconds()-minD.Nanoseconds())+1))

	return Wait(ctx, time.Duration(ns))
}

func WaitRandMs(ctx context.Context, minMs, maxMs int) error {
	if minMs > maxMs {
		return errWrongDuration
	}

	ms := minMs + rand.Intn(maxMs-minMs+1)

	return Wait(ctx, time.Millisecond*time.Duration(ms))
}
