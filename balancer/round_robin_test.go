package balancer

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestRoundRobinSuite(t *testing.T) {
	suite.Run(t, new(RoundRobinSuite))
}

type RoundRobinSuite struct {
	suite.Suite
}

// Each value should be iterated 'cycles' times
func (r *RoundRobinSuite) TestDefault() {
	cycles := 2
	testValues := []string{"a", "b", "c", "d"}

	b := newRoundRobin("")

	for _, v := range testValues {
		b.Add(v)
	}

	r.Require().Equal(len(testValues), b.Amount())

	for i := 0; i < cycles; i++ {
		for j := range b.Amount() {
			r.Require().Equal(testValues[j], b.Pick(), fmt.Sprintf("cycle: %d idx: %d", i, j))
		}
	}
}

func (r *RoundRobinSuite) TestEmpty() {
	cycles := 1
	defaultType := ""
	b := newRoundRobin(defaultType)

	r.Require().Equal(0, b.Amount())

	for i := 0; i < cycles; i++ {
		r.Require().Equal(defaultType, b.Pick())
	}
}
