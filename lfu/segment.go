package lfu

import (
	"bytes"
	"github.com/shura1014/common/hash"
	"github.com/shura1014/common/utils/stringutil"
)

const (
	DefaultCap         = 1024
	DefaultSegmentSize = 16
	MaxCap             = 1 << 30
)

type HashMap[K comparable, V any] struct {
	segments []*Cache[K, V]
	mask     uint32
}

func NewHashMap[K comparable, V any](args ...int) *HashMap[K, V] {
	var (
		capacity    = DefaultCap
		segmentSize = DefaultSegmentSize
	)

	if len(args) > 0 {
		capacity = args[0]
	}
	if len(args) > 1 {
		segmentSize = args[1]
	}
	if capacity <= 0 {
		capacity = DefaultCap
	}
	// capacity = 33
	// capacity = 48
	if 0 != capacity%segmentSize {
		capacity = (capacity/segmentSize + 1) * segmentSize
	}

	// 找到最佳的以为运算
	// DefaultSegmentSize = 16 那么 bestShift = 28  任何一个int值又移动28为必介于1-15之间，可以方便的寻找段
	size := 1
	for size < segmentSize {
		size <<= 1
	}

	lruMap := &HashMap[K, V]{
		segments: make([]*Cache[K, V], size),
		mask:     uint32(size - 1),
	}

	if capacity > MaxCap {
		capacity = MaxCap
	}
	c := capacity / size
	segmentCap := 1
	// segmentCap要求是2的幂 capacity = 48 DefaultSegmentSize=16 c=48/16=3 segmentCap=4
	for segmentCap < c {
		segmentCap <<= 1
	}

	for i := 0; i < len(lruMap.segments); i++ {
		lruMap.segments[i] = NewLfu[K, V](segmentCap)
	}
	return lruMap
}

func (l *HashMap[K, V]) Put(key K, value V, expire ...int64) error {
	return l.segments[l.getIndex(key)].Put(key, value, expire...)
}
func (l *HashMap[K, V]) PutIfAbsent(key K, value V, expire ...int64) error {
	return l.segments[l.getIndex(key)].PutIfAbsent(key, value, expire...)
}
func (l *HashMap[K, V]) PutAll(data map[K]V, expire ...int64) error {
	for key, value := range data {
		_ = l.Put(key, value, expire...)
	}
	return nil
}
func (l *HashMap[K, V]) Get(key K) (v V, ok bool) {
	return l.segments[l.getIndex(key)].Get(key)
}

func (l *HashMap[K, V]) Contains(key K) (ok bool) {
	return l.segments[l.getIndex(key)].Contains(key)
}

func (l *HashMap[K, V]) Evict(key K) V {
	return l.segments[l.getIndex(key)].Evict(key)
}

func (l *HashMap[K, V]) Clear() {
	for _, segment := range l.segments {
		segment.Clear()
	}
}

func (l *HashMap[K, V]) Size() int {
	sum := 0
	for i := 0; i < len(l.segments); i++ {
		sum += l.segments[i].Size()
	}
	return sum
}

func (l *HashMap[K, V]) DisPlay() string {
	buffer := bytes.NewBuffer(nil)

	for _, segment := range l.segments {
		buffer.WriteString(segment.DisPlay())
		buffer.WriteByte('\n')
	}

	return buffer.String()
}
func (l *HashMap[K, V]) ExpireCleanTask() {
	for _, segment := range l.segments {
		segment.ExpireCleanTask()
	}

}

func (l *HashMap[K, V]) getIndex(key K) uint32 {
	s := stringutil.ToString(key)
	index := Hash(hash.BKDRHash32(stringutil.StringToBytes(s)))
	return index & l.mask
}

func Hash(h uint32) uint32 {
	h += (h << 15) ^ 0xFFFFCD7D
	h ^= h >> 10
	h += h << 3
	h ^= h >> 6
	h += (h << 2) + (h << 14)
	return h ^ (h >> 16)
}
