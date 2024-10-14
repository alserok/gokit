package balancer

type Balancer[T comparable] interface {
	// Type returns balancer type
	Type() uint
	// Add adds value to balancer list
	Add(values ...T)
	// Pick retrieves value from balancer list
	Pick() T
	// Remove removes value from the balancer list
	Remove(T) bool
	// Amount returns number of values in the balancer list
	Amount() int
	// All returns balancer list
	All() []T
}

const (
	RoundRobin = iota
	StickyRoundRobin
	FastestResponse
)

func New[T comparable](t uint, defaultType T, customizer ...Customizer) Balancer[T] {
	switch t {
	case RoundRobin:
		return newRoundRobin(defaultType, customizer...)
	case StickyRoundRobin:
		return newStickyRoundRobin(defaultType, customizer...)
	case FastestResponse:
		return newFastestResponse(defaultType, customizer...)

	default:
		panic("invalid balancer type")
	}
}

type Customizer func(balancer any)
