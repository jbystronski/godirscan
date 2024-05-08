package filesystem

import (
	"sync"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
)

var (
	once sync.Once
	n    *pubsub.Node
	c    *FsController
)

func New(root string) *pubsub.Node {
	once.Do(func() {
		n = pubsub.NewNode()

		NewFsController(n)
	})

	c.root = root

	return n
}
