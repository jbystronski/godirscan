package boxes

import (
	"time"

	"github.com/jbystronski/godirscan/pkg/lib/pubsub"

	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

type ErrorController struct {
	*pubsub.Node

	view       termui.Section
	errMessage string
}

func NewError(err string) *pubsub.Node {
	n := pubsub.NewNode()

	c := ErrorController{
		n,

		termui.NewSection(),

		err,
	}
	c.view.SetBorder().SetPadding(1, 2, 1, 2).SetHeight(8).SetWidth(cols() / 3).CenterHorizontally().CenterVertically()

	n.OnGlobal(pubsub.RESIZE, func() {
		c.view.CenterVertically().CenterHorizontally()

		c.Print()
	})

	n.OnGlobal(pubsub.T, c.Print)

	n.On(pubsub.Q, func() {
		n.Unlink()
		n.Passthrough(pubsub.RENDER, n.First())
	})

	n.On(pubsub.RENDER, c.Print)

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
