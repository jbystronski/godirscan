package controller

import (
	"fmt"

	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem/utils"
)

func (c *Controller) updateTotalSize() {
	c.GoToCell(c.TotalLines()-2, c.Alt.ContentLineStart())
	fmt.Print(c.TotalSize(c.DataAccessor.Size()))
}

func (c *Controller) updateEntrySize(index, size int) {
	l := c.Line(index)
	c.ClearLine(l, c.ContentLineStart(), 12)
	c.GoToCell(l, c.ContentLineStart())
	fmt.Print(c.PrintSizeAsString(size))
}

func (c *Controller) updateParentStoreSize(updatedSize int) {
	path := c.path

	for path != common.GetRootDirectory() {
		if parentDir, ok := utils.GetParentDirectory(path); ok {
			if size, ok := c.cache.Get(parentDir); ok {
				c.cache.Set(parentDir, size+updatedSize)
				path = parentDir

				continue
			}
			return
		}
	}
}
