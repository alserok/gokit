package balancer

import (
	"context"
	"sync"
	"time"
)

func newFastestResponseTime[T comparable](defaultType T, customizers ...Customizer) *fastestResponseTime[T] {
	frr := &fastestResponseTime[T]{
		defaultType: defaultType,
		values:      make(map[T]struct{}),
	}

	for _, c := range customizers {
		c(frr)
	}

	return frr
}

type fastestResponseTime[T comparable] struct {
	defaultType T

	q          chan T
	values     map[T]struct{}
	currentLen int64

	mu sync.RWMutex
}

func (f *fastestResponseTime[T]) Type() uint {
	return FastestResponseTime
}

func (f *fastestResponseTime[T]) Add(values ...T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	f.mu.RLock()
	for _, val := range values {
		f.values[val] = struct{}{}

		select {
		case f.q <- val:
		case <-ctx.Done():
		}
	}
	f.mu.RUnlock()
}

func (f *fastestResponseTime[T]) Pick() T {
	if len(f.values) == 0 {
		return f.defaultType
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	select {
	case val := <-f.q:
		return val
	case <-ctx.Done():
		for key, _ := range f.values {
			return key
		}
	}

	return f.defaultType
}

func (f *fastestResponseTime[T]) Remove(target T) bool {
	f.mu.RLock()
	delete(f.values, target)
	f.mu.RUnlock()

	return false
}

func (f *fastestResponseTime[T]) Amount() int {
	return len(f.values)
}

func (f *fastestResponseTime[T]) All() []T {
	res := make([]T, 0, len(f.values))

	for val, _ := range f.values {
		res = append(res, val)
	}

	return res
}

func WithUpdater[T comparable](ch chan T) Customizer {
	return func(balancer any) {
		balancer.(*fastestResponseTime[T]).q = ch
	}
}
