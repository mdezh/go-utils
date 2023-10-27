package workerpool

import (
	"context"
	"sync"
)

type Worker[I, O any] func(context.Context, I) (O, error)

type empty struct{}

type WorkerPool[I, O any] struct {
	worker    Worker[I, O]
	workerNum int
	mx        sync.RWMutex
	started   bool
	done      chan empty
	out       chan O
	err       error
	errOnce   sync.Once
}

func New[I, O any](
	worker Worker[I, O],
	workerNum, outChanBufLen int,
) *WorkerPool[I, O] {
	return &WorkerPool[I, O]{
		worker:    worker,
		workerNum: int(workerNum),
		out:       make(chan O, outChanBufLen),
		done:      make(chan empty),
	}
}

// Returns output channel. Will be closed by worker pool after finishing of input processing.
func (p *WorkerPool[I, O]) Out() <-chan O {
	return p.out
}

// Returns channel that will be closed by worker pool after the end of work.
// Useful in case we want to know worker pool is finished, but doesn't want to read data from the output channel
func (p *WorkerPool[I, O]) Done() <-chan empty {
	return p.done
}

// Returns an error if at least one worker returned an error
func (p *WorkerPool[I, O]) Err() error {
	p.mx.RLock()
	defer p.mx.RUnlock()

	return p.err
}

func (p *WorkerPool[I, O]) Run(ctx context.Context, in <-chan I) {
	p.mx.Lock()
	defer p.mx.Unlock()

	if p.started {
		return
	}
	p.started = true

	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(p.workerNum)

		for i := 0; i < p.workerNum; i++ {
			go func() {
				defer wg.Done()

				for {
					select {
					case <-ctx.Done():
						return
					default:
						select {
						case <-ctx.Done():
							return
						case input, ok := <-in:
							if !ok {
								return
							}
							res, err := p.worker(ctx, input)
							if err != nil {
								p.errOnce.Do(func() {
									p.setErr(err)
									cancel()
								})
							} else {
								p.out <- res
							}
						}
					}
				}
			}()
		}

		wg.Wait()

		close(p.out)
		close(p.done)
	}()
}

func (p *WorkerPool[I, O]) setErr(err error) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.err = err
}
