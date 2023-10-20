package controller

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jbystronski/godirscan/pkg/filesystem/utils"
)

func (c *Controller) Copy() {
	c.mv(c.path, false)
}

func (c *Controller) Move() {
	c.mv(c.path, true)
}

func (c *Controller) mv(currentPath string, remove bool) {
	defer c.selected.Clear()

	if c.selected.Len() == 0 {
		return
	}

	prompt := "Copy"

	if remove {
		prompt = "Move"
	}

	answ := c.wrapInput(fmt.Sprintf("%s%s", prompt, " selected into the current directory? :"), "y")

	if answ != "y" {
		return
	}

	sem := make(chan struct{}, 40)

	entries := c.selected.Copy()
	c.ctx.Create()
	_, tickerDone := c.ObserveTicker(c.ctx.Ctx.Done(), time.Millisecond*200, func() {
		if c.path == currentPath {
			c.fullRender()
		}
	})

	var wg sync.WaitGroup

	// wg2.Add(1)
	go func() {
		//	defer wg2.Done()
		for source := range entries {

			dir, file := filepath.Split(source)

			if dir == currentPath {
				c.ErrorChan <- errors.New("moving file within same directory")
				return

			}

			if strings.HasPrefix(dir, currentPath) {

				c.ErrorChan <- errors.New("cannot move a folder into itself")
				return

			}

			target := addSuffixIfExists(filepath.Join(currentPath, file))
			// wg.Add(1)
			// go func(source, target string) {
			// 	defer wg.Done()

			c.ctx.Observe(func() {
				c.copy(source, target, c.ErrorChan, sem)
				info, err := os.Lstat(target)
				if err != nil {
					c.ErrorChan <- err
				}

				c.insertEntry(c.path, info)
				c.SetChunk(c.DataAccessor.Len())
			})
			// }(source, target)

		}
		wg.Wait()
		tickerDone <- struct{}{}
		if c.path == currentPath {
			c.SetStore(currentPath)
		}

		if c.Alt.path == currentPath {
			time.Sleep(time.Millisecond * 200)
			c.Alt.SetStore(c.Alt.path)
		}
	}()

	if remove {

		for path := range entries {
			err := os.RemoveAll(path)
			if err != nil {
				c.ErrorChan <- err
			}

		}

		entries = nil
	}
}

func (c *Controller) copy(source, target string, errChan chan<- error, sem chan struct{}) {
	info, err := os.Stat(source)
	if err != nil {
		errChan <- err
		return
	}

	switch true {
	case info.IsDir():
		err := os.Mkdir(target, info.Mode())
		if err != nil {
			errChan <- err
			return
		}

		contents, err := os.ReadDir(source)
		if err != nil {
			errChan <- err
			return
		}

		for _, file := range contents {
			// c.WithContext(func() {
			c.copy(filepath.Join(source, file.Name()), filepath.Join(target, file.Name()), errChan, sem)
			// })
		}

	default:
		err := utils.Copy(source, target)
		if err != nil {
			c.ErrorChan <- err
			return
		}

		info, _ := os.Stat(source)

		err = os.Chmod(target, info.Mode())
		if err != nil {
			errChan <- err
			return
		}

	}
}

func (c *Controller) truncate(path string) {
	err := os.Truncate(path, 0)
	if err != nil {
		c.ErrorChan <- err
	}
}

func tryCreateSymlink(srcPath, targetPath string) (bool, error) {
	fileInfo, err := os.Lstat(srcPath)
	if err != nil {
		return false, err
	}
	if fileInfo.Mode()&os.ModeSymlink != 0 {
		symlinkTarget, linkErr := os.Readlink(srcPath)

		if linkErr != nil {
			return false, linkErr
		}

		symlinkErr := os.Symlink(symlinkTarget, targetPath)
		if symlinkErr != nil {
			return false, symlinkErr
		}

		return true, nil

	}
	return false, nil
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
