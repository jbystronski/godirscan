package controller

import "github.com/jbystronski/godirscan/pkg/common"

func (c *Controller) Refresh() {
	c.SetTotalLines(common.NumVisibleLines())
	c.SetChunkLines(c.Lines())
	c.SetChunk(c.DataAccessor.Len())

	c.SetWidth(common.PaneWidth())

	if c.OffsetLeftStart() > 1 {
		c.SetOffsetLeftStart(common.PaneWidth() + 1)
	}
	c.PrintBox()
	c.fullRender()
}
