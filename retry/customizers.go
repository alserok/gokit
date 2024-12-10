package retry

import "time"

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
