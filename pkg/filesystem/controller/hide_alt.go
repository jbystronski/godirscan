package controller

import "github.com/jbystronski/godirscan/pkg/common"

func (c *Controller) hideAlt() {
	if !c.Active() {
		return
	}

	// c.Alt.hidden = !c.Alt.hidden

	if !c.Alt.hidden {
		c.Alt.hidden = true
		c.SetWidth(common.PaneWidth() * 2)
		c.SetOffsetLeftStart(1)
		c.fullRender()
		return

	}

	c.Alt.hidden = false
	c.SetWidth(common.PaneWidth())
	c.Alt.fullRender()
}
