package controller

import (
	"os"

	"github.com/jbystronski/godirscan/pkg/filesystem/utils"
)

func (c *Controller) left() {
	dir, ok := utils.GetParentDirectory(c.path)

	if !ok {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		c.ErrorChan <- err
	}

	len := len(entries)

	if len == 0 {
		return
	}

	if index, ok := c.backtrace.Get(dir); ok {

		if index > len-1 {
			c.SetIndex(len - 1)
		} else {
			c.SetIndex(index)
		}
		c.backtrace.Unset(dir)
		//	c.SetInitialData(dir)
		//	c.SetChunk(len)

		c.SetStore(dir)

		// c.DefaultSort()
		// c.fullRender()

	}
}
