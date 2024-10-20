package cache

import "context"

type Cache[T any] interface {
	Set(ctx context.Context, key string, val T)
	Get(ctx context.Context, key string) *T
}

const (
	LRU = iota
	LFU
)

func New[T any](cacheType uint, lim uint64) Cache[T] {
	switch cacheType {
	case LFU:
		return nil
	case LRU:
		return newLRU[T](int(lim))
	default:
		panic("invalid cache type")
	}
}
