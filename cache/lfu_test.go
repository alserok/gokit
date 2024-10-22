package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestLFUSuite(t *testing.T) {
	suite.Run(t, new(lfuSuite))
}

type lfuSuite struct {
	suite.Suite
}

func (s *lfuSuite) TestDefault() {
	type testStruct struct {
		value string
	}
	lim := 30
	skip := 10

	c := newLFU[testStruct](lim)

	// insert in cache
	for i := range lim {
		c.Set(context.Background(), fmt.Sprintf("key %d", i), testStruct{value: "value"})
	}

	// get values till 'lim' - 'skip' index, all got values frequency is incremented by 1
	for i := range lim - skip {
		val := c.Get(context.Background(), fmt.Sprintf("key %d", i))
		s.Require().NotNil(val)
	}

	// insert 'skip' new values and increment their frequency by 1
	for i := range skip {
		c.Set(context.Background(), fmt.Sprintf("some key %d", i), testStruct{value: "some value"})
		c.Get(context.Background(), fmt.Sprintf("some key %d", i))
	}

	// check if skipped values were extincted
	for i := lim - skip; i < lim; i++ {
		val := c.Get(context.Background(), fmt.Sprintf("key %d", i))
		s.Require().Nil(val)
	}
}
