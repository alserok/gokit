package cache

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestLRUSuite(t *testing.T) {
	suite.Run(t, new(lruSuite))
}

type lruSuite struct {
	suite.Suite
}

func (s *lruSuite) TestDefault() {
	type testStruct struct {
		value string
	}
	lim := 5

	cache := newLRU[testStruct](int64(lim))

	for i := range lim {
		cache.Set(context.Background(), fmt.Sprintf("key %d", i), testStruct{value: "a"})
	}

	for i := range lim {
		val := cache.Get(context.Background(), fmt.Sprintf("key %d", i))
		s.Require().NotNil(val)
	}

	cache.Set(context.Background(), "some extra value key", testStruct{value: "a"})

	val := cache.Get(context.Background(), "key 0")
	s.Require().Nil(val)

	cache.Set(context.Background(), "some extra value key 1", testStruct{value: "a"})

	val = cache.Get(context.Background(), "key 1")
	s.Require().Nil(val)
}
