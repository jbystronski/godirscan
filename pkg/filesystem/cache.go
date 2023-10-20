package filesystem

import (
	"github.com/jbystronski/godirscan/pkg/common"
)

type CacheStore struct {
	common.MapAccessor[string, int]
	minSize int
}

func NewCacheStore(min int) *CacheStore {
	return &CacheStore{
		MapAccessor: &common.GenericMap[string, int]{},
		minSize:     min,
	}
}

func (c *CacheStore) Set(k string, v int) {
	if v < c.minSize {
		return
	}

	c.MapAccessor.Set(k, v)
}
