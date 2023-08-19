package lru

import (
	"bytes"
	"fmt"
	"github.com/shura1014/common/container/list/generics"
	"github.com/shura1014/common/goerr"
)

// Simple 简单的缓存
// 不加锁，没有过期功能
type Simple[K comparable, V any] struct {
	cap  int
	data map[K]*generics.MapElement[K, V]
	list *generics.MapList[K, V]
}

// New cap容量
func New[K comparable, V any](cap int) *Simple[K, V] {
	return &Simple[K, V]{
		cap:  cap,
		list: generics.NewMapList[K, V](),
		data: make(map[K]*generics.MapElement[K, V], cap),
	}
}

func (l *Simple[K, V]) Put(key K, value V) {
	if element, ok := l.data[key]; ok {
		// 移到首节点
		l.list.MoveToFront(element)
		element.Value = value
		return
	}

	// 移除最后一个元素
	if l.list.Len() >= l.cap {
		element := l.list.Back()
		if element == nil {
			return
		}
		l.list.Remove(element)
		delete(l.data, element.Key)
	}

	// 放置首节点
	element := l.list.PushFront(key, value)
	l.data[key] = element
}

func (l *Simple[K, V]) PutIfAbsent(key K, value V) error {
	_, ok := l.data[key]
	if ok {
		return goerr.Text("The key %v already exists", key)
	}
	l.Put(key, value)
	return nil
}

func (l *Simple[K, V]) PutAll(data map[K]V) error {
	for k, v := range data {
		l.Put(k, v)
	}
	return nil
}

func (l *Simple[K, V]) Contains(key K) bool {
	_, ok := l.data[key]
	return ok
}

func (l *Simple[K, V]) Get(key K) (v V, ok bool) {
	if element, ok := l.data[key]; ok {
		l.list.MoveToFront(element)
		return element.Value, true
	}
	return
}

func (l *Simple[K, V]) Evict(key K) (v V) {
	if element, ok := l.data[key]; ok {
		value := element.Value
		l.list.Remove(element)
		delete(l.data, element.Key)
		return value
	}
	return
}

func (l *Simple[K, V]) Size() int {
	return l.list.Len()
}

func (l *Simple[K, V]) Clear() {
	l.list.Init()
	l.data = make(map[K]*generics.MapElement[K, V])
}

func (l *Simple[K, V]) Join(glue string) string {
	if l.list == nil {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	length := l.list.Len()
	if length > 0 {
		for i, e := 0, l.list.Front(); i < length; i, e = i+1, e.Next() {
			if e == nil {
				continue
			}
			buffer.WriteString(fmt.Sprintf("%v", e.Value))
			if i != length-1 {
				buffer.WriteString(glue)
			}
		}
	}
	return buffer.String()
}

func (l *Simple[K, V]) String() string {
	if l == nil {
		return ""
	}
	return "[" + l.Join(",") + "]"
}
