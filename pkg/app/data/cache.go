package data

import "github.com/jbystronski/godirscan/pkg/lib/maps"

type CacheStore struct {
	maps.MapAccessor[string, int]
	minSize int
}

func NewCacheStore(min int) *CacheStore {
	return &CacheStore{
		MapAccessor: &maps.GenericMap[string, int]{},
		minSize:     min,
	}
}

func (c *CacheStore) Set(k string, v int) {
	if v < c.minSize {
		return
	}

	c.MapAccessor.Set(k, v)
}
