package breaker

import (
	"sync/atomic"
	"time"
)

func newBreaker(timeout time.Duration, failToClose int64) *breaker {
	return &breaker{timeout: timeout, failToClose: failToClose, status: open}
}

type breaker struct {
	status int64

	timeout   time.Duration
	closeTime time.Time

	failToClose int64
	failed      int64
}

func (b *breaker) Execute(fn func() bool) bool {
	if atomic.LoadInt64(&b.status) == closed {
		if time.Since(b.closeTime) >= b.timeout {
			atomic.StoreInt64(&b.status, openClosed)
		} else {
			return false
		}
	}

	isServerError := fn()

	if isServerError {
		if atomic.LoadInt64(&b.status) == openClosed {
			atomic.StoreInt64(&b.status, closed)
			b.closeTime = time.Now()
		} else {
			atomic.AddInt64(&b.failed, 1)
			if atomic.LoadInt64(&b.failed) == b.failToClose {
				atomic.StoreInt64(&b.status, closed)
				atomic.StoreInt64(&b.failed, 0)
				b.closeTime = time.Now()
			}
		}
	} else {
		if atomic.LoadInt64(&b.status) == openClosed {
			atomic.StoreInt64(&b.status, open)
		}
	}

	return true
}

func (b *breaker) Executable() bool {
	return atomic.LoadInt64(&b.status) == open
}
