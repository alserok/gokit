package worker_pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

func newWorkerPool[T any](exec Executable[T], maxWorkers int64) *workerPool[T] {
	return &workerPool[T]{
		exec:          exec,
		maxWorkers:    maxWorkers,
		input:         make(chan T, maxWorkers),
		errors:        make(chan error, maxWorkers),
		workers:       make(chan struct{}, maxWorkers),
		workersUpdate: make(chan struct{}, 1),
		shutdownCh:    make(chan struct{}),
	}
}

type workerPool[T any] struct {
	exec          Executable[T]
	input         chan T
	errors        chan error
	workersUpdate chan struct{}

	workers     chan struct{}
	maxWorkers  int64
	currWorkers int64

	shutdownCh chan struct{}
}

func (w *workerPool[T]) Shutdown() {
	close(w.shutdownCh)

	for atomic.LoadInt64(&w.currWorkers) != 0 {
	}
}

func (w *workerPool[T]) Stop() {
	for len(w.input) != 0 {
	}

	close(w.shutdownCh)

	for atomic.LoadInt64(&w.currWorkers) != 0 {
	}
}

func (w *workerPool[T]) Errors() <-chan error {
	return w.errors
}

func (w *workerPool[T]) SetWorkers(amount uint) {
	atomic.StoreInt64(&w.maxWorkers, int64(amount))
	select {
	case w.workersUpdate <- struct{}{}:
	default:
	}
}

func (w *workerPool[T]) Start() {
	go func() {
		wg := &sync.WaitGroup{}
		for range int(atomic.LoadInt64(&w.maxWorkers)) {
			select {
			case w.workers <- struct{}{}:
			default:
				return
			}
		}

		for range w.workers {
			wg.Add(1)
			atomic.AddInt64(&w.currWorkers, 1)

			go func() {
				defer wg.Done()
				defer atomic.AddInt64(&w.currWorkers, -1)

				for {
					select {
					case val, ok := <-w.input:
						if !ok {
							return
						}

						if err := w.exec(val); err != nil {
							select {
							case w.errors <- err:
							default:
							}
						}
					case _, ok := <-w.workersUpdate:
						if !ok {
							return
						}

						if atomic.LoadInt64(&w.currWorkers) >= atomic.LoadInt64(&w.maxWorkers) {
							w.workersUpdate <- struct{}{}
							return
						}

						for atomic.LoadInt64(&w.currWorkers) < atomic.LoadInt64(&w.maxWorkers) {
							w.workers <- struct{}{}
							time.Sleep(time.Microsecond * 10)
						}
					case <-w.shutdownCh:
						return
					}
				}
			}()
		}

		wg.Wait()

		close(w.input)
		close(w.errors)
		close(w.workers)
		close(w.workersUpdate)
	}()
}

func (w *workerPool[T]) Add(ctx context.Context, val T) bool {
	select {
	case w.input <- val:
		return true
	case <-ctx.Done():
		return false
	default:
		return false
	}
}
