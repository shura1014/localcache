package expire

import (
	"github.com/shura1014/common/utils/stringutil"
	"github.com/shura1014/localcache"
	"testing"
	"time"
)

func Test_LruAlg_Expire(t *testing.T) {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(1024).Build()
	_ = local.Put("name", "shura", 1000)
	value, ok := local.Get("name")
	t.Log(value, ok)
	time.Sleep(time.Second)
	value, ok = local.Get("name")
	t.Log(value, ok)
}

func Test_LruAlg_Cap(t *testing.T) {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LRU).Cap(16).EnableRecord().Build()
	for i := 0; i < 17; i++ {
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("shura", i), 10000)
	}
	for i := 0; i < 17; i++ {
		value, ok := local.Get(stringutil.Join("name", i))
		t.Log(value, ok)
	}
	t.Log(local.DisPlay())
}
