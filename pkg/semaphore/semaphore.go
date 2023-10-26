package semaphore

type empty struct{}

type Semaphore struct {
	c chan empty
}

func New(maxInParallel int) *Semaphore {
	return &Semaphore{
		c: make(chan empty, maxInParallel),
	}
}

func (s *Semaphore) Acquire() {
	s.c <- empty{}
}

func (s *Semaphore) Release() {
	<-s.c
}
