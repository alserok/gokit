package breaker

import (
	"sync/atomic"
	"time"
)

func newBreaker(timeout time.Duration, failToClose int64) *breaker {
	return &breaker{timeout: timeout, failToClose: failToClose, status: open, closeTime: time.Now().UnixNano()}
}

type breaker struct {
	status int64

	timeout   time.Duration
	closeTime int64

	failToClose int64
	failed      int64
}

func (b *breaker) Execute(fn func() bool) bool {
	if time.Now().UnixNano()-atomic.LoadInt64(&b.closeTime) >= b.timeout.Nanoseconds() {
		if atomic.LoadInt64(&b.status) == closed {
			atomic.StoreInt64(&b.status, openClosed)
		} else {
			atomic.StoreInt64(&b.failed, 0)
		}
	}

	if atomic.LoadInt64(&b.status) == closed {
		return false
	}

	isServerError := fn()

	if isServerError {
		if atomic.LoadInt64(&b.status) == openClosed {
			atomic.StoreInt64(&b.status, closed)
			atomic.StoreInt64(&b.closeTime, time.Now().UnixNano())
		} else {
			atomic.AddInt64(&b.failed, 1)
			if atomic.LoadInt64(&b.failed) == b.failToClose {
				atomic.StoreInt64(&b.status, closed)
				atomic.StoreInt64(&b.failed, 0)
				atomic.StoreInt64(&b.closeTime, time.Now().UnixNano())
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
