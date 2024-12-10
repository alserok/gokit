package limiter

import "time"

func WithCapacity(cap uint) Customizer {
	return func(limiter any) {
		switch l := limiter.(type) {
		case *leakyBucket:
			l.cap = cap
		case *fixedWindowCounter:
			l.lim = int64(cap)
		}
	}
}

func WithTick(tick time.Duration) Customizer {
	return func(limiter any) {
		switch l := limiter.(type) {
		case *leakyBucket:
			l.tick = tick
		case *fixedWindowCounter:
			l.period = tick
		}
	}
}
