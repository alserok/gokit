package balancer

import (
	"testing"
)

func BenchmarkAddPickRoundRobin(b *testing.B) {
	rr := newRoundRobin("")

	for i := 0; i < b.N; i++ {
		rr.Add("q")
		rr.Pick()
	}
}

func BenchmarkPickRoundRobin(b *testing.B) {
	rr := newRoundRobin("")
	rr.Add("q1")
	rr.Add("q1")

	for i := 0; i < b.N; i++ {
		rr.Pick()
	}
}

func BenchmarkAddRoundRobin(b *testing.B) {
	rr := newRoundRobin("")

	for i := 0; i < b.N; i++ {
		rr.Add("q")
	}
}
