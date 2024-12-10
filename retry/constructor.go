package retry

import "time"

type Retry interface {
	Execute(func() error) error
}

const (
	FixedInterval = iota
	ExponentialBackoff
	LinearBackoff
	SuccessiveInterval

	defaultInterval = time.Millisecond * 10
	defaultAmount   = 5

	defaultStep = 0.25
)

func New(t uint, customs ...Customizer) Retry {
	switch t {
	case FixedInterval:
		return newFixedInterval(customs...)
	case ExponentialBackoff:
		return newExponentialBackoff(customs...)
	case LinearBackoff:
		return newLinearBackoff(customs...)
	case SuccessiveInterval:
		return newSuccessiveInterval(customs...)
	default:
		panic("invalid retry type")
	}
}

type Customizer func(retry any)
