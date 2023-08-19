package lfu

import (
	"fmt"
	"github.com/shura1014/common/container/list/generics"
	"github.com/shura1014/common/goerr"
	"github.com/shura1014/common/type/atom"
	"github.com/shura1014/common/utils/timeutil"
)

// Simple 不加锁
type Simple[K comparable, V any] struct {
	defaultExpire int64
	data          map[K]*lfuEntry[K, V]
	freqList      *generics.List[*FreqEntry[K, V]]
	cap           int
	len           int
	// 状态
	recordState bool
	miss        *atom.Int32
	hit         *atom.Int32
}

func New[K comparable, V any](cap int) *Simple[K, V] {
	if cap == 0 {
		cap = 1024
	}
	l := &Simple[K, V]{
		cap:      cap,
		data:     make(map[K]*lfuEntry[K, V]),
		freqList: generics.NewList[*FreqEntry[K, V]](),
		miss:     atom.NewInt32(),
		hit:      atom.NewInt32(),
	}
	// 初始化一条，默认第一次put频次就是1
	l.freqList.PushFront(&FreqEntry[K, V]{
		freq:    1,
		entries: make(map[*lfuEntry[K, V]]byte),
	})
	return l
}

func (l *Simple[K, V]) RecordState(recordState bool) {
	l.recordState = recordState
}

func (l *Simple[K, V]) Get(key K) (v V, ok bool) {
	if entry, ok := l.data[key]; ok {
		if !entry.IsExpired() {
			if l.recordState {
				l.hit.Add(1)
			}
			l.increment(entry)
			return entry.value, true
		} else {
			// todo
		}
	}
	if l.recordState {
		l.miss.Add(1)
	}
	return
}

func (l *Simple[K, V]) Put(key K, value V, expireMilli ...int64) error {
	expire := computeExpire(l.defaultExpire, expireMilli...)
	if element, ok := l.data[key]; ok {
		element.value = value
		element.expire = expire
		l.increment(element)
		return nil
	}
	// 检查一下空间，如果空间满了 ，那么删除一个最近不常用的
	if l.len >= l.cap {
		l.remove(1)
	}
	e := &lfuEntry[K, V]{
		key:    key,
		value:  value,
		expire: expire,
	}

	freqEntry := l.freqList.Front()
	freqEntry.Value.entries[e] = 1
	e.freqList = freqEntry
	l.data[key] = e

	l.len++
	return nil
}

func (l *Simple[K, V]) PutIfAbsent(key K, value V, expireMilli ...int64) error {
	_, ok := l.data[key]
	if ok {
		return goerr.Text("The key %v already exists", key)
	}
	return l.Put(key, value, expireMilli...)
}

func (l *Simple[K, V]) PutAll(data map[K]V, expireMilli ...int64) error {
	for k, v := range data {
		_ = l.Put(k, v, expireMilli...)
	}
	return nil
}

func (l *Simple[K, V]) Evict(key K) (v V) {
	if entry, ok := l.data[key]; ok {
		l.removeItem(entry)
		l.len -= 1
		return entry.value
	}

	return
}

func (l *Simple[K, V]) Clear() {
	l.data = make(map[K]*lfuEntry[K, V])
	l.freqList = l.freqList.Init()
	l.len = 0
	l.miss.Swap(0)
	l.hit.Swap(0)
}

func (l *Simple[K, V]) Size() int {
	return l.len
}

func (l *Simple[K, V]) DisPlay() string {
	return fmt.Sprintf("len: %d hit: %d miss: %d\n", l.len, l.hit.Load(), l.miss.Load())
}

func (l *Simple[K, V]) ExpireCleanTask() {

}

func (l *Simple[K, V]) remove(n int) {
	freqEntry := l.freqList.Front()
	for i := 0; i < n; {
		if freqEntry == nil {
			return
		} else {
			for item := range freqEntry.Value.entries {
				if i >= n {
					return
				}
				l.removeItem(item)
				i++
				l.len -= 1
			}
			// 上一个entry删除完了还没有达到n的数量，那么下一个节点继续删除
			freqEntry = freqEntry.Next()
		}
	}

}

func (l *Simple[K, V]) removeItem(entry *lfuEntry[K, V]) {
	freqEntry := entry.freqList.Value
	delete(l.data, entry.key)
	delete(freqEntry.entries, entry)
	if l.isDeleteAble(freqEntry) {
		l.freqList.Remove(entry.freqList)
	}
}

func (l *Simple[K, V]) increment(entry *lfuEntry[K, V]) {
	freqElement := entry.freqList
	nextFreq := freqElement.Value.freq + 1
	nextFreqElement := freqElement.Next()
	delete(freqElement.Value.entries, entry)
	// 如果说已经是尾节点，或者没有下一个节点的freq比较大，创建一个新的freqElement
	if nextFreqElement == nil || nextFreqElement.Value.freq > nextFreq {
		// 	复用
		if l.isDeleteAble(freqElement.Value) {
			freqElement.Value.freq = nextFreq
			nextFreqElement = freqElement
		} else {
			nextFreqElement = l.freqList.InsertAfter(&FreqEntry[K, V]{
				freq:    nextFreq,
				entries: make(map[*lfuEntry[K, V]]byte),
			}, freqElement)
		}
	} else if nextFreqElement.Value.freq == nextFreq && l.isDeleteAble(freqElement.Value) {
		l.freqList.Remove(freqElement)
	}

	nextFreqElement.Value.entries[entry] = 1
	entry.freqList = nextFreqElement
}

// 是否满足删除条件 e.freq != 1 为默认初始值，不删除
func (l *Simple[K, V]) isDeleteAble(e *FreqEntry[K, V]) bool {
	if len(e.entries) == 0 && e.freq != 1 {
		return true
	}
	return false
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

func (l *Simple[K, V]) Contains(key K) (ok bool) {
	_, ok = l.data[key]
	return ok
}
