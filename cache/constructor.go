package cache

import "context"

type Cache[T any] interface {
	Set(ctx context.Context, key string, val T)
	Get(ctx context.Context, key string) T
}

const (
	LRU = iota
	LFU
)

func New[T any](cacheType uint) Cache[T] {
	switch cacheType {
	case LFU:
		return nil
	case LRU:
		return nil
	default:
		panic("invalid cache type")
	}
}
