package cache

import (
	"context"
	"sync"
)

func newLRU[T any](lim int) *lru[T] {
	return &lru[T]{
		lim:    lim,
		values: make(map[string]*nodeLRU[T]),
	}
}

type lru[T any] struct {
	head *nodeLRU[T]
	tail *nodeLRU[T]

	lim int

	values map[string]*nodeLRU[T]

	mu sync.RWMutex
}

func (l *lru[T]) Set(ctx context.Context, key string, val T) {
	l.mu.RLock()
	_, ok := l.values[key]
	l.mu.RUnlock()

	if ok {
		return
	}

	n := &nodeLRU[T]{
		key: key,
		val: val,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.head == nil {
		l.head = n
		l.tail = n
	} else {
		n.next = l.head
		l.head.prev = n
		l.head = n
	}

	l.values[key] = n

	if len(l.values) > l.lim {
		delete(l.values, l.tail.key)
		l.tail = l.tail.prev
	}
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

	if val == l.head {
		return &val.val
	}

	if val == l.tail {
		l.tail = val.prev
		l.tail.next = nil
	} else {
		if val.prev != nil {
			val.prev.next = val.next
		}
		if val.next != nil {
			val.next.prev = val.prev
		}
	}

	val.next = l.head
	l.head.prev = val
	val.prev = nil
	l.head = val

	return &val.val
}

type nodeLRU[T any] struct {
	val T
	key string

	next *nodeLRU[T]
	prev *nodeLRU[T]
}
