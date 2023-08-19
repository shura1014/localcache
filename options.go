package localcache

import (
	"github.com/shura1014/localcache/lfu"
	"github.com/shura1014/localcache/lru"
)

const (
	LRU = iota
	LFU
)

var defaultAlg = LRU

type EvictedHandler func(string, any)

type Options[K comparable, V any] struct {
	// 毫秒值
	defaultExpire int64
	onEvicted     EvictedHandler
	cap           int
	alg           int
	recordState   bool
}

func NewBuilder[K comparable, V any]() Options[K, V] {
	return Options[K, V]{}
}

func (option Options[K, V]) SetEvictHandler(handler EvictedHandler) Options[K, V] {
	option.onEvicted = handler
	return option
}

func (option Options[K, V]) SetGlobalExpire(expireMilli int64) Options[K, V] {
	option.defaultExpire = expireMilli
	return option
}

func (option Options[K, V]) EnableRecord() Options[K, V] {
	option.recordState = true
	return option
}

func (option Options[K, V]) Cap(cap int) Options[K, V] {
	option.cap = cap
	return option
}

func (option Options[K, V]) Alg(alg int) Options[K, V] {
	option.alg = alg
	return option
}

func (option Options[K, V]) Build() ICache[K, V] {
	cache := option.getAlg(option.alg)
	cache.ExpireCleanTask()
	return cache
}

// 内置的缓存
func (option Options[K, V]) getAlg(alg int) ICache[K, V] {
	if alg == 0 {
		alg = defaultAlg
	}
	switch alg {
	case LRU:
		l := lru.NewLruCache[K, V](option.cap)
		l.RecordState(option.recordState)
		return l
	case LFU:
		l := lfu.NewLfu[K, V](option.cap)
		l.RecordState(option.recordState)
		return l
	}
	return nil

}
