# Gokit

---
* Load balancer
* Rate limiter
* Circuit breaker

## Kit for microservices written in Go

### Load balancer

---

#### Round Robin

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

#### Sticky round Robin

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

### Circuit breaker

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