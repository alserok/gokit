package balancer

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestConstructorSuite(t *testing.T) {
	suite.Run(t, new(ConstructorSuite))
}

type ConstructorSuite struct {
	suite.Suite
}

func (c *ConstructorSuite) TestDefault() {
	types := []uint{RoundRobin, StickyRoundRobin}

	for _, t := range types {
		c.Require().Equal(t, New(t, "").Type())
	}
}
