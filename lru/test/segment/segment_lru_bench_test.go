package segment

import (
	"github.com/shura1014/localcache/lru"
	"strconv"
	"sync"
	"testing"
)

/*
go test . -bench . -count 3 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/shura1014/localcache/lru/test/segment
Benchmark_Segment_LruMap_Put-8                   9788989               121.0 ns/op            60 B/op          2 allocs/op
Benchmark_Segment_LruMap_Put-8                  10695516               123.5 ns/op            62 B/op          2 allocs/op
Benchmark_Segment_LruMap_Put-8                  10688553               120.8 ns/op            61 B/op          2 allocs/op
Benchmark_Segment_SyncMap_Store-8                5067858               234.9 ns/op            47 B/op          3 allocs/op
Benchmark_Segment_SyncMap_Store-8                5082151               253.4 ns/op            47 B/op          3 allocs/op
Benchmark_Segment_SyncMap_Store-8                4805108               241.4 ns/op            48 B/op          3 allocs/op
Benchmark_Segment_LruMap_Get-8                  32287502                37.50 ns/op           16 B/op          1 allocs/op
Benchmark_Segment_LruMap_Get-8                  30411794                37.50 ns/op           16 B/op          1 allocs/op
Benchmark_Segment_LruMap_Get-8                  30742035                37.31 ns/op           16 B/op          1 allocs/op
Benchmark_Segment_SyncMap_Load-8                1000000000               1.012 ns/op           0 B/op          0 allocs/op
Benchmark_Segment_SyncMap_Load-8                1000000000               1.011 ns/op           0 B/op          0 allocs/op
Benchmark_Segment_SyncMap_Load-8                1000000000               1.015 ns/op           0 B/op          0 allocs/op
Benchmark_Segment_LruMap_GetAndPut-8             9525494               152.4 ns/op            77 B/op          4 allocs/op
Benchmark_Segment_LruMap_GetAndPut-8             8260795               153.3 ns/op            76 B/op          4 allocs/op
Benchmark_Segment_LruMap_GetAndPut-8             8599512               150.6 ns/op            74 B/op          4 allocs/op
Benchmark_Segment_SyncMap_LoadAndStore-8         2456072               576.8 ns/op            97 B/op          4 allocs/op
Benchmark_Segment_SyncMap_LoadAndStore-8         2425150               520.1 ns/op            95 B/op          4 allocs/op
Benchmark_Segment_SyncMap_LoadAndStore-8         2511205               500.6 ns/op            96 B/op          4 allocs/op
PASS
ok      github.com/shura1014/localcache/lru/test/segment        25.860s
*/
func Benchmark_Segment_LruMap_Put(b *testing.B) {
	b.ReportAllocs()
	hashMap := lru.NewHashMap[int, int](1024, 64)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Put(i, i)
			i++
		}
	})
}

// Benchmark_Segment_SyncMap_Store-8   	 4019116	       309.4 ns/op	      50 B/op	       3 allocs/op
func Benchmark_Segment_SyncMap_Store(b *testing.B) {
	syncMap := sync.Map{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			syncMap.Store(i, i)
			i++
		}
	})
}

// Benchmark_Segment_LruMap_Get-8   	27959509	        41.31 ns/op	      16 B/op	       1 allocs/op
func Benchmark_Segment_LruMap_Get(b *testing.B) {

	hashMap := lru.NewHashMap[int, int](1024, 64)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Get(i)
			i++
		}
	})
}

// Benchmark_Segment_SyncMap_Load-8   	994049503	         1.085 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Segment_SyncMap_Load(b *testing.B) {
	syncMap := sync.Map{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			syncMap.Load(i)
			i++
		}
	})
}

// Benchmark_Segment_LruMap_GetAndPut-8   	 7490858	       164.0 ns/op	      76 B/op	       4 allocs/op
func Benchmark_Segment_LruMap_GetAndPut(b *testing.B) {
	hashMap := lru.NewHashMap[int, int](1024, 64)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Put(i, i)
			hashMap.Get(i)
			i++
		}
	})
}

// Benchmark_Segment_SyncMap_LoadAndStore-8   	 2171938	       580.3 ns/op	      96 B/op	       4 allocs/op
func Benchmark_Segment_SyncMap_LoadAndStore(b *testing.B) {
	syncMap := sync.Map{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			syncMap.Store(i, i)
			syncMap.Load(strconv.Itoa(i))
			i++
		}
	})
}
