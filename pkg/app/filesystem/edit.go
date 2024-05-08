package filesystem

import (
	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/app/data"
)

func (c *FsController) edit() {
	if c.data.Len() == 0 {
		return
	}

	en, ok := c.activeEntry()

	if !ok {
		return
	}

	if en.FsType() == data.FileDatatype {
		c.executeCmd(config.Running().DefaultEditor, en.FullPath())
	}
}
