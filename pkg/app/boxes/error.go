package boxes

import (
	"time"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub/event"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type ErrorController struct {
	*event.Node

	view       termui.Section
	errMessage string
}

func NewError(err string) *event.Node {
	n := event.NewNode()

	c := ErrorController{
		n,

		termui.NewSection(),

		err,
	}
	c.view.SetBorder().SetPadding(1, 2, 1, 2).SetHeight(8).SetWidth(cols() / 3).CenterHorizontally().CenterVertically()

	n.OnGlobal(event.RESIZE, func() {
		c.view.CenterVertically().CenterHorizontally()

		c.Print()
	})

	n.OnGlobal(event.T, c.Print)

	n.On(event.Q, func() {
		n.Unlink()
		n.Passthrough(event.RENDER, n.First())
	})

	n.On(event.RENDER, c.Print)

	return n
}

func (c *ErrorController) Print() {
	time.Sleep(time.Millisecond * 100)

	updateDimensions()

	c.view.Content = []string{c.formatRow("ERROR"), c.formatRow(c.errMessage), "\n", c.formatRow("Press 'q' to dismiss")}

	c.view.Print(themeMain())
	c.view.PrintContent()
}

func (c *ErrorController) formatRow(s string) string {
	return buildString(fmtBold(alignCenter(c.view.ContentWidth(), s, " ")))
}
