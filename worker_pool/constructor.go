package worker_pool

import "context"

type Executable[T any] func(T) error

type WorkerPool[T any] interface {
	// Add adds value to input channel for workers
	Add(ctx context.Context, val T) bool

	// Start starts workers
	Start()
	// Shutdown immediately stops workers
	Shutdown()
	// Stop stops workers after they finish all data processing
	Stop()

	// SetWorkers sets max amount of workers
	SetWorkers(workers uint)

	// Errors returns errors if any
	Errors() <-chan error
}

func NewWorkerPool[T any](executable Executable[T], maxWorkers int64) WorkerPool[T] {
	return newWorkerPool(executable, maxWorkers)
}
