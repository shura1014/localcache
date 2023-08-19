package expire

import (
	"github.com/shura1014/localcache"
	"testing"
)

/*
go test . -bench . -count 3 -benchmem
goos: darwin
goarch: arm64
pkg: github.com/shura1014/localcache/lru/test/expire
Benchmark_Lru_Expire_Put-8       3952137               315.6 ns/op            64 B/op          2 allocs/op
Benchmark_Lru_Expire_Put-8       4373032               331.6 ns/op            64 B/op          2 allocs/op
Benchmark_Lru_Expire_Put-8       4280791               316.4 ns/op            61 B/op          2 allocs/op
Benchmark_Lru_Expire_Get-8      10133486               119.0 ns/op             0 B/op          0 allocs/op
Benchmark_Lru_Expire_Get-8      10197744               119.4 ns/op             0 B/op          0 allocs/op
Benchmark_Lru_Expire_Get-8      10125225               120.7 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/shura1014/localcache/lru/test/expire 12.373s
*/
func Benchmark_Lru_Expire_Put(b *testing.B) {
	local := localcache.NewBuilder[int, int]().Alg(localcache.LRU).Cap(1024).Build()
	b.RunParallel(func(pb *testing.PB) {
		i := 1
		for pb.Next() {
			_ = local.Put(i, 1, 1000)
			i++
		}
	})
}

func Benchmark_Lru_Expire_Get(b *testing.B) {
	local := localcache.NewBuilder[int, int]().Alg(localcache.LRU).Cap(1024).Build()
	b.RunParallel(func(pb *testing.PB) {
		i := 1
		for pb.Next() {
			local.Get(i)
			i++
		}
	})
}
