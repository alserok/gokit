package balancer

import (
	"sync"
	"sync/atomic"
)

func newRoundRobin[T comparable](defaultType T, customizers ...Customizer) *roundRobin[T] {
	return &roundRobin[T]{defaultType: defaultType}
}

type roundRobin[T comparable] struct {
	defaultType T

	values  []T
	current int64

	mu sync.RWMutex
}

func (b *roundRobin[T]) Type() uint {
	return RoundRobin
}

func (b *roundRobin[T]) Add(values ...T) {
	b.mu.RLock()
	b.values = append(b.values, values...)
	b.mu.RUnlock()
}

func (b *roundRobin[T]) Pick() T {
	if len(b.values) == 0 {
		return b.defaultType
	}

	val := b.values[atomic.LoadInt64(&b.current)]

	if int(atomic.LoadInt64(&b.current)) >= len(b.values)-1 {
		atomic.StoreInt64(&b.current, 0)
	} else {
		atomic.AddInt64(&b.current, 1)
	}

	return val
}

func (b *roundRobin[T]) Remove(target T) bool {
	b.mu.RLock()
	for idx, val := range b.values {
		if val == target {
			b.values = append(b.values[:idx], b.values[idx+1:]...)
			return true
		}
	}
	b.mu.RUnlock()

	return false
}

func (b *roundRobin[T]) Amount() int {
	return len(b.values)
}

func (b *roundRobin[T]) All() []T {
	return b.values
}
