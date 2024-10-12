package balancer

import (
	"sync"
	"sync/atomic"
)

func newStickyRoundRobin[T comparable](defaultType T, customizers ...Customizer) *stickyRoundRobin[T] {
	srr := &stickyRoundRobin[T]{stick: 1, defaultType: defaultType}

	for _, c := range customizers {
		c(srr)
	}

	return srr
}

type stickyRoundRobin[T comparable] struct {
	defaultType T

	values  []T
	current int64

	stick        int64
	currentStick int64

	mu sync.RWMutex
}

func (b *stickyRoundRobin[T]) Type() uint {
	return StickyRoundRobin
}

func (b *stickyRoundRobin[T]) Add(values ...T) {
	b.mu.RLock()
	b.values = append(b.values, values...)
	b.mu.RUnlock()
}

func (b *stickyRoundRobin[T]) Pick() T {
	if len(b.values) == 0 {
		return b.defaultType
	}

	val := b.values[atomic.LoadInt64(&b.current)]

	if atomic.LoadInt64(&b.stick)-1 == atomic.LoadInt64(&b.currentStick) {
		atomic.StoreInt64(&b.currentStick, 0)

		if int(atomic.LoadInt64(&b.current)) >= len(b.values)-1 {
			atomic.StoreInt64(&b.current, 0)
		} else {
			atomic.AddInt64(&b.current, 1)
		}
	} else {
		atomic.AddInt64(&b.currentStick, 1)
	}

	return val
}

func (b *stickyRoundRobin[T]) Remove(target T) bool {
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

func (b *stickyRoundRobin[T]) Amount() int {
	return len(b.values)
}

func (b *stickyRoundRobin[T]) All() []T {
	return b.values
}

func WithStick[T comparable](amount int64) Customizer {
	return func(balancer any) {
		balancer.(*stickyRoundRobin[T]).stick = amount
	}
}
