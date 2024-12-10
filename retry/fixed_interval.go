package retry

import (
	"fmt"
	"time"
)

func newFixedInterval(customs ...Customizer) *fixedInterval {
	fi := fixedInterval{
		interval: defaultInterval,
		amount:   defaultAmount,
	}
	for _, custom := range customs {
		custom(&fi)
	}

	return &fi
}

type fixedInterval struct {
	interval time.Duration
	amount   uint
}

func (h *fixedInterval) Execute(fn func() error) error {
	var err error
	for range h.amount {
		if err = fn(); err == nil {
			return nil
		}

		<-time.After(h.interval)
	}

	return fmt.Errorf("failed after %d retries: %w", h.amount, err)
}
