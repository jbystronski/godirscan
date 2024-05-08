package filesystem

import (
	"os"
	"runtime"
	"strings"

	"github.com/jbystronski/godirscan/pkg/app/config"
)

func (c *FsController) openBookmark(b string) {
	info, _ := os.Stat(b)

	switch true {
	case info.IsDir():

		err := c.Active().setStore(b)
		if err != nil {
			c.sendError(err)
		}

	default:
		var args string
		switch runtime.GOOS {
		case "darwin":
			{
				args = strings.Join([]string{"open"}, " ")
			}
		case "windows":
			{
				args = strings.Join([]string{"cmd", "/c", "start"}, " ")
			}
		default:
			{
				args = strings.Join([]string{"xdg-open"}, " ")
			}
		}
		c.executeCmd(args, b)

	}
}

func (c *FsController) bookmark(group string) {
	if en, ok := c.activeEntry(); ok {
		path := c.getInput("press enter to bookmark this entry:", en.FullPath())

		if path == "" {
			return
		}

		if _, err := os.Stat(path); err != nil {

			c.sendError(err)
			return
		}

		config.Running().AddBookmark(group, path)

	}
}

func (c *FsController) removeBookmark(group, bookmark string) {
	answ := c.getInput("remove bookmark "+bookmark, "y")

	if answ != "y" {
		return
	}

	config.Running().RemoveBookmark(group, bookmark)
}

func (c *FsController) addBookmarkGroup() {
	answ := c.getInput("add bookmark group ", "")

	if answ == "" {
		return
	}

	config.Running().AddBookmarkGroup(answ)
}

func (c *FsController) removeBookmarkGroup(name string) {
	answ := c.getInput("remove group "+name, "y")

	if answ != "y" {
		return
	}
	config.Running().RemoveBookmarkGroup(name)
}
