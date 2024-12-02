package worker_pool

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	suite.Run(t, new(WorkerPoolSuite))
}

type WorkerPoolSuite struct {
	suite.Suite
}

func (suite *WorkerPoolSuite) TestNewWorkerPool() {
	p := NewWorkerPool(func(str string) error {
		return nil
	}, 1)

	suite.Require().NotNil(p)
}

func (suite *WorkerPoolSuite) TestDefaultCounter() {
	counter := int64(0)
	fn := func(c *int64) error {
		atomic.AddInt64(&counter, 1)
		return nil
	}

	workers := 3
	p := newWorkerPool(fn, int64(workers))
	suite.Require().NotNil(p)
	suite.Require().NotNil(p.errors)
	suite.Require().NotNil(p.input)
	suite.Require().NotNil(p.workers)
	suite.Require().NotNil(p.exec)
	suite.Require().Equal(p.maxWorkers, int64(workers))
	suite.Require().Equal(p.currWorkers, int64(0))

	p.Start()

	for range workers {
		suite.Require().True(p.Add(context.Background(), &counter))
	}

	time.Sleep(time.Microsecond * 100)

	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.currWorkers))
	suite.Require().Equal(workers, int(atomic.LoadInt64(&counter)))

	p.Stop()

	suite.Require().Equal(p.currWorkers, int64(0))
}

func (suite *WorkerPoolSuite) TestSetWorkers() {
	fn := func(arg bool) error {
		return nil
	}

	workers := 1
	p := newWorkerPool(fn, int64(workers))
	suite.Require().NotNil(p)

	p.Start()

	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.maxWorkers))
	time.Sleep(time.Millisecond * 10)
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.currWorkers))

	workers = 2
	p.SetWorkers(uint(workers))
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.maxWorkers))
	time.Sleep(time.Millisecond * 10)
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.currWorkers))

	workers = 5
	p.SetWorkers(uint(workers))
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.maxWorkers))
	time.Sleep(time.Millisecond * 10)
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.currWorkers))

	workers = 2
	p.SetWorkers(uint(workers))
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.maxWorkers))
	time.Sleep(time.Millisecond * 10)
	suite.Require().Equal(int64(workers), atomic.LoadInt64(&p.currWorkers))

	p.Stop()
	suite.Require().Equal(p.currWorkers, int64(0))
}

func (suite *WorkerPoolSuite) TestError() {
	p := newWorkerPool(func(_ any) error {
		return errors.New("error")
	}, 1)
	suite.Require().NotNil(p)

	p.Start()

	suite.Require().True(p.Add(context.Background(), nil))
	suite.Require().NotNil(<-p.errors)
	suite.Require().True(p.Add(context.Background(), nil))
	suite.Require().NotNil(<-p.errors)
	suite.Require().True(p.Add(context.Background(), nil))
	suite.Require().NotNil(<-p.errors)
}

func (suite *WorkerPoolSuite) TestShutdown() {
	finished := false
	p := newWorkerPool(func(b *bool) error {
		time.Sleep(time.Millisecond * 10)
		*b = true
		return nil
	}, 1)

	p.Start()

	suite.Require().True(p.Add(context.Background(), &finished))
	p.Shutdown()

	suite.Require().Equal(0, int(atomic.LoadInt64(&p.currWorkers)))
	suite.Require().False(finished)
}

func (suite *WorkerPoolSuite) TestStop() {
	finished := false
	p := newWorkerPool(func(b *bool) error {
		time.Sleep(time.Millisecond * 10)
		*b = true
		return nil
	}, 1)

	p.Start()

	suite.Require().True(p.Add(context.Background(), &finished))
	p.Stop()

	suite.Require().Equal(0, int(p.currWorkers))
	suite.Require().True(finished)
}
