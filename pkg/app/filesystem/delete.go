package filesystem

import (
	"errors"
	"os"
)

func (c *FsController) delete() error {
	if c.root == c.selected.BasePath() {
		c.selected.Clear()
	}

	for key := range c.selected.Self() {
		if error := os.RemoveAll(key); error != nil {
			if errors.Is(error, os.ErrNotExist) {
				return error
			}
			c.selected.Toggle(key)
		}
	}

	return nil
}
