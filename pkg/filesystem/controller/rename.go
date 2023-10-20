package controller

import (
	"fmt"
	"os"
	"strings"
)

func (c *Controller) rename() {
	if c.DataAccessor.Len() == 0 {
		return
	}

	en, ok := c.Find(c.Index())

	if !ok {
		return
	}

	newName := c.wrapInput("Rename:", en.Name())

	if newName == "" || newName == en.Name() {
		return
	}

	if strings.Contains(newName, string(os.PathSeparator)) {
		c.ErrorChan <- fmt.Errorf("path separator can't be used inside name")
		return
	}

	err := en.Rename(newName)
	if err != nil {
		c.ErrorChan <- err
		return
	}

	c.render()
	if c.path == c.Alt.path {
		c.Alt.SetStore(c.Alt.path)
	}
}
