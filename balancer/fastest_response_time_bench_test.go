package balancer

import (
	"testing"
)

func BenchmarkAddPickFastestResponse(b *testing.B) {
	frt := newFastestResponse("")

	for i := 0; i < b.N; i++ {
		frt.Add("q")
		frt.Pick()
	}
}

func BenchmarkPickFastestResponse(b *testing.B) {
	updates := make(chan string, b.N)
	go func() {
		for i := 0; i < b.N; i++ {
			updates <- "val"
		}
	}()

	frt := newFastestResponse("", WithUpdater(updates))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		frt.Pick()
	}
}

func BenchmarkAddFastestResponse(b *testing.B) {
	updates := make(chan string, b.N)
	frt := newFastestResponse("", WithUpdater(updates))

	for i := 0; i < b.N; i++ {
		frt.Add("q")
	}
}
