package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jbystronski/godirscan/pkg/app/data"
)

func (c *FsController) newFile() {
	response := c.getInput("Create a new file", "")

	if response == "" {
		return
	}

	path := filepath.Join(c.root, response)

	_, err := os.Stat(path)

	if err == nil && response != "" {

		answ := c.getInput(fmt.Sprintf("%s (%s) %s", "File", response, "already exists, do you wish to override it?"), "n")

		if answ == "y" || answ == "Y" {
			truncate(path)
		} else {
			c.render()
			c.alt.render()

		}

	} else {

		en := &data.FsEntry{}

		en.SetFsType(data.FileDatatype)

		err := c.Create(response, en)

		if err != nil {
			c.sendError(err)
		} else {
			c.setStore(c.root)

			c.alt.setStore(c.alt.root)

			c.render()
			c.alt.render()

		}

	}
}
