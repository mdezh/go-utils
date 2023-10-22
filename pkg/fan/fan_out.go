package fan

import "github.com/mdezh/go-utils/pkg/channels"

func FanOut[T any](in <-chan T, n uint) []<-chan T {
	return FanOutBuffered(0, in, n)
}

func FanOutBuffered[T any](bufLen uint, in <-chan T, n uint) []<-chan T {
	outs := make([]chan T, n)

	for i := range outs {
		outs[i] = make(chan T, bufLen)
	}

	go func() {
		var i uint
		for m := range in {
			outs[i] <- m
			i++
			i = i % n
		}

		for i := range outs {
			close(outs[i])
		}
	}()

	return channels.ToReadOnly(outs)
}
