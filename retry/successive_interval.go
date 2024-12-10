package retry

import (
	"fmt"
	"sync/atomic"
	"time"
)

func newSuccessiveInterval(customs ...Customizer) *successiveInterval {
	si := successiveInterval{
		interval: defaultInterval,
		amount:   defaultAmount,
	}
	for _, custom := range customs {
		custom(&si)
	}

	return &si
}

type successiveInterval struct {
	interval time.Duration
	amount   uint

	prevFailures int64
}

func (s *successiveInterval) Execute(fn func() error) error {
	mult := float32(1 + uint(atomic.LoadInt64(&s.prevFailures))/s.amount)

	currFailures := 0
	defer func() {
		atomic.StoreInt64(&s.prevFailures, int64(currFailures))
	}()

	var err error
	for range s.amount {
		if err = fn(); err == nil {
			return nil
		}

		<-time.After(s.interval * time.Duration(mult))
		currFailures++
	}

	return fmt.Errorf("failed after %d retries: %w", s.amount, err)
}
