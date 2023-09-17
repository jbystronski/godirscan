package filesystem

import (
	"fmt"

	"github.com/jbystronski/godirscan/pkg/common"
)

type FsNavigator struct {
	common.ChunkNavigator
	backtrace map[string]int
}

func (n *FsNavigator) SetBacktrace(key string) {
	if n.backtrace == nil {
		n.backtrace = make(map[string]int)
	}

	n.backtrace[key] = n.Index()
}

func (n *FsNavigator) Backtrace(key string, len int) bool {
	if len == 0 {
		return false
	}

	if index, ok := n.backtrace[key]; ok {
		if index > len-1 {
			n.SetIndex(len - 1)
		} else {
			n.SetIndex(index)
		}
		n.SetChunk(len)
		delete(n.backtrace, key)

		return true

	}

	return false
}

func (n *FsNavigator) Print() {
	for k, v := range n.backtrace {
		fmt.Println("path: ", k, " index: ", v)
	}
}
