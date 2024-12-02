# Gokit

---
* Load balancer
* Rate limiter
* Circuit breaker
* Cache
* Worker Pool

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

```go
package main

import (
	"github.com/alserok/gokit/worker_pool"
	"sync"
)

func main() {
	counter := 0
	fn := func(c *int) error {
		mu := sync.Mutex{}

		mu.Lock()
		*c += 1
		mu.Unlock()

		return nil
	}

	workers := 3
	// init worker pool with 3 workers
	p := worker_pool.NewWorkerPool(fn, int64(workers))
	defer p.Stop()
	
	// launch goroutines(workers)
	go p.Start()
	
	// added data for worker fn
	p.Add(context.Background(), &counter)
}
```