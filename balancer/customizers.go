package balancer

func WithUpdater[T comparable](ch chan T) Customizer {
	return func(balancer any) {
		switch b := balancer.(type) {
		case *fastestResponseTime[T]:
			b.q = ch
		default:
			panic("invalid balancer type")
		}
	}
}

func WithStick[T comparable](amount int64) Customizer {
	return func(balancer any) {
		switch b := balancer.(type) {
		case *stickyRoundRobin[T]:
			b.stick = amount
		default:
			panic("invalid balancer type")
		}
	}
}
