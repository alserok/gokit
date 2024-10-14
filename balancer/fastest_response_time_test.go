package balancer

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestFastestResponseSuite(t *testing.T) {
	suite.Run(t, new(FastestResponseSuite))
}

type FastestResponseSuite struct {
	suite.Suite
}

func (r *FastestResponseSuite) TestDefault() {
	cycles := 2
	testValues := []string{"a", "b", "c", "d"}

	updates := make(chan string, len(testValues))
	b := newFastestResponse("", WithUpdater(updates))

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

func (r *FastestResponseSuite) TestRemove() {
	testValues := []string{"a", "b", "c", "d"}

	updates := make(chan string, len(testValues))
	b := newFastestResponse("", WithUpdater(updates))

	for _, v := range testValues {
		b.Add(v)
	}

	// if value is already added it will be listed at least 1 time
	for i := range b.Amount() {
		r.Require().Equal(testValues[i], b.Pick())
		if i != 0 {
			updates <- testValues[i]
		}
	}

	b.Remove(testValues[0])
	r.Require().Equal(len(testValues)-1, b.Amount())

	for i := range b.Amount() {
		r.Require().Equal(testValues[i+1], b.Pick())
	}
}

func (r *FastestResponseSuite) TestEmpty() {
	cycles := 1
	defaultType := ""
	b := newFastestResponse(defaultType)

	r.Require().Equal(0, b.Amount())

	for i := 0; i < cycles; i++ {
		r.Require().Equal(defaultType, b.Pick())
	}
}
