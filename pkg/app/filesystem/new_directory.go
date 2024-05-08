package filesystem

import (
	"github.com/jbystronski/godirscan/pkg/app/data"
)

func (c *FsController) newDirectory() {
	response := c.getInput("Create directory", "")

	if response == "" {
		return
	}

	en := &data.FsEntry{}

	en.SetFsType(data.DirDatatype)

	if err := c.Create(response, en); err == nil {

		c.setStore(c.root)

		c.alt.setStore(c.alt.root)
	} else {
		c.sendError(err)
	}
}
