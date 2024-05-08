package boxes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub/message"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type ProgressController struct {
	*event.Node

	currentMsg string
	view       termui.Section
	sync.Mutex
}

func NewProgressBox(cancelCtx context.CancelFunc) *event.Node {
	n := event.NewNode()

	c := ProgressController{
		n,

		"",
		termui.NewSection(),
		sync.Mutex{},
	}
	updateDimensions()

	c.view.SetHeight(1).SetWidth(cols())
	c.view.SetTop(rows())
	c.view.SetLeft(1)
	// c.view.SetBorder().SetPadding(2, 2, 2, 2).SetHeight(6).SetWidth(cols() / 2).CenterVertically().CenterHorizontally()

	n.On(event.ESC, func() {
		cancelCtx()
	})

	n.On(event.RENDER, func() {
		// c.currentMsg = string(c.DequeueMessage("process_message"))

		c.render()
		// c.render(string(c.DequeueMessage("process_message")))
	})

	n.Subscribe("progress_message", func(m message.Message) {
		// c.currentMsg = string(m)
		c.Lock()
		c.printMessage(string(m))
		c.Unlock()
		// c.update(string(m))

		// n.EnqueueMessage("progress_message", m)
	})

	return n
}

func (c *ProgressController) updateDimensions() {
}

func (c *ProgressController) printMessage(msg string) {
	clear(rows(), 1, cols())
	// clear(c.view.ContentFirstLine(), c.view.ContentStart(), c.view.ContentWidth())
	// msg = trimEnd(msg, c.view.ContentWidth(), c.view.ContentWidth(), 2, '.')
	// c.view.Content = []string{c.currentMsg, "Press ESC to abort"}

	fmt.Print(fmtBold(trimEnd(msg, c.view.ContentWidth(), c.view.ContentWidth(), 2, '.')))

	// c.view.Content = []string{fmtBold(msg)}
	// c.view.PrintContent()
}

func (c *ProgressController) render() {
	time.Sleep(time.Millisecond * 100)

	updateDimensions()
	//	c.view.Width = cols() / 2

	c.view.SetWidth(cols())
	//	c.view.Print(themeMain())
	//
	// c.printMessage(m)
}

// func (c *ProgressController) update(msg string) {
// 	c.currentMsg = msg
// 	c.printMessage()
// }
