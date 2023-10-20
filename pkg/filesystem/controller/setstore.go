package controller

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/jbystronski/godirscan/pkg/filesystem"
	"github.com/jbystronski/godirscan/pkg/filesystem/utils"
)

func (c *Controller) insertEntry(path string, info fs.FileInfo) {
	entry := c.pool.Get()

	switch true {

	case utils.IsSymlink(info):
		entry.SetFsType(filesystem.Symlink)

	case info.IsDir():
		entry.SetSize(int(info.Size()))
		entry.SetFsType(filesystem.Dir)

	default:
		entry.SetSize(int(info.Size()))
		entry.SetFsType(filesystem.File)

	}

	entry.SetName(filepath.Join(path, info.Name()))

	c.DataAccessor.Insert(entry)
}

func (c *Controller) SetInitialData(dir string) error {
	utils.ResolveUserDirectory(&dir)

	c.SetPath(dir)

	for _, v := range c.DataAccessor.All() {
		v.SetSize(0)
		v.SetName("")
		c.pool.Put(v)
	}

	dc, err := utils.ReadIgnorePermission(dir)
	if err != nil {
		return err
	}

	c.DataAccessor.Reset()

	for _, en := range dc {
		//	entry := c.pool.Get()

		info, err := os.Lstat(filepath.Join(dir, en.Name()))
		if err != nil {
			continue
		}

		c.insertEntry(dir, info)

		// switch true {

		// case utils.IsSymlink(info):
		// 	entry.SetFsType(filesystem.Symlink)

		// case info.IsDir():
		// 	entry.SetSize(int(info.Size()))
		// 	entry.SetFsType(filesystem.Dir)

		// default:
		// 	entry.SetSize(int(info.Size()))
		// 	entry.SetFsType(filesystem.File)

		// }

		// entry.SetName(filepath.Join(dir, info.Name()))
		// c.DataAccessor.Insert(entry)

	}

	if c.DataAccessor.Len() > 0 && c.Index() > c.DataAccessor.Len()-1 {
		c.SetIndex(c.DataAccessor.Len() - 1)
	}

	c.SetChunk(c.DataAccessor.Len())

	c.fullRender()
	return nil
}

func (c *Controller) ScanSize(path string) (total int) {
	if utils.IsVirtualFs(path) {
		return
	}
	pathInfo, err := os.Stat(path)
	if err != nil {
		c.ErrorChan <- err
		return
	}

	total += int(pathInfo.Size())

	contents, err := utils.ReadIgnorePermission(path)
	if err != nil {
		c.ErrorChan <- err
		// error = err

		return
	}

	for _, dirEntry := range contents {

		info, err := dirEntry.Info()
		if err != nil {
			// error = err

			return
		}

		if dirEntry.IsDir() {
			c.ctx.Observe(func() {
				size := c.ScanSize(filepath.Join(path, info.Name()))
				total += size
			})
		}

		total += int(info.Size())

	}

	return
}

func (c *Controller) CalculateStoreSize(dir string) {
	var wg sync.WaitGroup

	go func() {
		for _, e := range c.DataAccessor.All() {
			switch e.FsType() {
			case filesystem.Dir:

				if size, ok := c.cache.Get(e.FullPath()); ok {
					e.SetSize(size)
				} else {

					wg.Add(1)

					go func(e *filesystem.FsEntry) {
						defer wg.Done()
						c.ctx.Observe(func() {
							size := c.ScanSize(e.FullPath())

							e.SetSize(size)
							c.cache.Set(e.FullPath(), e.Size())
						})
					}(e)
				}
			}
		}
		wg.Wait()

		c.DefaultSort()

		c.fullRender()
	}()
}

func (c *Controller) SetStore(dir string) {
	c.ctx.Create()

	err := c.SetInitialData(dir)
	if err != nil {
		c.ErrorChan <- err
	}

	c.CalculateStoreSize(dir)
}
