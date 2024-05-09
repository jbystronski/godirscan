package filesystem

import (
	"sync"

	"github.com/jbystronski/pubsub"
)

var (
	once sync.Once
	n    *pubsub.Node
	c    *FsController
)

func New(root string) *pubsub.Node {
	once.Do(func() {
		n = pubsub.NewNode(pubsub.GlobalBroker())

		NewFsController(n)
	})

	c.root = root

	return n
}
