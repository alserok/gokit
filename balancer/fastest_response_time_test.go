package balancer

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestFastestResponseTimeSuite(t *testing.T) {
	suite.Run(t, new(FastestResponseTimeSuite))
}

type FastestResponseTimeSuite struct {
	suite.Suite
}

// Each value should be iterated 'cycles' times
func (r *FastestResponseTimeSuite) TestDefault() {
	cycles := 2
	testValues := []string{"a", "b", "c", "d"}

	updates := make(chan string, len(testValues))
	b := newFastestResponseTime("", WithUpdater(updates))

	for _, v := range testValues {
		b.Add(v)
	}

	r.Require().Equal(len(testValues), b.Amount())

	for i := 0; i < cycles; i++ {
		for j := range b.Amount() {
			r.Require().Equal(testValues[j], b.Pick(), fmt.Sprintf("cycle: %d idx: %d", i, j))
			updates <- testValues[j]
		}
	}
}

func (r *FastestResponseTimeSuite) TestEmpty() {
	cycles := 1
	defaultType := ""
	b := newFastestResponseTime(defaultType)

	r.Require().Equal(0, b.Amount())

	for i := 0; i < cycles; i++ {
		r.Require().Equal(defaultType, b.Pick())
	}
}
