package controller

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func (c *Controller) newFile() {
	newFileName := c.wrapInput("Create a new file: ", "")

	path := filepath.Join(c.path, newFileName)

	_, err := os.Stat(path)

	if err == nil && newFileName != "" {
		answ := c.wrapInput(fmt.Sprintf("%s (%s) %s", "File", newFileName, "already exists, do you wish to override it?"), "n")
		if answ == "y" || answ == "Y" {
			c.truncate(path)
		} else {
			c.render()
			c.Alt.render()

		}

	} else {

		en := &filesystem.FsEntry{}

		en.SetFsType(filesystem.File)

		ok := c.Create(newFileName, en)

		if ok {

			c.SetStore(c.path)

			c.Alt.SetStore(c.Alt.path)

			if err != nil {
				c.ErrorChan <- err
			}

			c.fullRender()
			c.Alt.fullRender()
		} else {
			c.render()
			c.Alt.render()
		}
	}
}
