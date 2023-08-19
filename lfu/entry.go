package lfu

import (
	"github.com/shura1014/common/container/list/generics"
	"github.com/shura1014/common/utils/timeutil"
)

type lfuEntry[K comparable, V any] struct {
	key    K
	value  V
	expire int64

	freqList *generics.Element[*FreqEntry[K, V]]
}

type FreqEntry[K comparable, V any] struct {
	freq    int
	entries map[*lfuEntry[K, V]]byte
}

func (e *lfuEntry[K, V]) IsExpired() bool {
	if e.expire == 0 {
		return false
	}
	return e.expire < timeutil.MilliSeconds()
}
