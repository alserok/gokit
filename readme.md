# Gokit

---
* Load balancer
* Rate limiter
* Circuit breaker

## Kit for microservices written in Go

### Load balancer

---

#### Round Robin
```go
package main

import "github.com/alserok/gokit/balancer"

func main() {
	defaultType := ""
	
	b := balancer.New(balancer.RoundRobin, defaultType)
}
```

#### Sticky round Robin
```go
package main

import "github.com/alserok/gokit/balancer"

func main() {
	defaultType := ""
	
	b := balancer.New(balancer.StickyRoundRobin, defaultType)
}
```