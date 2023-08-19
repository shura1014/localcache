package sync

import (
	"github.com/shura1014/localcache/lru"
	"sync"
	"testing"
)

/*
go test . -bench . -count 3 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/shura1014/localcache/lru/test/sync
Benchmark_Sync_Map-8                     4575249               251.3 ns/op            48 B/op          3 allocs/op
Benchmark_Sync_Map-8                     5262469               271.8 ns/op            47 B/op          3 allocs/op
Benchmark_Sync_Map-8                     4752892               282.4 ns/op            48 B/op          3 allocs/op
Benchmark_Sync_LruMap-8                  4647603               271.6 ns/op            40 B/op          0 allocs/op
Benchmark_Sync_LruMap-8                  4343979               277.9 ns/op            41 B/op          0 allocs/op
Benchmark_Sync_LruMap-8                  3900535               270.0 ns/op            42 B/op          0 allocs/op
Benchmark_Sync_Map_Load-8               1000000000               1.017 ns/op           0 B/op          0 allocs/op
Benchmark_Sync_Map_Load-8               1000000000               1.037 ns/op           0 B/op          0 allocs/op
Benchmark_Sync_Map_Load-8               1000000000               1.018 ns/op           0 B/op          0 allocs/op
Benchmark_Sync_LruMap_Get-8             10447573               116.7 ns/op             0 B/op          0 allocs/op
Benchmark_Sync_LruMap_Get-8             10361385               115.6 ns/op             0 B/op          0 allocs/op
Benchmark_Sync_LruMap_Get-8             10676599               117.2 ns/op             0 B/op          0 allocs/op
Benchmark_Sync_LruMap_PutAndGet-8        3134730               404.7 ns/op            37 B/op          0 allocs/op
Benchmark_Sync_LruMap_PutAndGet-8        2810965               405.5 ns/op            38 B/op          0 allocs/op
Benchmark_Sync_LruMap_PutAndGet-8        3227448               400.3 ns/op            34 B/op          0 allocs/op
Benchmark_Sync_Map_LoadAndStore-8        2861736               474.4 ns/op            83 B/op          3 allocs/op
Benchmark_Sync_Map_LoadAndStore-8        2869120               452.5 ns/op            83 B/op          3 allocs/op
Benchmark_Sync_Map_LoadAndStore-8        2837672               482.8 ns/op            86 B/op          3 allocs/op
PASS
ok      github.com/shura1014/localcache/lru/test/sync   27.101s
*/
func Benchmark_Sync_Map(b *testing.B) {
	sm := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Store(i, i)
			i++
		}
	})
}

func Benchmark_Sync_LruMap(b *testing.B) {
	lruMap := lru.NewLruMap[int, int](1024)

	//for i := 0; i < 1000; i++ {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lruMap.Put(i, i)
			i++
		}
	})
	//}
}

func Benchmark_Sync_Map_Load(b *testing.B) {
	sm := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Load(i)
			i++
		}
	})
}

func Benchmark_Sync_LruMap_Get(b *testing.B) {
	lruMap := lru.NewLruMap[int, int](1024)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lruMap.Get(i)
			i++
		}

	})
}

func Benchmark_Sync_LruMap_PutAndGet(b *testing.B) {
	lruMap := lru.NewLruMap[int, int](1024)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			lruMap.Put(i, i)
			lruMap.Get(i)
			i++
		}

	})
}

func Benchmark_Sync_Map_LoadAndStore(b *testing.B) {
	sm := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Store(i, i)
			sm.Load(i)
			i++
		}

	})
}
