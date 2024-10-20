package limiter

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestLeakyBucketSuite(t *testing.T) {
	suite.Run(t, new(LeakyBucketSuite))
}

type LeakyBucketSuite struct {
	suite.Suite
}

func (s *LeakyBucketSuite) TestDefault() {
	capacity := uint(5)
	tick := time.Millisecond * 100
	lim := newLeakyBucket(WithTick(tick), WithCapacity(capacity))

	for range capacity {
		s.Require().True(lim.Allow(context.Background()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), tick/2)
	defer cancel()

	s.Require().False(lim.Allow(ctx))

	time.Sleep(tick + time.Millisecond*30)

	s.Require().True(lim.Allow(context.Background()))
}
