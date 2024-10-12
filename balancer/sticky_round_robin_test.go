package balancer

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestStickyRoundRobinSuite(t *testing.T) {
	suite.Run(t, new(StickyRoundRobinSuite))
}

type StickyRoundRobinSuite struct {
	suite.Suite
}

// Each value should be iterated 'cycles' times
func (r *StickyRoundRobinSuite) TestDefault() {
	cycles := 3
	testValues := []string{"a", "b", "c", "d", "e"}
	stickAmount := 3

	b := newStickyRoundRobin("", WithStick[string](int64(stickAmount)))

	for _, v := range testValues {
		b.Add(v)
	}

	r.Require().Equal(len(testValues), b.Amount())

	for i := 0; i < cycles; i++ {
		for j := range b.Amount() {
			for range stickAmount {
				r.Require().Equal(testValues[j], b.Pick(), fmt.Sprintf("cycle: %d idx: %d", i, j))
			}
		}
	}
}

func (r *StickyRoundRobinSuite) TestEmpty() {
	cycles := 1
	defaultType := ""
	b := newStickyRoundRobin(defaultType)

	r.Require().Equal(0, b.Amount())

	for i := 0; i < cycles; i++ {
		r.Require().Equal(defaultType, b.Pick())
	}
}
