package localcache

type ICache[K comparable, V any] interface {
	Put(key K, value V, expire ...int64) error
	PutIfAbsent(key K, value V, expire ...int64) error
	PutAll(data map[K]V, expire ...int64) error
	Get(key K) (v V, ok bool)
	Contains(key K) bool
	Evict(key K) (v V)
	Clear()
	Size() int
	DisPlay() string
	ExpireCleanTask()
}
