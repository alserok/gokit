package cache

import (
	"context"
	"strconv"
	"testing"
)

func BenchmarkLFUGet(b *testing.B) {
	c := newLFU[int](b.N)
	stringKeys := make([]string, 0, b.N)
	for i := 0; i < b.N; i++ {
		stringKeys = append(stringKeys, strconv.Itoa(b.N))
		c.Set(context.Background(), strconv.Itoa(b.N), b.N)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set(context.Background(), stringKeys[i], b.N)
	}
}

func BenchmarkLFUSet(b *testing.B) {
	c := newLFU[int](b.N)
	for i := 0; i < b.N; i++ {
		c.Set(context.Background(), strconv.Itoa(b.N), b.N)
	}
}
