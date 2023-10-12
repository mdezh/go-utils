package channels

func ToWriteOnly[T any](chans []chan T) []chan<- T {
	res := make([]chan<- T, len(chans))
	for i := range chans {
		res[i] = chans[i]
	}

	return res
}
