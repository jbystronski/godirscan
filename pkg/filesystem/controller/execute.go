package controller

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jbystronski/godirscan/pkg/common"
	"github.com/jbystronski/godirscan/pkg/filesystem"
	"github.com/jbystronski/godirscan/pkg/filesystem/utils"
)

func (c *Controller) execute() {
	en, ok := c.Find(c.Index())

	if !ok {
		return
	}

	switch en.FsType() {

	case filesystem.Dir:

		c.right()

	case filesystem.File:
		err := utils.ExecuteFile(en.FullPath())
		if err != nil {
			c.ErrorChan <- err
		}
		c.restoreScreen()
	case filesystem.Symlink:
		linkPath, err := os.Readlink(en.FullPath())
		if err != nil {
			c.ErrorChan <- err
			return
		}

		path := filepath.Join(c.path, linkPath)

		info, err := os.Stat(path)
		if err != nil {
			c.ErrorChan <- err
			return
		}

		if !info.IsDir() {
			return
		}

		parts := strings.Split(linkPath, string(os.PathSeparator))
		part := c.path

		c.backtrace.Set(part, c.Index())
		if part == common.GetRootDirectory() {
			part = ""
		}

		c.Reset()
		for _, v := range parts {
			part += string(os.PathSeparator) + v
			c.backtrace.Set(part, c.Index())
		}
		c.PrintBox()
		c.Alt.PrintBox()
		c.SetStore(part)

		c.DefaultSort()
		c.fullRender()

	}
}
