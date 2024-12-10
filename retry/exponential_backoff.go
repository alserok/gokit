package retry

import (
	"fmt"
	"time"
)

func newExponentialBackoff(customs ...Customizer) *exponentialBackoff {
	eb := exponentialBackoff{
		interval: defaultInterval,
		amount:   defaultAmount,
		exponent: 1.5,
	}
	for _, custom := range customs {
		custom(&eb)
	}

	return &eb
}

type exponentialBackoff struct {
	interval time.Duration
	amount   uint
	exponent float32
}

func (e *exponentialBackoff) Execute(f func() error) error {
	mult := float32(1)

	var err error
	for range e.amount {
		if err = f(); err == nil {
			return nil
		}

		<-time.After(e.interval * time.Duration(mult))
		mult *= e.exponent
	}

	return fmt.Errorf("failed after %d retries: %w", e.amount, err)
}
