package lfu

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	*Simple[K, V]
	lock sync.RWMutex
}

func NewLfu[K comparable, V any](cap int) *Cache[K, V] {
	return &Cache[K, V]{
		Simple: New[K, V](cap),
	}
}

func (l *Cache[K, V]) RecordState(recordState bool) {
	l.Simple.recordState = recordState
}

func (l *Cache[K, V]) Get(key K) (v V, ok bool) {
	l.lock.Lock()
	v, ok = l.Simple.Get(key)
	l.lock.Unlock()
	return
}

func (l *Cache[K, V]) Put(key K, value V, expireMilli ...int64) error {
	l.lock.Lock()
	err := l.Simple.Put(key, value, expireMilli...)
	l.lock.Unlock()
	return err
}

func (l *Cache[K, V]) PutIfAbsent(key K, value V, expireMilli ...int64) error {
	l.lock.Lock()
	err := l.Simple.PutIfAbsent(key, value, expireMilli...)
	l.lock.Unlock()
	return err
}

func (l *Cache[K, V]) PutAll(data map[K]V, expireMilli ...int64) error {
	l.lock.Lock()
	err := l.Simple.PutAll(data, expireMilli...)
	l.lock.Unlock()
	return err
}

func (l *Cache[K, V]) Evict(key K) V {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.Simple.Evict(key)
}

func (l *Cache[K, V]) Clear() {
	l.lock.Lock()
	l.Simple.Clear()
	l.lock.Unlock()
}

func (l *Cache[K, V]) Size() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.Simple.Size()
}

func (l *Cache[K, V]) DisPlay() string {
	return l.Simple.DisPlay()
}

func (l *Cache[K, V]) ExpireCleanTask() {

}

func (l *Cache[K, V]) remove(n int) {
	l.lock.Lock()
	l.Simple.remove(n)
	l.lock.Unlock()
}

func (l *Cache[K, V]) removeItem(entry *lfuEntry[K, V]) {
	l.lock.Lock()
	l.Simple.removeItem(entry)
	l.lock.Unlock()
}

func (l *Cache[K, V]) Contains(key K) (ok bool) {
	l.lock.RLock()
	ok = l.Simple.Contains(key)
	l.lock.RUnlock()
	return ok
}
