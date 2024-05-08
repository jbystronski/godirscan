package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
)

// var (
// 	initResizeListener sync.Once
// 	resizeListenerNode    *e.Node
// 	r    *Resizer
// )

type ResizeListener struct {
	*pubsub.Node
	ch   chan os.Signal
	lock bool
}

func NewResizeListener() *pubsub.Node {
	// once.Do(func() {
	n := pubsub.NewNode()

	r := &ResizeListener{n, make(chan os.Signal, 1), false}

	r.Init()
	// })

	return n
}

func (r *ResizeListener) Init() {
	signal.Notify(r.ch, syscall.SIGWINCH)

	go func() {
		for {
			select {
			// case <-ctx.Done():
			// 	return

			case <-r.ch:

				if !r.lock {

					r.lock = true
					global.ClearScreen()
					time.Sleep(time.Millisecond * 500)
					r.Passthrough(pubsub.RESIZE, r.Next)

					r.lock = false

				}
			}
		}
	}()
}
