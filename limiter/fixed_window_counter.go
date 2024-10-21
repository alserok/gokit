package limiter

import (
	"context"
	"sync/atomic"
	"time"
)

const (
	defaultLimit  = 10_000
	defaultPeriod = time.Millisecond * 500
)

func newFixedWindowCounter(customizers ...Customizer) *fixedWindowCounter {
	fwc := &fixedWindowCounter{
		lim:         defaultLimit,
		period:      defaultPeriod,
		lastResetAt: time.Now().UnixNano(),
	}

	for _, c := range customizers {
		c(fwc)
	}

	return fwc
}

type fixedWindowCounter struct {
	lim     int64
	counter int64

	period      time.Duration
	lastResetAt int64
}

func (f *fixedWindowCounter) Allow(ctx context.Context) bool {
	if f.period.Nanoseconds() < time.Now().UnixNano()-f.lastResetAt {
		atomic.StoreInt64(&f.counter, 1)
		atomic.StoreInt64(&f.lastResetAt, time.Now().UnixNano())
		return true
	}

	if atomic.LoadInt64(&f.counter) < f.lim {
		atomic.AddInt64(&f.counter, 1)
		return true
	}

	return false
}
