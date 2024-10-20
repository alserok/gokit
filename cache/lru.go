package cache

import (
	"context"
	"sync"
)

func newLRU[T any](lim int) *lru[T] {
	return &lru[T]{
		lim:    lim,
		values: make(map[string]*node[T]),
	}
}

type lru[T any] struct {
	head *node[T]
	tail *node[T]

	lim int

	values map[string]*node[T]

	mu sync.RWMutex
}

func (l *lru[T]) Set(ctx context.Context, key string, val T) {
	n := &node[T]{
		key:  key,
		val:  val,
		next: l.head,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.head == nil {
		l.tail = n
	} else {
		l.head.prev = n
	}

	if len(l.values) >= l.lim {
		delete(l.values, l.tail.key)
		l.tail = l.tail.prev
	}

	l.head = n
	l.values[key] = n
}

func (l *lru[T]) Get(ctx context.Context, key string) *T {
	l.mu.RLock()
	val, ok := l.values[key]
	l.mu.RUnlock()

	if !ok {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if val.prev == nil {
		return &val.val
	}

	val.prev.next = val.next

	if val.next == nil {
		val.prev.next = nil
		l.tail = val.prev
	} else {
		val.next.prev = val.prev
	}

	val.next, l.head.prev = l.head, val
	l.head = val

	return &val.val
}

type node[T any] struct {
	val T
	key string

	next *node[T]
	prev *node[T]
}
