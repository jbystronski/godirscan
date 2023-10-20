package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jbystronski/godirscan/pkg/filesystem"
)

func (c *Controller) newDirectory() {
	newName := c.wrapInput("Create directory: ", "")

	en := &filesystem.FsEntry{}

	en.SetFsType(filesystem.Dir)

	if ok := c.Create(newName, en); ok {

		c.SetStore(c.path)

		time.Sleep(time.Millisecond * 200)

		c.Alt.SetStore(c.Alt.path)

		// info, _ := os.Stat(filepath.Join(c.FsStoreAccessor.Name(), newName))

		// c.updateParentStoreSize(int(info.Size()))

		// c.Alt.fullRender()
	}
}

func (c *Controller) Create(name string, en *filesystem.FsEntry) bool {
	var err error
	name = strings.TrimSpace(name)

	if name == "" {
		return false
	}

	if strings.ContainsAny(name, string(os.PathSeparator)) {
		c.ErrorChan <- fmt.Errorf("%s \"%v\"", "Name cannot contain", string(os.PathSeparator))
		return false
	}

	path := filepath.Join(c.path, name)

	switch en.FsType() {
	case filesystem.Dir:
		err = os.Mkdir(path, 0o777)

	case filesystem.File:
		_, err = os.Create(path)

	}

	if err != nil {
		c.ErrorChan <- err
		return false
	}

	return true
}
