package breaker

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestBreakerSuite(t *testing.T) {
	suite.Run(t, new(BreakerSuite))
}

type BreakerSuite struct {
	suite.Suite
}

func (s *BreakerSuite) TestDefault() {
	failToClose := 3
	b := newBreaker(time.Millisecond*10, int64(failToClose))

	fn := func() bool {
		return false // No error
	}
	errFn := func() bool {
		return true // Simulating a server error
	}

	// ok, default call
	s.Require().True(b.Execute(fn))

	// failToClose errors happen
	for i := 0; i < failToClose; i++ {
		s.Require().Equal(open, int(b.status))
		s.Require().True(b.Execute(errFn))
	}

	// call is canceled
	s.Require().False(b.Execute(fn))
	s.Require().Equal(closed, int(b.status))

	time.Sleep(time.Millisecond * 20)

	// check if openClosed => closed
	s.Require().True(b.Execute(errFn))
	s.Require().Equal(closed, int(b.status))

	time.Sleep(time.Millisecond * 20)

	// check if openClosed = > open
	s.Require().True(b.Execute(fn))
	s.Require().Equal(open, int(b.status))
}
