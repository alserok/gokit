package limiter

import (
	"context"
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
