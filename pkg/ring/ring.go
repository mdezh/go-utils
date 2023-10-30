package ring

func Ring[T any](in <-chan T, ringLen int) <-chan T {
	if ringLen < 1 {
		panic("ring buffer length should be greater than 0")
	}

	out := make(chan T, ringLen)

	go func() {
		for v := range in {
			select {
			case out <- v:
			default:
				<-out
				out <- v
			}
		}

		close(out)
	}()

	return out
}
