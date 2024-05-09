package filesystem

import "github.com/jbystronski/godirscan/pkg/app/data"

func (c *FsController) right() {
	if c.data.Len() == 0 {
		return
	}

	en, ok := c.activeEntry()

	if !ok {
		return
	}

	if en.FsType() == data.DirDatatype {
		c.backtrace.Set(c.root, c.Index())
		c.Navigator.Reset()

		c.setStore(en.FullPath())
	}
}
