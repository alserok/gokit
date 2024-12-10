# Gokit

---

* [Load balancer](#load-balancer)
* [Rate limiter](#rate-limiter)
* [Circuit breaker](#circuit-breaker)
* [Cache](#cache)
* [Worker Pool](#worker-pool)
* [Retry](#retry)

## Kit for microservices written in Go

## Load balancer

---

### Round Robin

Default implementation of round-robin algo
```text
1.          2.          3.          4.
() <-       ()          ()          () <-
()          () <-       ()          () 
()          ()          () <-       ()
```

```go
package main

import "github.com/alserok/gokit/balancer"

func main() {
	defaultType := ""
	
	b := balancer.New(balancer.RoundRobin, defaultType)
}
```

### Sticky round Robin

The same round-robin, but each value is being chosen n times in a row
```text
(n = 2)
1.          2.          3.          4.
() <-       () <-       ()          ()
()          ()          () <-       () <- 
()          ()          ()          ()
```

```go
package main

import "github.com/alserok/gokit/balancer"

func main() {
	defaultType := ""
	
	b := balancer.New(balancer.StickyRoundRobin, defaultType, balancer.WithStick(3))
}
```

### Fastest response

Uses a queue implemented by channel

```text
<- () <- () <- () <-
```

```go
package main

import "github.com/alserok/gokit/balancer"

func main() {
	defaultType := ""
	updates := make(chan string, 10) // should be the same length as the number of values added, so as not to block
	
	b := balancer.New(balancer.FastestResponse, defaultType, balancer.WithUpdater(updates))
}
```

---

## Circuit breaker

Prevents clients from calling services which are not available

---

```text
() => () ✅
() => () ❌ (service is not answering bcs of internal errors, change breaker status to closed)
() => () ⛔⏳ (cancel incoming requests, wait for timeout to end and switch to open-closed)
() => () ❌ (again service is not answering bcs of internal errors, change breaker status to closed)
() => () ⛔⏳ (cancel incoming requests, wait for timeout to end and switch to open-closed)
() => () ✅ (service is available, switch to open)
```

```go
package main

import (
	"github.com/alserok/gokit/breaker"
	"time"
)

func main() {
        timeout := time.Second
	failToClose := 100
	
	b := breaker.New(timeout,failToClose)
}
```

---

## Rate limiter

---

Prevents services from being failed because of too many requests

### Leaky bucket

Initially, limiter has cap tickets, and after every tick 1 ticket is being added.

```go
package main

import (
	"github.com/alserok/gokit/limiter"
	"time"
)

func main() {
	tick := time.Second
	cap := 100
	
	l := limiter.New(limiter.LeakyBucket, limiter.WithCapacity(cap), limiter.WithTick(tick))
}
```

### Fixed window

Cap - maximum number of requests that may pass per 'tick' unit of time. After 'tick' counter rests to 0.

```go
package main

import (
	"github.com/alserok/gokit/limiter"
	"time"
)

func main() {
	tick := time.Second
	cap := 100
	
	l := limiter.New(limiter.FixedWindowCounter, limiter.WithCapacity(cap), limiter.WithTick(tick))
}
```

---

## Cache

---

### LRU 

Extincts least recently used value, uses linked list

```go
package main

import (
	"github.com/alserok/gokit/cache"
)

func main() {
	lim := uint64(10)
	
	c := cache.New(cache.LRU, lim)
}
```

#### Benchmarks

Get
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkLRUGet
BenchmarkLRUGet-12      57143128                20.62 ns/op
```

Set
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkLRUSet
BenchmarkLRUSet-12      25298149                46.60 ns/op
```

### LFU

Extincts least frequently used value, uses min heap

```go
package main

import (
	"github.com/alserok/gokit/cache"
)

func main() {
	lim := uint64(10)
	
	c := cache.New(cache.LFU, lim)
}
```

#### Benchmarks

Get
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkLFUGet
BenchmarkLFUGet-12      80156571                14.72 ns/op
```

Set
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkLFUSet
BenchmarkLFUSet-12      28314708                40.92 ns/op
```

---

## Worker pool

For time-consuming tasks recommended to use `workers=runtime.NumCPU()`, 
for fast tasks will be more efficient to use less workers.

```go
package main

import (
	"github.com/alserok/gokit/worker_pool"
	"sync"
	"sync/atomic"
	"runtime"
)

func main() {
	counter := int64(0)
	fn := func(c *int64) error {
		atomic.AddInt64(c, 1)
		return nil
	}

	workers := 3
	// init worker pool with 3 workers
	p := worker_pool.NewWorkerPool(fn, int64(workers))
	// or p.Shutdown() to stop immediately
	defer p.Stop() // wait for workers to process all data already added
	
	// launch goroutines(workers)
	p.Start()
	
	// adds data for worker fn, returns true if successfully
	p.Add(context.Background(), &counter)
}
```

#### Benchmarks

Function executed by worker
```go
func(_ any) error {
	time.Sleep(time.Millisecond * 30)
	return nil
}
```

100 workers
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkNewWorkerPoolWith100Workers
BenchmarkNewWorkerPoolWith100Workers-12    	29727832	        34.60 ns/op
```

10 workers
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkNewWorkerPoolWith10Workers
BenchmarkNewWorkerPoolWith10Workers-12    	30755013	        33.37 ns/op
```

1 worker
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkNewWorkerPoolWith1Worker
BenchmarkNewWorkerPoolWith1Worker-12    	30688845	        33.42 ns/op
```

NumCPU(12) workers
```text
cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
BenchmarkNewWorkerPoolWithNumCPUWorkers
BenchmarkNewWorkerPoolWithNumCPUWorkers-12    	34146645	        33.58 ns/op
```

---

## Retry

Mechanism that allows to repeat action on failure.

If you omit customizers, will be set default values. 

### Exponential Backoff

Interval between retries multiplies by `exponent`

**Customizers**
```go
func WithAmount(amount uint) Customizer {}

func WithInterval(interval time.Duration) Customizer {}

func WithExponent(exp float32) Customizer {}
```

```go
package main

import (
	"errors"
	"fmt"
	"github.com/alserok/gokit/retry"
	"time"
)

func main() {
	r := retry.New(retry.ExponentialBackoff, retry.WithAmount(uint(3)), retry.WithInterval(time.Millisecond*5), retry.WithExponent(1.25))

	errorFn := func() error {
		return errors.New("error")
	}

	if err := r.Execute(errorFn); err != nil {
		panic("all retries finished with error: " + err.Error())
	}

	fmt.Println("some retry has finished without error")
}

```

### Fixed Interval

Interval between retries is constant

**Customizers**
```go
func WithAmount(amount uint) Customizer {}

func WithInterval(interval time.Duration) Customizer {}
```

```go
package main

import (
	"errors"
	"fmt"
	"github.com/alserok/gokit/retry"
	"time"
)

func main() {
	r := retry.New(retry.FixedInterval, retry.WithAmount(uint(3)), retry.WithInterval(time.Millisecond*5))

	errorFn := func() error {
		return errors.New("error")
	}

	if err := r.Execute(errorFn); err != nil {
		panic("all retries finished with error: " + err.Error())
	}

	fmt.Println("some retry has finished without error")
}
```

### Linear Backoff

Interval between retries increases by `step`

**Customizers**
```go
func WithAmount(amount uint) Customizer {}

func WithInterval(interval time.Duration) Customizer {}

func WithStep(step float32) Customizer {}
```

```go
package main

import (
	"errors"
	"fmt"
	"github.com/alserok/gokit/retry"
	"time"
)

func main() {
	r := retry.New(retry.LinearBackoff, retry.WithAmount(uint(3)), retry.WithInterval(time.Millisecond*5), retry.WithStep(0.5))

	errorFn := func() error {
		return errors.New("error")
	}

	if err := r.Execute(errorFn); err != nil {
		panic("all retries finished with error: " + err.Error())
	}

	fmt.Println("some retry has finished without error")
}
```

### Successive Interval

Interval between retries depends on previous amount of failures

**Customizers**
```go
func WithAmount(amount uint) Customizer {}

func WithInterval(interval time.Duration) Customizer {}
```

```go
package main

import (
	"errors"
	"fmt"
	"github.com/alserok/gokit/retry"
	"time"
)

func main() {
	r := retry.New(retry.SuccessiveInterval, retry.WithAmount(uint(3)), retry.WithInterval(time.Millisecond*5))

	errorFn := func() error {
		return errors.New("error")
	}

	if err := r.Execute(errorFn); err != nil {
		panic("all retries finished with error: " + err.Error())
	}

	fmt.Println("some retry has finished without error")
}
```