package controller

import (
	"os"
	"os/exec"

	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func (c *Controller) edit(editor string) {
	if c.DataAccessor.Len() == 0 {
		return
	}

	en, ok := c.Find(c.Index())

	if !ok {
		return
	}

	if en.FsType() == filesystem.File {
		common.ClearScreen()
		cmd := exec.Command(editor, en.FullPath())

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			c.ErrorChan <- err
			return
		}
		common.HideCursor()
		c.restoreScreen()

	}
}
