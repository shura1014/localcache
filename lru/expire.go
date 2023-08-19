package lru

import (
	"fmt"
	"github.com/shura1014/common/goerr"
	"github.com/shura1014/common/type/atom"
	"github.com/shura1014/common/utils/timeutil"
)

type lruEntry[V any] struct {
	value  V
	expire int64
}

func (e *lruEntry[V]) IsExpired() bool {
	if e.expire == 0 {
		return false
	}
	return e.expire < timeutil.MilliSeconds()
}

// Cache
// 1. ICache的一个实现
// 2. 带有过期时间
// 3. 并发安全的lru
type Cache[K comparable, V any] struct {
	lru           *Map[K, *lruEntry[V]]
	defaultExpire int64
	// 状态
	recordState bool
	miss        *atom.Int32
	hit         *atom.Int32
}

func NewLruCache[K comparable, V any](cap int) *Cache[K, V] {
	if cap == 0 {
		cap = 1024
	}
	return &Cache[K, V]{
		lru:  NewLruMap[K, *lruEntry[V]](cap),
		miss: atom.NewInt32(),
		hit:  atom.NewInt32(),
	}
}

func (l *Cache[K, V]) RecordState(recordState bool) {
	l.recordState = recordState
}

// Put expireMilli 毫秒值
func (l *Cache[K, V]) Put(key K, value V, expireMilli ...int64) error {
	expire := computeExpire(l.defaultExpire, expireMilli...)
	e := &lruEntry[V]{
		value:  value,
		expire: expire,
	}
	l.lru.Put(key, e)
	return nil
}

// PutIfAbsent expireMilli 毫秒值
func (l *Cache[K, V]) PutIfAbsent(key K, value V, expireMilli ...int64) error {
	_, ok := l.lru.Simple.data[key]
	if ok {
		return goerr.Text("The key %v already exists", key)
	}
	return l.Put(key, value, expireMilli...)
}

// PutAll expireMilli 毫秒值
func (l *Cache[K, V]) PutAll(data map[K]V, expireMilli ...int64) error {
	expireTime := computeExpire(l.defaultExpire, expireMilli...)
	for k, v := range data {
		l.lru.Put(k, &lruEntry[V]{
			value:  v,
			expire: expireTime,
		})
	}
	return nil
}

func (l *Cache[K, V]) Get(key K) (v V, ok bool) {
	l.lru.mu.Lock()
	defer l.lru.mu.Unlock()
	if element, ok := l.lru.Simple.data[key]; ok {
		e := element.Value
		// 如果过期不移动到最前面，自动容量满了后删除
		if !e.IsExpired() {
			if l.recordState {
				l.hit.Add(1)
			}
			l.lru.Simple.list.MoveToFront(element)

			return element.Value.value, true
		} else {
			// 快速删除
			l.lru.Simple.list.MoveToBack(element)
		}
	}
	if l.recordState {
		l.miss.Add(1)
	}
	return
}

func (l *Cache[K, V]) Contains(key K) (ok bool) {
	ok = l.lru.Contains(key)
	return ok
}

func (l *Cache[K, V]) Evict(key K) (v V) {
	evict := l.lru.Evict(key)
	if evict != nil {
		return evict.value
	}
	return
}

func (l *Cache[K, V]) Clear() {
	l.lru.Clear()
	l.miss.Swap(0)
	l.hit.Swap(0)
}

func (l *Cache[K, V]) Size() int {
	return l.lru.Size()
}

func (l *Cache[K, V]) DisPlay() string {
	return fmt.Sprintf("len: %d hit: %d miss: %d\n", l.lru.Size(), l.hit.Load(), l.miss.Load())
}

// ExpireCleanTask 不需要清理
func (l *Cache[K, V]) ExpireCleanTask() {
	return
}

// defaultExpire 默认的过期时间
// appoint 指定的过期时间
func computeExpire(defaultExpire int64, appoint ...int64) int64 {
	expire := defaultExpire
	if len(appoint) > 0 {
		expire = appoint[0]
	}

	return timeutil.MilliSeconds() + expire
}
