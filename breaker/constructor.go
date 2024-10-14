package breaker

import "time"

type Breaker interface {
	// Execute runs function, function should return true if it is a server error
	Execute(func() bool) bool
	// Executable returns true if it is no timeout
	Executable() bool
}

const (
	closed = iota
	openClosed
	open
)

func New(timeout time.Duration, failToClose int64) Breaker {
	return newBreaker(timeout, failToClose)
}
