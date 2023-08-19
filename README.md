# localcache

本地缓存
lru 简单的lru 并发安全的lru 分段lru 带过期时间的lru
lfu 简单的lfu 并发安全的lfu 分段lfu 带过期时间的lfu

# 快速使用

## lru

### 带过期时间
```go
func lruExpire() {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(1024).Build()
	_ = local.Put("name", "shura", 1000)
	value, ok := local.Get("name")
	fmt.Println(value, ok)
	time.Sleep(time.Second)
	value, ok = local.Get("name")
	fmt.Println(value, ok)
}

shura true
 false
```

### cap容量 
超出容量清理最近没用到的

```go
func lruCap() {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(16).EnableRecord().Build()
	for i := 0; i < 17; i++ {
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("shura", i), 10000)
	}
	for i := 0; i < 17; i++ {
		value, ok := local.Get(stringutil.Join("name", i))
		fmt.Println(value, ok)
	}
	fmt.Println(local.DisPlay())
}
```

```text
 false
shura1 true
shura2 true
shura3 true
shura4 true
shura5 true
shura6 true
shura7 true
shura8 true
shura9 true
shura10 true
shura11 true
shura12 true
shura13 true
shura14 true
shura15 true
shura16 true
len: 16 hit: 16 miss: 1

容量16 最近没用到的为第一个，缓存命中 16 没命中 1
```

## lfu

### 带过期时间

```go
func lfuExpire() {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(1024).Build()
	_ = local.Put("name", "shura", 1000)
	value, ok := local.Get("name")
	fmt.Println(value, ok)
	time.Sleep(time.Second)
	value, ok = local.Get("name")
	fmt.Println(value, ok)
}

shura true
 false
```

### cap容量

清理最近使用频次最低的

```go
func lfuCap() {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LFU).Cap(10).EnableRecord().Build()
	_ = local.Put("name1", "shura1", 10000)
	_ = local.Put("name1", "shura1", 10000)
	_ = local.Put("name1", "shura1", 10000)
	_ = local.Put("name1", "shura1", 10000)
	_ = local.Put("name1", "shura1", 10000)

	_ = local.Put("name2", "shura2", 10000)
	_ = local.Put("name2", "shura2", 10000)

	_ = local.Put("name3", "shura3", 10000)
	_ = local.Put("name3", "shura3", 10000)
	_ = local.Put("name3", "shura3", 10000)
	_ = local.Put("name3", "shura3", 10000)

	for i := 4; i < 16; i++ {
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("shura", i), 10000)
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("shura", i), 10000)
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("shura", i), 10000)
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("shura", i), 10000)
	}

	for i := 1; i < 16; i++ {
		value, ok := local.Get(stringutil.Join("name", i))
		fmt.Println(value, ok)
	}
	fmt.Println(local.DisPlay())
}
```

结果

```text
shura1 true
 false
 false
shura4 true
 false
shura6 true
shura7 true
 false
shura9 true
shura10 true
shura11 true
shura12 true
 false
shura14 true
shura15 true
len: 10 hit: 10 miss: 5

由于shura1最近使用频次为5没被清除，shura2 shura3 使用频次较低，已被清除
容量 10 命中10 没命中5
```


## 分段lru

并发测试

原生同步sync.Map与分段lru测试

写操作有明显的性能提升

```text
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
```