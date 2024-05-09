package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jbystronski/godirscan/pkg/app/boxes"
	"github.com/jbystronski/godirscan/pkg/global"
	"github.com/jbystronski/godirscan/pkg/global/event"
	"github.com/jbystronski/pubsub"
)

func (c *FsController) Copy(deleteAfter bool) error {
	var errInstance error
	currentPath := c.root

	if c.HasNext() {
		c.Next().Unlink()
	}

	sem := make(chan struct{}, 40)

	entries := c.selected.Copy()

	// ctx := utils.NewCancelContext()

	c.ctx.Create()

	progress := boxes.NewProgressBox(c.ctx.CancelFunc)

	c.LinkTo(progress)
	c.Passthrough(event.RENDER, c.Next())

	for source := range entries {

		dir, file := filepath.Split(source)

		if dir == currentPath {
			c.ctx.CancelFunc()

			errInstance = errors.New("moving file within same directory")
			break

		}

		if strings.HasPrefix(dir, currentPath) {
			c.ctx.CancelFunc()
			errInstance = errors.New("cannot move a folder into itself")
			break

		}

		target := addSuffixIfExists(filepath.Join(currentPath, file))

		c.ctx.Observe(func() {
			c.copy(source, target, sem)
			info, err := os.Lstat(target)
			if err != nil {
				errInstance = err
				return
			}

			c.insertEntry(c.root, info)
			c.Publish("progress_message", pubsub.Message("copying "+source))

			c.selected.Toggle(source)

			if deleteAfter {

				err := os.RemoveAll(source)
				if err != nil {
					errInstance = err
					return
				}

			}
		})

	}

	c.Next().Unlink()

	if c.root == currentPath {
		c.setStore(currentPath)
	}

	if c.alt.root == currentPath {
		c.alt.setStore(c.alt.root)
	}

	return errInstance
}

func (c *FsController) copy(source, target string, sem chan struct{}) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	switch true {
	case info.IsDir():
		err := os.Mkdir(target, info.Mode())
		if err != nil {
			return err
		}

		contents, err := os.ReadDir(source)
		if err != nil {
			return err
		}

		for _, file := range contents {
			c.copy(filepath.Join(source, file.Name()), filepath.Join(target, file.Name()), sem)
		}

	default:

		err := global.Copy(source, target)
		if err != nil {
			return err
		}

		info, _ := os.Stat(source)

		err = os.Chmod(target, info.Mode())
		if err != nil {
			return err
		}

	}
	return nil
}

func addSuffixIfExists(targetPath string) string {
	_, err := os.Stat(targetPath)
	if err != nil {
		return targetPath
	}

	for num := 1; ; num++ {
		targetCopy := targetPath + " COPY " + strconv.Itoa(num)
		_, err := os.Stat(targetCopy)
		if err != nil {
			return targetCopy
		}
	}
}
