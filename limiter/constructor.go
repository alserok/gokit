package limiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context) bool
}

const (
	LeakyBucket = iota
	FixedWindowCounter
)

func New(t uint, customizers ...Customizer) Limiter {
	switch t {
	case LeakyBucket:
		return newLeakyBucket(customizers...)
	case FixedWindowCounter:
		return newFixedWindowCounter()
	default:
		panic("invalid limiter type")
	}
}

type Customizer func(limiter any)

func WithCapacity(cap uint) Customizer {
	return func(limiter any) {
		switch limiter.(type) {
		case *leakyBucket:
			limiter.(*leakyBucket).cap = cap
		case *fixedWindowCounter:
			limiter.(*fixedWindowCounter).lim = int64(cap)
		}
	}
}

func WithTick(tick time.Duration) Customizer {
	return func(limiter any) {
		switch limiter.(type) {
		case *leakyBucket:
			limiter.(*leakyBucket).tick = tick
		case *fixedWindowCounter:
			limiter.(*fixedWindowCounter).period = tick
		}
	}
}
