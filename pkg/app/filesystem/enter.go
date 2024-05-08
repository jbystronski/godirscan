package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jbystronski/godirscan/pkg/app/config"
	"github.com/jbystronski/godirscan/pkg/app/data"
	"github.com/jbystronski/godirscan/pkg/global"
)

func (c *FsController) executeFile(en *data.FsEntry) {
	fp := en.FullPath()
	ext := filepath.Ext(fp)
	var input string

	if customCmd, ok := config.Running().Executors[ext]; ok {
		input = strings.Join([]string{customCmd}, " ")
	} else {
		switch runtime.GOOS {
		case "darwin":
			{
				input = strings.Join([]string{"open"}, " ")
			}
		case "windows":
			{
				input = strings.Join([]string{"cmd", "/c", "start"}, " ")
			}
		default:
			{
				input = strings.Join([]string{"xdg-open"}, " ")
			}
		}
	}
	c.executeCmd(input, fp)
}

func (c *FsController) execute() {
	en, ok := c.activeEntry()

	if !ok {
		return
	}

	switch en.FsType() {

	case data.DirDatatype:

		c.right()

	case data.FileDatatype:

		c.executeFile(en)

	case data.SymlinkDatatype:
		linkPath, err := os.Readlink(en.FullPath())
		if err != nil {
			c.sendError(err)

			return
		}

		path := filepath.Join(c.root, linkPath)

		info, err := os.Stat(path)
		if err != nil {
			c.sendError(err)

			return
		}

		if !info.IsDir() {
			c.executeFile(en)
		} else {

			parts := strings.Split(linkPath, string(os.PathSeparator))
			part := c.root

			c.backtrace.Set(part, c.Index())
			if part == global.GetRootDirectory() {
				part = ""
			}

			c.Navigator.Reset()
			for _, v := range parts {
				part += string(os.PathSeparator) + v
				c.backtrace.Set(part, c.Index())
			}
			c.panel.Print(themeMain())
			c.alt.panel.Print(themeMain())
			c.setStore(part)

			c.render()
		}

	}
}
