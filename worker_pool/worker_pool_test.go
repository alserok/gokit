package worker_pool

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"sync"
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
	counter := 0
	fn := func(c *int) error {
		mu := sync.Mutex{}

		mu.Lock()
		*c += 1
		mu.Unlock()

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

	go p.Start()

	for range workers {
		p.Add(context.Background(), &counter)
	}

	time.Sleep(time.Millisecond * 20)

	suite.Require().Equal(p.currWorkers, int64(workers))
	suite.Require().Equal(counter, workers)

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

	go p.Start()

	suite.Require().Equal(int64(workers), p.maxWorkers)
	time.Sleep(time.Millisecond * 5)
	suite.Require().Equal(int64(workers), p.currWorkers)

	workers = 2
	p.SetWorkers(uint(workers))
	suite.Require().Equal(int64(workers), p.maxWorkers)
	time.Sleep(time.Millisecond * 5)
	suite.Require().Equal(int64(workers), p.currWorkers)

	workers = 5
	p.SetWorkers(uint(workers))
	suite.Require().Equal(int64(workers), p.maxWorkers)
	time.Sleep(time.Millisecond * 5)
	suite.Require().Equal(int64(workers), p.currWorkers)

	workers = 2
	p.SetWorkers(uint(workers))
	suite.Require().Equal(int64(workers), p.maxWorkers)
	time.Sleep(time.Millisecond * 5)
	suite.Require().Equal(int64(workers), p.currWorkers)

	p.Stop()
	suite.Require().Equal(p.currWorkers, int64(0))
}

func (suite *WorkerPoolSuite) TestError() {
	p := newWorkerPool(func(_ any) error {
		return errors.New("error")
	}, 1)
	suite.Require().NotNil(p)

	go p.Start()

	p.Add(context.Background(), nil)
	suite.Require().NotNil(<-p.errors)
	p.Add(context.Background(), nil)
	suite.Require().NotNil(<-p.errors)
	p.Add(context.Background(), nil)
	suite.Require().NotNil(<-p.errors)
}
