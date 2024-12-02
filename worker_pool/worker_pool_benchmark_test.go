package worker_pool

import (
	"context"
	"runtime"
	"testing"
)

func BenchmarkNewWorkerPoolWith100Workers(b *testing.B) {
	p := newWorkerPool(func(_ any) error { return nil }, 100)
	p.Start()
	defer p.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Add(context.Background(), nil)
	}
}

func BenchmarkNewWorkerPoolWith10Workers(b *testing.B) {
	p := newWorkerPool(func(_ any) error { return nil }, 10)
	p.Start()
	defer p.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Add(context.Background(), nil)
	}
}

func BenchmarkNewWorkerPoolWith1Worker(b *testing.B) {
	p := newWorkerPool(func(_ any) error { return nil }, 1)
	p.Start()
	defer p.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Add(context.Background(), nil)
	}
}

func BenchmarkNewWorkerPoolWithNumCPUWorkers(b *testing.B) {
	p := newWorkerPool(func(_ any) error { return nil }, int64(runtime.NumCPU()))
	p.Start()
	defer p.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Add(context.Background(), nil)
	}
}
