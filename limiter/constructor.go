package limiter

import "context"

type Limiter interface {
	Allow(ctx context.Context) bool
}

const (
	LeakyBucket = iota
)

func New(t uint, customizers ...Customizer) Limiter {
	switch t {
	case LeakyBucket:
		return newLeakyBucket(customizers...)
	default:
		panic("invalid limiter type")
	}
}

type Customizer func(limiter any)
