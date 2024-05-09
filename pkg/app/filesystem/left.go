package filesystem

import (
	"os"

	"github.com/jbystronski/godirscan/pkg/global"
)

func (c *FsController) left() {
	dir, ok := global.GetParentDirectory(c.root)

	if !ok {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		c.sendError(err)
	}

	len := len(entries)

	if len == 0 {
		return
	}

	if index, ok := c.backtrace.Get(dir); ok {

		if index > len-1 {
			c.SetIndex(len - 2)
		} else {
			c.SetIndex(index)
		}
		c.backtrace.Unset(dir)

		c.setStore(dir)

	}
}
