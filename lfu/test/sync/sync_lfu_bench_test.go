package sync

import (
	"github.com/shura1014/localcache"
	"testing"
)

/*
go test . -bench . -count 3 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/shura1014/localcache/lfu/test/sync
Benchmark_SyncLfu_Put-8          2726065               559.3 ns/op            40 B/op          1 allocs/op
Benchmark_SyncLfu_Put-8          2359809               658.6 ns/op            40 B/op          1 allocs/op
Benchmark_SyncLfu_Put-8          1951159               603.9 ns/op            40 B/op          1 allocs/op
Benchmark_SyncLfu_Get-8         10107876               118.0 ns/op             0 B/op          0 allocs/op
Benchmark_SyncLfu_Get-8         10296013               118.3 ns/op             0 B/op          0 allocs/op
Benchmark_SyncLfu_Get-8         10297777               117.1 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/shura1014/localcache/lfu/test/sync   13.472s
*/
func Benchmark_SyncLfu_Put(b *testing.B) {
	local := localcache.NewBuilder[int, int]().Alg(localcache.LFU).Cap(1024).Build()
	b.RunParallel(func(pb *testing.PB) {
		i := 1
		for pb.Next() {
			_ = local.Put(i, 1, 1000)
			i++
		}
	})
}

func Benchmark_SyncLfu_Get(b *testing.B) {
	local := localcache.NewBuilder[int, int]().Alg(localcache.LFU).Cap(1024).Build()
	b.RunParallel(func(pb *testing.PB) {
		i := 1
		for pb.Next() {
			local.Get(i)
			i++
		}
	})
}
