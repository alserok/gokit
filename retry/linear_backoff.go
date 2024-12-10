package retry

import (
	"fmt"
	"time"
)

func newLinearBackoff(customs ...Customizer) *linearBackoff {
	lb := linearBackoff{
		amount:   defaultAmount,
		interval: defaultInterval,
		step:     defaultStep,
	}
	for _, custom := range customs {
		custom(&lb)
	}

	return &lb
}

type linearBackoff struct {
	interval time.Duration
	amount   uint
	step     float32
}

func (l *linearBackoff) Execute(fn func() error) error {
	mult := float32(1)

	var err error
	for range l.amount {
		if err = fn(); err == nil {
			return nil
		}

		<-time.After(l.interval * time.Duration(mult))
		mult += l.step
	}

	return fmt.Errorf("failed after %d retries: %w", l.amount, err)
}
