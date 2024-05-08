package filesystem

import (
	"sync"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
)

var (
	once sync.Once
	n    *event.Node
	c    *FsController
)

func New(root string) *event.Node {
	once.Do(func() {
		n = event.NewNode()

		NewFsController(n)
	})

	c.root = root

	return n
}
