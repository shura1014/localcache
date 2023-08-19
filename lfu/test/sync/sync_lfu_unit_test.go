package sync

import (
	"github.com/shura1014/common/utils/stringutil"
	"github.com/shura1014/localcache"
	"testing"
	"time"
)

func Test_LfuAlg_Expire(t *testing.T) {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LFU).Cap(1024).Build()
	_ = local.Put("name", "shura", 1000)
	value, ok := local.Get("name")
	t.Log(value, ok)
	time.Sleep(time.Second)
	value, ok = local.Get("name")
	t.Log(value, ok)
}

func Test_LfuAlg_Cap(t *testing.T) {
	local := localcache.NewBuilder[string, string]().Alg(localcache.LFU).Cap(10).Build()
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
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
		_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
	}

	//for i := 16; i < 26; i++ {
	//	_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
	//	_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
	//	_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
	//	_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
	//	_ = local.Put(stringutil.Join("name", i), stringutil.Join("wendell", i), 10000)
	//}

	for i := 1; i < 26; i++ {
		value, ok := local.Get(stringutil.Join("name", i))
		t.Log(value, ok)
	}
	t.Log(local.DisPlay())
}
