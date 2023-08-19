package segment

import (
	"github.com/shura1014/common"
	"github.com/shura1014/common/assert"
	"github.com/shura1014/localcache/lru"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_Segment_Size(t *testing.T) {
	lruMap := lru.NewHashMap[int, int](1024, 64)
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

// NewLruMap[int, int](1024, 16)
// === RUN   TestConcurrent
//
//	segment_lru_test.go:41: start
//	segment_lru_test.go:43: 写线程start
//	segment_lru_test.go:48: 耗时 1790 毫秒
//	segment_lru_test.go:51: 1024
//	segment_lru_test.go:53: 读线程start
//	segment_lru_test.go:58: 耗时 845 毫秒
func TestConcurrent(t *testing.T) {
	lruMap := lru.NewHashMap[int, int](1024, 64)
	t.Log("start")
	go func() {
		t.Log("写线程start")
		start := time.Now()
		for i := 0; i < 10000000; i++ {
			lruMap.Put(i, i)
		}
		t.Logf("耗时 %v 毫秒", time.Now().Sub(start).Milliseconds())
	}()
	time.Sleep(8 * time.Second)
	t.Log(lruMap.Size())
	go func() {
		t.Log("读线程start")
		start := time.Now()
		for i := 0; i < 10000000; i++ {
			lruMap.Get(i)
		}
		t.Logf("耗时 %v 毫秒", time.Now().Sub(start).Milliseconds())
	}()

	common.Wait()
}

// === RUN   TestSyncMap
//
//	segment_lru_test.go:76: start
//	segment_lru_test.go:78: 写线程start
//	segment_lru_test.go:83: 耗时 5374 毫秒
//	segment_lru_test.go:88: 读线程start
//	segment_lru_test.go:93: 耗时 1730 毫秒
func TestSyncMap(t *testing.T) {
	syncMap := sync.Map{}
	t.Log("start")
	go func() {
		t.Log("写线程start")
		start := time.Now()
		for i := 0; i < 8000000; i++ {
			syncMap.Store(strconv.Itoa(i), i)
		}
		t.Logf("耗时 %v 毫秒", time.Now().Sub(start).Milliseconds())
	}()

	time.Sleep(10 * time.Second)
	go func() {
		t.Log("读线程start")
		start := time.Now()
		for i := 0; i < 8000000; i++ {
			syncMap.Load(strconv.Itoa(i))
		}
		t.Logf("耗时 %v 毫秒", time.Now().Sub(start).Milliseconds())
	}()

	common.Wait()
}
