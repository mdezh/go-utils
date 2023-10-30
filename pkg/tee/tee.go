package tee

import (
	"github.com/mdezh/go-utils/pkg/channels"
)

type empty struct{}

// 'Out' channels will be closed immediately after closing 'in' channel, so blocked writes will be lost
func Tee[T any](in <-chan T, numOuts int) []<-chan T {
	return TeeBuffered(in, numOuts, 0)
}

// 'Out' channels will be closed immediately after closing 'in' channel, so blocked writes will be lost
func TeeBuffered[T any](in <-chan T, numOuts, bufLen int) []<-chan T {
	outs := make([]chan T, numOuts)
	for i := 0; i <= numOuts; i++ {
		outs[i] = make(chan T, bufLen)
	}

	go func() {
		done := make(chan empty)

		for v := range in {
			for i := range outs {
				v, i := v, i
				select {
				case outs[i] <- v:
				default:
					go func() {
						select {
						case <-done:
						case outs[i] <- v:
						}
					}()
				}
			}
		}

		close(done)

		for i := range outs {
			close(outs[i])
		}
	}()

	return channels.ToReadOnly(outs)
}
