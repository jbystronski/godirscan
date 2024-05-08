package filesystem

import (
	"errors"
	"os"
	"strings"
)

func (c *FsController) rename() {
	if c.data.Len() == 0 {
		return
	}

	en, ok := c.activeEntry()

	if !ok {
		return
	}

	newName := c.getInput("Rename", en.Name())

	if newName == "" || newName == en.Name() {
		return
	}

	if strings.Contains(newName, string(os.PathSeparator)) {
		c.sendError(errors.New("path separator can't be used inside name"))

		return
	}

	err := en.Rename(newName)
	if err != nil {
		c.sendError(err)

		return
	}

	c.render()
	if c.root == c.alt.root {
		c.alt.setStore(c.alt.root)
	}
}
