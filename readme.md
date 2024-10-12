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