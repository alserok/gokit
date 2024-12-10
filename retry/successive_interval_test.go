package retry

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestSuccessiveInterval(t *testing.T) {
	suite.Run(t, new(successiveIntervalSuite))
}

type successiveIntervalSuite struct {
	suite.Suite

	retry *successiveInterval
}

func (s *successiveIntervalSuite) SetupTest() {
	s.retry = newSuccessiveInterval()
}

func (s *successiveIntervalSuite) TestDefault() {
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

func (s *successiveIntervalSuite) TestWithoutErrors() {
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
