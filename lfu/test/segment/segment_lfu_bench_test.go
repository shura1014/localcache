package segment

import (
	"github.com/shura1014/localcache/lfu"
	"strconv"
	"sync"
	"testing"
)

/*
go test . -bench . -count 3 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/shura1014/localcache/lfu/test/segment
Benchmark_Segment_LfuCache_Put-8                 7792705               150.2 ns/op            48 B/op          3 allocs/op
Benchmark_Segment_LfuCache_Put-8                 8119678               149.3 ns/op            48 B/op          3 allocs/op
Benchmark_Segment_LfuCache_Put-8                 8269495               148.0 ns/op            48 B/op          3 allocs/op
Benchmark_Segment_SyncMap_Store-8                5204067               271.5 ns/op            47 B/op          3 allocs/op
Benchmark_Segment_SyncMap_Store-8                5670402               238.7 ns/op            46 B/op          3 allocs/op
Benchmark_Segment_SyncMap_Store-8                5300437               236.1 ns/op            47 B/op          3 allocs/op
Benchmark_Segment_LfuCache_Get-8                33102306                32.37 ns/op           15 B/op          1 allocs/op
Benchmark_Segment_LfuCache_Get-8                36224090                32.51 ns/op           16 B/op          1 allocs/op
Benchmark_Segment_LfuCache_Get-8                36060987                33.40 ns/op           15 B/op          1 allocs/op
Benchmark_Segment_SyncMap_Load-8                1000000000               1.013 ns/op           0 B/op          0 allocs/op
Benchmark_Segment_SyncMap_Load-8                1000000000               1.043 ns/op           0 B/op          0 allocs/op
Benchmark_Segment_SyncMap_Load-8                1000000000               1.141 ns/op           0 B/op          0 allocs/op
Benchmark_Segment_LfuCache_GetAndPut-8           4395160               279.5 ns/op           251 B/op          8 allocs/op
Benchmark_Segment_LfuCache_GetAndPut-8           4369868               280.1 ns/op           252 B/op          8 allocs/op
Benchmark_Segment_LfuCache_GetAndPut-8           4336438               287.3 ns/op           252 B/op          8 allocs/op
Benchmark_Segment_SyncMap_LoadAndStore-8         2168714               497.4 ns/op            98 B/op          4 allocs/op
Benchmark_Segment_SyncMap_LoadAndStore-8         2387632               539.5 ns/op           100 B/op          4 allocs/op
Benchmark_Segment_SyncMap_LoadAndStore-8         2181237               581.2 ns/op            96 B/op          4 allocs/op
PASS
ok      github.com/shura1014/localcache/lfu/test/segment        26.235s
*/
func Benchmark_Segment_LfuCache_Put(b *testing.B) {
	hashMap := lfu.NewHashMap[int, int](1024, 64)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = hashMap.Put(i, i, 1000)
			i++
		}
	})
}

func Benchmark_Segment_SyncMap_Store(b *testing.B) {
	syncMap := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			syncMap.Store(i, i)
			i++
		}
	})
}

func Benchmark_Segment_LfuCache_Get(b *testing.B) {

	hashMap := lfu.NewHashMap[int, int](1024, 64)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Get(i)
			i++
		}
	})
}

func Benchmark_Segment_SyncMap_Load(b *testing.B) {
	syncMap := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			syncMap.Load(i)
			i++
		}
	})
}

func Benchmark_Segment_LfuCache_GetAndPut(b *testing.B) {
	hashMap := lfu.NewHashMap[int, int](1024, 64)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = hashMap.Put(i, i, 1000)
			hashMap.Get(i)
			i++
		}
	})
}

func Benchmark_Segment_SyncMap_LoadAndStore(b *testing.B) {
	syncMap := sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			syncMap.Store(i, i)
			syncMap.Load(strconv.Itoa(i))
			i++
		}
	})
}
