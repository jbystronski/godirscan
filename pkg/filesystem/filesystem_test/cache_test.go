package filesystem_test

import (
	"testing"

	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func TestCache(t *testing.T) {
	min := 104857600
	cache := filesystem.NewCacheStore(min)
	cache.Set("entry1", min+1)
	cache.Set("entry2", min+500)
	cache.Set("entry3", min+min*2)

	t.Run("Cache, test Len", func(t *testing.T) {
		got := cache.Len()
		want := 3
		if want != got {
			t.Errorf("testing Len, want %v, got %v", want, got)
		}
	})
}
