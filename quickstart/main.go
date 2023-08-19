package main

import (
	"fmt"
	"github.com/shura1014/common/utils/stringutil"
	"github.com/shura1014/localcache"
	"time"
)

func main() {
	//lruCap()
	//lruExpire()
	lfuCap()
	//lfuExpire()
}

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

func lruExpire() {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(1024).Build()
	_ = local.Put("name", "shura", 1000)
	value, ok := local.Get("name")
	fmt.Println(value, ok)
	time.Sleep(time.Second)
	value, ok = local.Get("name")
	fmt.Println(value, ok)
}

func lfuExpire() {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(1024).Build()
	_ = local.Put("name", "shura", 1000)
	value, ok := local.Get("name")
	fmt.Println(value, ok)
	time.Sleep(time.Second)
	value, ok = local.Get("name")
	fmt.Println(value, ok)
}

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
