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

func WithAmount(amount uint) Customizer {
	return func(retry any) {
		switch r := retry.(type) {
		case *exponentialBackoff:
			r.amount = amount
		case *linearBackoff:
			r.amount = amount
		case *fixedInterval:
			r.amount = amount
		case *successiveInterval:
			r.amount = amount
		default:
			panic("invalid retry type")
		}
	}
}

func WithInterval(interval time.Duration) Customizer {
	return func(retry any) {
		switch r := retry.(type) {
		case *exponentialBackoff:
			r.interval = interval
		case *linearBackoff:
			r.interval = interval
		case *fixedInterval:
			r.interval = interval
		case *successiveInterval:
			r.interval = interval
		default:
			panic("invalid retry type")
		}
	}
}

func WithExponent(exp float32) Customizer {
	return func(retry any) {
		switch r := retry.(type) {
		case *exponentialBackoff:
			r.exponent = exp
		default:
			panic("invalid retry type")
		}
	}
}

func WithStep(step float32) Customizer {
	return func(retry any) {
		switch r := retry.(type) {
		case *linearBackoff:
			r.step = step
		default:
			panic("invalid retry type")
		}
	}
}
