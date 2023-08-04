package cache

import (
	"math"

	"github.com/jbystronski/godirscan/pkg/entry"
)

type cachedEntries struct {
	Size    int
	Entries []*entry.Entry
}

var (
	cache          = make(map[string]cachedEntries)
	cacheEntrySize = math.Pow(1024, 3)
)

func Clear() {
	cache = make(map[string]cachedEntries)
}

func Get(key string) (cachedEntries, bool) {
	if entries, exist := cache[key]; exist {
		return entries, true
	}

	return cachedEntries{}, false
}

func Store(key string, size int, data []*entry.Entry) {
	if size >= int(cacheEntrySize) {
		cache[key] = cachedEntries{Size: size, Entries: data}
	}
}
