package filesystem

import (
	"math"
)

var (
	cache          = make(map[string]FsStore)
	cacheEntrySize = math.Pow(1000, 3)
)

func Clear() {
	cache = make(map[string]FsStore)
}

func Get(key string) (FsStore, bool) {
	if store, exist := cache[key]; exist {
		return store, true
	}

	return FsStore{}, false
}

func Set(store FsStore) {
	if store.Size() >= int(cacheEntrySize) {
		cache[store.Name()] = store
	}
}
