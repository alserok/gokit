package breaker

type Breaker interface {
	// Execute runs function, function should return true if it is a server error
	Execute(func() bool)
	// Executable returns true if it is no timeout
	Executable() bool
}

const (
	GRPC = iota
	HTTP
)

func New(t uint) Breaker {
	switch t {
	case GRPC:
		return nil
	case HTTP:
		return nil
	default:
		panic("invalid breaker type")
	}
}
