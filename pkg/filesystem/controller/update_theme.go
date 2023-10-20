package controller

import "github.com/jbystronski/godirscan/pkg/common"

func (c *Controller) updateTheme(config *common.Config) {
	if c.Active() {

		config.ChangeTheme()

		c.SetTheme(common.CurrentTheme)
		c.Alt.SetTheme(common.CurrentTheme)
		c.PrintBox()
		c.Alt.PrintBox()
		c.fullRender()
		c.Alt.fullRender()
	}
}
