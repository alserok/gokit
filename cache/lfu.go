package cache

import (
	"container/heap"
	"context"
	"sync"
)

func newLFU[T any](lim int) *lfu[T] {
	mHeap := &minHeap{}
	heap.Init(mHeap)

	return &lfu[T]{
		lim:    lim,
		heap:   mHeap,
		values: make(map[string]*nodeLFU[any]),
	}
}

type lfu[T any] struct {
	heap *minHeap

	lim int

	values map[string]*nodeLFU[any]

	mu sync.RWMutex
}

func (l *lfu[T]) Get(ctx context.Context, key string) *T {
	l.mu.Lock()
	defer l.mu.Unlock()

	val, ok := l.values[key]
	if !ok {
		return nil
	}

	val.freq++
	heap.Fix(l.heap, val.index)

	res := val.val.(T)

	return &res
}

func (l *lfu[T]) Set(ctx context.Context, key string, val T) {
	l.mu.RLock()
	_, ok := l.values[key]
	l.mu.RUnlock()

	if ok {
		return
	}

	node := &nodeLFU[any]{
		freq: 0,
		val:  val,
		key:  key,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.values) == l.lim {
		outVal := heap.Pop(l.heap)
		delete(l.values, outVal.(*nodeLFU[any]).key)
	}

	heap.Push(l.heap, node)
	l.values[key] = node
}

type nodeLFU[T any] struct {
	freq int

	key string
	val T

	index int
}

type minHeap []*nodeLFU[any]

func (h *minHeap) Len() int {
	return len(*h)
}

func (h *minHeap) Less(i, j int) bool {
	return (*h)[i].freq < (*h)[j].freq
}

func (h *minHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index = i
	(*h)[j].index = j
}

func (h *minHeap) Push(x any) {
	node := x.(*nodeLFU[any])
	node.index = len(*h)

	*h = append(*h, node)

}

func (h *minHeap) Pop() any {
	n := len(*h)

	node := (*h)[n-1]
	*h = (*h)[:n-1]

	return node
}
