package lru

import (
	"sync"
)

// Map 在 Simple 的基础上加锁，并发安全的
type Map[K comparable, V any] struct {
	mu sync.RWMutex
	*Simple[K, V]
}

func NewLruMap[K comparable, V any](cap int) *Map[K, V] {
	return &Map[K, V]{
		Simple: New[K, V](cap),
	}
}

func (l *Map[K, V]) Put(key K, value V) {
	l.mu.Lock()
	l.Simple.Put(key, value)
	l.mu.Unlock()
}

func (l *Map[K, V]) Get(key K) (v V, ok bool) {
	l.mu.Lock()
	v, ok = l.Simple.Get(key)
	l.mu.Unlock()
	return
}

func (l *Map[K, V]) Evict(key K) (v V) {
	l.mu.Lock()
	v = l.Simple.Evict(key)
	l.mu.Unlock()
	return
}

func (l *Map[K, V]) Size() int {
	return l.Simple.Size()
}

func (l *Map[K, V]) Clear() {
	l.mu.Lock()
	l.Simple.Clear()
	l.mu.Unlock()
}

func (l *Map[K, V]) Join(glue string) string {
	l.mu.RLock()
	result := l.Simple.Join(glue)
	l.mu.RUnlock()
	return result
}

func (l *Map[K, V]) Contains(key K) (ok bool) {
	l.mu.RLock()
	ok = l.Simple.Contains(key)
	l.mu.RUnlock()
	return
}

func (l *Map[K, V]) String() string {
	return l.Simple.String()
}
