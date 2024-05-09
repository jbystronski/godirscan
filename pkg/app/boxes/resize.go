package boxes

import (
	"fmt"

	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/global/event"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
	"github.com/jbystronski/pubsub"
)

type ResizeWarnController struct {
	*pubsub.Node
	outer termui.Section
	inner termui.Section
}

func NewResizeWarning() *pubsub.Node {
	n := pubsub.NewNode(pubsub.GlobalBroker())

	c := &ResizeWarnController{
		n,
		termui.NewSection(),
		termui.NewSection(),
	}

	n.On(event.RENDER, c.print)

	n.On(event.ESC, func() {
		n.Passthrough(event.QUIT_APP, n.First())
	})

	return n
}

func (c *ResizeWarnController) print() {
	hideCursor()
	cls()

	c.outer.SetBorder().SetWidth(cols()).SetHeight(rows())

	c.inner.SetWidth(cols() / 2).SetHeight(4).CenterHorizontally().CenterVertically()

	c.outer.Print(themeMain())

	content := []string{}
	content = append(content, fmtBold(alignCenter(c.inner.ContentWidth(), "Plesase resize your termui", " ")))

	content = append(content, fmtBold(alignCenter(c.inner.ContentWidth(), fmt.Sprintf("Required width %d, has %d", config.MIN_WIDTH, cols()), " ")))

	c.inner.Content = content
	c.inner.PrintContent()
}
