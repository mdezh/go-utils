package fan

import "sync"

func FanIn[T any](ins ...<-chan T) <-chan T {
	return FanInBuffered(0, ins...)
}

func FanInBuffered[T any](bufLen uint, ins ...<-chan T) <-chan T {
	out := make(chan T, bufLen)

	go func() {
		wg := sync.WaitGroup{}

		wg.Add(len(ins))

		for i := range ins {
			ch := ins[i]
			go func() {
				for m := range ch {
					out <- m
				}

				wg.Done()
			}()
		}

		wg.Wait()

		close(out)
	}()

	return out
}
