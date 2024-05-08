package boxes

import (
	"fmt"
	"sync"

	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/lib/pubsub"
	"github.com/jbystronski/godirscan/pkg/lib/termui"
)

var (
	initResizeWarning    sync.Once
	resizeWarnNode       *pubsub.Node
	resizeWarnController *ResizeWarnController
)

type ResizeWarnController struct {
	*pubsub.Node
	outer termui.Section
	inner termui.Section
}

func NewResizeWarning() *pubsub.Node {
	initResizeWarning.Do(func() {
		resizeWarnNode = pubsub.NewNode()

		resizeWarnController = &ResizeWarnController{
			resizeWarnNode,
			termui.NewSection(),
			termui.NewSection(),
		}

		resizeWarnNode.On(pubsub.RENDER, resizeWarnController.print)

		resizeWarnNode.On(pubsub.ESC, func() {
			resizeWarnNode.Passthrough(pubsub.QUIT_APP, resizeWarnNode.First())
		})
	})

	return resizeWarnNode
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
