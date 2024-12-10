package retry

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestLinearBackoff(t *testing.T) {
	suite.Run(t, new(linearBackoffSuite))
}

type linearBackoffSuite struct {
	suite.Suite

	retry *linearBackoff
}

func (s *linearBackoffSuite) SetupTest() {
	s.retry = newLinearBackoff()
}

func (s *linearBackoffSuite) TestDefault() {
	s.retry.amount = uint(3)
	s.retry.interval = time.Millisecond * 3

	executed := 0
	start := time.Now()

	s.Require().Error(s.retry.Execute(func() error {
		executed++
		return errors.New("error")
	}))

	s.Require().Greater(time.Since(start), s.retry.interval*time.Duration(s.retry.amount))
	s.Require().Equal(s.retry.amount, uint(executed))
}

func (s *linearBackoffSuite) TestWithoutErrors() {
	s.retry.amount = uint(3)
	s.retry.interval = time.Millisecond * 3

	executed := 0
	start := time.Now()

	s.Require().NoError(s.retry.Execute(func() error {
		executed++
		return nil
	}))

	s.Require().Less(time.Since(start), s.retry.interval*time.Duration(s.retry.amount))
	s.Require().Equal(1, executed)
}
