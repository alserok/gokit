package balancer

import (
	"testing"
)

func BenchmarkAddPickStickyRoundRobin(b *testing.B) {
	rr := newStickyRoundRobin("")

	for i := 0; i < b.N; i++ {
		rr.Add("q")
		rr.Pick()
	}
}

func BenchmarkPickStickyRoundRobin(b *testing.B) {
	rr := newStickyRoundRobin("")
	rr.Add("q1")
	rr.Add("q1")

	for i := 0; i < b.N; i++ {
		rr.Pick()
	}
}

func BenchmarkAddStickyRoundRobin(b *testing.B) {
	rr := newStickyRoundRobin("")

	for i := 0; i < b.N; i++ {
		rr.Add("q")
	}
}
