package limiter

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestFixedWindowCounterSuite(t *testing.T) {
	suite.Run(t, new(fixedWindowCounterSuite))
}

type fixedWindowCounterSuite struct {
	suite.Suite
}

func (s *fixedWindowCounterSuite) TestDefault() {
	lim := uint(50)
	period := time.Millisecond * 500

	l := newFixedWindowCounter(WithCapacity(lim), WithTick(period))

	for range lim {
		s.Require().True(l.Allow(context.Background()))
	}

	s.Require().False(l.Allow(context.Background()))

	time.Sleep(period)

	s.Require().True(l.Allow(context.Background()))
}
