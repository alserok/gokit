package worker_pool

import (
	"context"
	"testing"
)

func BenchmarkNewWorkerPoolWith100Workers(b *testing.B) {
	p := newWorkerPool(func(_ any) error { return nil }, 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Add(context.Background(), nil)
	}

	for p.currWorkers != 0 {
	}
}
