package sync

import (
	"github.com/shura1014/common/assert"
	"github.com/shura1014/localcache/lru"
	"testing"
)

func TestSize(t *testing.T) {
	lruMap := lru.NewLruMap[int, int](12)
	assert.Assert(lruMap.Size(), 0)
	lruMap.Put(1, 1)
	assert.Assert(lruMap.Size(), 1)
	lruMap.Put(1, 2)
	assert.Assert(lruMap.Size(), 1)
	lruMap.Put(2, 2)
	assert.Assert(lruMap.Size(), 2)
	lruMap.Put(3, 3)
	assert.Assert(lruMap.Size(), 3)
	lruMap.Evict(3)
	assert.Assert(lruMap.Size(), 2)
	t.Log(lruMap.Get(1))
	lruMap.Clear()
	assert.Assert(lruMap.Size(), 0)
	t.Log(lruMap.Get(1))

}
